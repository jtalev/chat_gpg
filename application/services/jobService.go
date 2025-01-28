package application

import (
	"database/sql"

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
