package application

import (
	"database/sql"
	"log"

	"github.com/jtalev/chat_gpg/domain/models"
	"github.com/jtalev/chat_gpg/infrastructure/repository"
)

func GetJobById(id int, db *sql.DB) (domain.Job, error) {
	job, err := infrastructure.GetJobById(id, db)
	if err != nil {
		return domain.Job{}, err
	}
	return job, nil
}

func GetJobs(db *sql.DB) ([]domain.Job, error) {
	outJobs, err := infrastructure.GetJobs(db)
	if err != nil {
		return outJobs, err
	}
	return outJobs, nil
}

func PostJob(inJob domain.Job, db *sql.DB) (domain.Job, error) {
	outJob, err := infrastructure.PostJob(inJob, db)
	if err != nil {
		log.Printf("Error posting job: %v", err)
		return outJob, err
	}
	return outJob, err
}
