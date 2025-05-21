package task_queue

import (
	"context"
	"database/sql"
	"encoding/json"
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

type DbBackupPayload struct {
	NextRunAt time.Time `json:"next_run_at"`
}

func CreateDbBackupPayload(nextRunAt time.Time) DbBackupPayload {
	return DbBackupPayload{
		NextRunAt: nextRunAt,
	}
}

type DbBackupHandler struct {
}

func init() {
	RegisterTaskHandler("db_backup", &DbBackupHandler{})
}

func (d *DbBackupHandler) ProcessTask(task Task, queue chan Task, db *sql.DB) error {
	var p DbBackupPayload
	err := json.Unmarshal(task.Payload, &p)
	if err != nil {
		return err
	}

	if p.NextRunAt.After(time.Now()) {
		queue <- task
		return nil
	}

	err = backupDb()
	if err != nil {
		return err
	}

	p.NextRunAt = p.NextRunAt.AddDate(0, 0, 7)

	payload, err := json.Marshal(p)
	if err != nil {
		return err
	}
	task.Payload = payload

	err = UpdateTask(task, db)
	if err != nil {
		return err
	}
	queue <- task

	log.Println("db backed up successfully")

	return nil
}

const (
	sourceDb   = "../../chat_gpg/infrastructure/db/dev.db"
	backupDir  = "backups"
	bucketName = "db-backup-chat-gpg"
	region     = "ap-southeast-2"
)

type backup struct {
	BackupFile     *os.File
	BackupFileName string
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

func backupDb() error {
	var b backup

	err := b.createBackup()
	if err != nil {
		return err
	}

	err = b.uploadToS3()
	if err != nil {
		return err
	}

	return nil
}
