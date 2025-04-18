package application

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

const (
	bucketName = "db-backup-chat-gpg"
	region     = "ap-southeast-2"
)

type Pdf struct {
	UUID string

	JsonTemplatePath string
	OutJsonName      string
	OutJsonPath      string
	OutJsonFile      *os.File

	InPdfPath  string
	OutPdfName string
	OutPdfPath string
	OutPdfFile *os.File

	S3FileName   string
	S3StorageDir string
	S3Key        string

	Data any
}

var alphabet = map[string]string{
	"a": "A",
	"b": "B",
	"c": "C",
	"d": "D",
	"e": "E",
	"f": "F",
	"g": "G",
	"h": "H",
	"i": "I",
	"j": "J",
	"k": "K",
	"l": "L",
	"m": "M",
	"n": "N",
	"o": "O",
	"p": "P",
	"q": "Q",
	"r": "R",
	"s": "S",
	"t": "T",
	"u": "U",
	"v": "V",
	"w": "W",
	"x": "X",
	"y": "Y",
	"z": "Z",
}

func (p *Pdf) FormatJsonTemplate() error {
	inputFile, err := os.Open(p.JsonTemplatePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	templateVariable := "."
	var outputLines []string

	scanner := bufio.NewScanner(inputFile)
	const maxCapacity = 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, `"value"`) {
			switch {
			case strings.HasSuffix(line, "false,"):
				line = strings.TrimSuffix(line, "false,")
				line = fmt.Sprintf(`%s{{ %s }},`, line, templateVariable)
			case strings.HasSuffix(line, "true,"):
				line = strings.TrimSuffix(line, "true,")
				line = fmt.Sprintf(`%s{{ %s }},`, line, templateVariable)
			case strings.HasSuffix(line, `"",`):
				line = strings.TrimSuffix(line, `"",`)
				line = fmt.Sprintf(`%s"{{ %s }}",`, line, templateVariable)
			default:
				fmt.Println("nothing trimmed")
				outputLines = append(outputLines, line)
				continue
			}

			templateVariable = "."
			outputLines = append(outputLines, line)

		} else if strings.HasPrefix(trimmed, `"name"`) {
			arr := strings.Split(trimmed, `"`)
			if len(arr) < 4 {
				outputLines = append(outputLines, line)
				continue
			}
			temp := arr[3]
			parts := strings.Split(temp, "_")
			for i, word := range parts {
				if len(word) == 0 {
					continue
				}
				first := string(word[0])
				upper, ok := alphabet[first]
				if ok {
					parts[i] = upper + word[1:]
				}
				templateVariable += parts[i]
			}
			outputLines = append(outputLines, line)

		} else {
			outputLines = append(outputLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error at line %d: %w", len(outputLines), err)
	}

	outputFile, err := os.Create(p.JsonTemplatePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	for _, line := range outputLines {
		_, err = writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func (p *Pdf) GenerateJsonTemplate() error {
	conf := api.LoadConfiguration()
	err := api.ExportFormFile(p.InPdfPath, p.JsonTemplatePath, conf)
	if err != nil {
		return err
	}
	return nil
}

// returns path of processed json file
// json used to fill pdf
func (p *Pdf) ExecuteJsonTemplate() error {
	if p.JsonTemplatePath == "" {
		panic("type Pdf requires member InJsonPath to execute json templating")
	}

	inJsonFile, err := os.ReadFile(p.JsonTemplatePath)
	if err != nil {
		return err
	}

	p.OutJsonName = strings.Split(strings.Split(p.JsonTemplatePath, "/")[len(strings.Split(p.JsonTemplatePath, "/"))-1], ".")[0]

	dotIdx := len(p.JsonTemplatePath) - 5
	p.OutJsonPath = p.JsonTemplatePath[:dotIdx] + "_output" + p.JsonTemplatePath[dotIdx:]
	p.OutJsonFile, err = os.Create(p.OutJsonPath)
	if err != nil {
		return err
	}
	defer p.OutJsonFile.Close()

	t := template.Must(template.New(p.OutJsonName).Parse(string(inJsonFile)))
	err = t.Execute(p.OutJsonFile, p.Data)
	if err != nil {
		return err
	}

	return nil
}

func setCfg() (aws.Config, error) {
	s3CredKey := os.Getenv("S3_KEY")
	s3CredSecret := os.Getenv("S3_SECRET")

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			aws.NewCredentialsCache(
				credentials.NewStaticCredentialsProvider(s3CredKey, s3CredSecret, ""))))
	if err != nil {
		return aws.Config{}, err
	}
	return cfg, nil
}

func (p *Pdf) Store() error {
	if p.OutPdfName == "" || p.UUID == "" || p.S3StorageDir == "" || p.OutPdfPath == "" {
		log.Println(p.OutPdfName)
		panic("type Pdf needs further configuration")
	}

	p.S3FileName = fmt.Sprintf("%s_%s.pdf", p.UUID, p.OutPdfName)
	p.S3Key = fmt.Sprintf("%s/%s", p.S3StorageDir, p.S3FileName)

	cfg, err := setCfg()
	if err != nil {
		log.Printf("error loading AWS config: %v", err)
		return err
	}

	client := s3.NewFromConfig(cfg)

	file, err := os.Open(p.OutPdfPath)
	if err != nil {
		log.Printf("error opening pdf to be stored: %v", err)
		return err
	}
	defer file.Close()

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(p.S3Key),
		Body:   file,
	})

	log.Printf("%s updloaded to %s/%s", p.S3FileName, bucketName, p.S3StorageDir)

	return nil
}

func (p *Pdf) GeneratePdf() error {
	if p.InPdfPath == "" || p.OutPdfPath == "" || p.OutJsonPath == "" {
		panic("type Pdf requires member InPdfPath, OutPdfPath & OutJsonPath to generate Pdf")
	}

	err := p.ExecuteJsonTemplate()
	if err != nil {
		log.Printf("failed to execute json template: %v", err)
		return err
	}

	err = api.FillFormFile(p.InPdfPath, p.OutJsonPath, p.OutPdfPath, nil)
	if err != nil {
		log.Printf("failed to fill form file: %v", err)
		return err
	}

	err = p.Store()
	if err != nil {
		log.Printf("error storing pdf: %v", err)
		return err
	}
	return nil
}

func (p *Pdf) GetPresignedURL(expirySeconds int64) (string, error) {
	if p.S3FileName == "" || p.S3StorageDir == "" {
		panic("type Pdf requires S3FileName and S3StorageDir to generate a pre-signed URL")
	}

	p.S3Key = fmt.Sprintf("%s/%s", p.S3StorageDir, p.S3FileName)

	cfg, err := setCfg()
	if err != nil {
		log.Printf("error loading AWS config: %v", err)
		return "", err
	}

	client := s3.NewFromConfig(cfg)

	presigner := s3.NewPresignClient(client)

	req, err := presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket:              aws.String(bucketName),
		Key:                 aws.String(p.S3Key),
		ResponseContentType: aws.String("application/pdf"),
	}, s3.WithPresignExpires(time.Duration(expirySeconds)*time.Second))

	if err != nil {
		log.Printf("error generating pre-signed URL: %v", err)
		return "", err
	}

	return req.URL, nil
}

func (p *Pdf) Delete() error {
	if p.S3FileName == "" || p.S3StorageDir == "" {
		panic("type Pdf requires S3FileName and S3StorageDir to delete s3 object")
	}

	p.S3Key = fmt.Sprintf("%s/%s", p.S3StorageDir, p.S3FileName)

	cfg, err := setCfg()
	if err != nil {
		log.Printf("error loading AWS config: %v", err)
		return err
	}

	client := s3.NewFromConfig(cfg)

	_, err = client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(p.S3Key),
	})
	if err != nil {
		log.Printf("error deleting s3 object: %v", err)
		return err
	}
	return nil
}
