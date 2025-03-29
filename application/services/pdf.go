package application

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

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

	InJsonPath  string
	OutJsonName string
	OutJsonPath string
	OutJsonFile *os.File

	InPdfPath  string
	OutPdfName string
	OutPdfPath string
	OutPdfFile *os.File

	S3FileName   string
	S3StorageDir string
	S3Key        string

	Data any
}

// returns path of processed json file
// json used to fill pdf
func (p *Pdf) ExecuteJsonTemplate() error {
	if p.InJsonPath == "" {
		panic("type Pdf requires member InJsonPath to execute json templating")
	}

	inJsonFile, err := os.ReadFile(p.InJsonPath)
	if err != nil {
		return err
	}

	p.OutJsonName = strings.Split(strings.Split(p.InJsonPath, "/")[len(strings.Split(p.InJsonPath, "/"))-1], ".")[0]

	dotIdx := len(p.InJsonPath) - 5
	p.OutJsonPath = p.InJsonPath[:dotIdx] + "_output" + p.InJsonPath[dotIdx:]
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

func (p *Pdf) Store() error {
	if p.OutPdfName == "" || p.UUID == "" || p.S3StorageDir == "" || p.OutPdfPath == "" {
		log.Println(p.OutPdfName)
		panic("type Pdf needs further configuration")
	}

	p.S3FileName = fmt.Sprintf("%s_%s.pdf", p.UUID, p.OutPdfName)
	p.S3Key = fmt.Sprintf("%s/%s", p.S3StorageDir, p.S3FileName)

	s3CredKey := os.Getenv("S3_KEY")
	s3CredSecret := os.Getenv("S3_SECRET")

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			aws.NewCredentialsCache(
				credentials.NewStaticCredentialsProvider(s3CredKey, s3CredSecret, ""))))
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
	if p.InPdfPath == "" {
		panic("type Pdf requires member InPdfPath to generate Pdf")
	}

	err := p.ExecuteJsonTemplate()
	if err != nil {
		return err
	}

	err = api.FillFormFile(p.InPdfPath, p.OutJsonPath, p.OutPdfPath, nil)
	if err != nil {
		return err
	}

	p.Store()
	return nil
}
