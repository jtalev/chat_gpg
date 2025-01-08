package services

import (
	"database/sql"

	"github.com/jtalev/chat_gpg/models"
	"github.com/jtalev/chat_gpg/repository"
)

func GetJobById(id int, db *sql.DB) (models.Job, error) {
	job, err := repository.GetJobById(id, db)
	if err != nil {
		return models.Job{}, err
	}
	return job, nil
}
