package application

import (
	"database/sql"
	"log"
	"strconv"

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

type JobDto struct {
	ID          string
	Name        string
	Number      string
	Address     string
	Suburb      string
	PostCode    string
	City        string
	IsComplete  string
	NameErr     string
	NumberErr   string
	AddressErr  string
	SuburbErr   string
	PostCodeErr string
	CityErr     string
	SuccessMsg  string
}

func jobDtoToJob(jobDto JobDto) (domain.Job, error) {
	id, err := strconv.Atoi(jobDto.ID)
	if err != nil {
		return domain.Job{}, nil
	}
	number, err := strconv.Atoi(jobDto.Number)
	if err != nil {
		jobDto.NumberErr = "Numbers only"
		return domain.Job{}, nil
	}
	isComplete := false
	if jobDto.IsComplete == "true" {
		isComplete = true
	}
	job := domain.Job{
		ID:         id,
		Name:       jobDto.Name,
		Number:     number,
		Address:    jobDto.Address,
		Suburb:     jobDto.Suburb,
		PostCode:   jobDto.PostCode,
		City:       jobDto.City,
		IsComplete: isComplete,
	}

	return job, nil
}

func mapErrorsToJobDto(errors domain.JobErrors, jobDto JobDto) JobDto {
	jobDto.NameErr = errors.NameErr
	jobDto.NumberErr = errors.NumberErr
	jobDto.AddressErr = errors.AddressErr
	jobDto.SuburbErr = errors.SuburbErr
	jobDto.PostCodeErr = errors.PostCodeErr
	jobDto.CityErr = errors.CityErr

	return jobDto
}

func PutJob(jobDto JobDto, db *sql.DB) (JobDto, error) {
	job, err := jobDtoToJob(jobDto)
	if err != nil {
		return JobDto{}, err
	}

	_, err = strconv.Atoi(jobDto.Number)
	if err != nil {
		jobDto.NumberErr = "Numbers only"
	}

	errors := job.Validate()
	if errors.IsSuccessful == false {
		jobDto = mapErrorsToJobDto(errors, jobDto)
		return jobDto, nil
	} else {
		job, err = infrastructure.PutJob(job.ID, job, db)
		jobDto.SuccessMsg = "Job successfully updated."
		log.Println(job.IsComplete)
		return jobDto, nil
	}
}
