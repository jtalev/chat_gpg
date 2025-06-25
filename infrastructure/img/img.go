package img

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	bucketName = "chat-gpg-img-store"
	region     = "ap-southeast-2"
)

type ImgStore struct {
	bucketName string
	region     string
}

func InitImgStore() ImgStore {
	return ImgStore{
		bucketName: bucketName,
		region:     region,
	}
}

func setConfig() (aws.Config, error) {
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

func (i *ImgStore) Store(imgPath, imgUuid, s3Dir string) error {
	s3FileName := fmt.Sprintf("%s.jpg", imgUuid)
	s3Key := fmt.Sprintf("%s/%s", s3Dir, s3FileName)

	cfg, err := setConfig()
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)

	file, err := os.Open(imgPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(i.bucketName),
		Key:    aws.String(s3Key),
		Body:   file,
	})

	log.Printf("%s uploaded to %s %s", s3FileName, i.bucketName, s3Dir)

	return nil
}

func (i *ImgStore) GetImgUrl(uuid, s3Dir string) (string, error) {
	s3FileName := fmt.Sprintf("%s.jpg", uuid)
	s3Key := fmt.Sprintf("%s/%s", s3Dir, s3FileName)

	cfg, err := setConfig()
	if err != nil {
		log.Printf("error loading AWS config: %v", err)
		return "", err
	}

	client := s3.NewFromConfig(cfg)

	presigner := s3.NewPresignClient(client)

	req, err := presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(s3Key),
	}, s3.WithPresignExpires(time.Duration(100)*time.Second))

	if err != nil {
		log.Printf("error generating pre-signed URL: %v", err)
		return "", err
	}

	return req.URL, nil
}
