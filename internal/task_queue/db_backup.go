package task_queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	sourceDb   = "../chat_gpg/infrastructure/db/prod.db"
	backupDir  = "backups"
	bucketName = "db-backup-chat-gpg"
	region     = "ap-southeast-2"
)

type backup struct {
	BackupFile     *os.File
	BackupFileName string

	NextRunAt time.Time
}

func (b *backup) createBackup() error {
	timestamp := time.Now().Format(time.DateOnly)
	backupFileName := fmt.Sprintf("backup_%s.db", timestamp)

	srcFile, err := os.Open(sourceDb)
	if err != nil {
		log.Printf("error opening src file: %v", err)
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(backupFileName)
	if err != nil {
		log.Printf("error creating destination file: %v", err)
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		log.Printf("error copying source to dest file: %v", err)
		return err
	}

	b.BackupFile = destFile
	b.BackupFileName = backupFileName
	return nil
}

func (b *backup) uploadToS3() error {
	s3key := os.Getenv("S3_KEY")
	s3secret := os.Getenv("S3_SECRET")
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			aws.NewCredentialsCache(
				credentials.NewStaticCredentialsProvider(s3key, s3secret, ""))))
	if err != nil {
		log.Printf("error loading AWS config: %v", err)
		return err
	}

	client := s3.NewFromConfig(cfg)

	file, err := os.Open(b.BackupFileName)
	if err != nil {
		log.Printf("error opening backup file: %v", err)
		return err
	}
	defer file.Close()

	key := fmt.Sprintf("%v/%v", backupDir, b.BackupFileName)

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		log.Printf("error uploading object to s3 bucket: %v", err)
		return err
	}

	log.Println("Backup uploaded to s3:", key)

	return nil
}

func BackupDb(data []byte) error {
	var b backup
	err := json.Unmarshal(data, &b)
	if err != nil {
		return err
	}
	if b.NextRunAt.After(time.Now()) {
		return errors.New("task not scheduled to run yet")
	}
	err = b.createBackup()
	if err != nil {
		return err
	}
	err = b.uploadToS3()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	HandlerRegistry["backup_db"] = BackupDb
}
