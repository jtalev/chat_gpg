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
		return domain.Job{}, err
	}
	isComplete := false
	if jobDto.IsComplete == "true" {
		isComplete = true
	}
	job := domain.Job{
		ID:         id,
		Name:       jobDto.Name,
		Address:    jobDto.Address,
		Suburb:     jobDto.Suburb,
		PostCode:   jobDto.PostCode,
		City:       jobDto.City,
		IsComplete: isComplete,
	}

	if jobDto.Number == "" {
		job.Number = 0
	} else {
		number, err := strconv.Atoi(jobDto.Number)
		if err != nil {
			job.Number = number
			jobDto.NumberErr = "Numbers only"
			return job, nil
		}
		job.Number = number
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

func fillJobProperties(job domain.Job) domain.Job {
	if job.Address == "" {
		job.Address = "n/a"
	}
	if job.Suburb == "" {
		job.Suburb = "n/a"
	}
	if job.PostCode == "" {
		job.PostCode = "n/a"
	}
	if job.City == "" {
		job.City = "n/a"
	}
	return job
}

func fillJobDtoProperties(jobDto JobDto) JobDto {
	if jobDto.Number == "" {
		jobDto.Number = "0"
	}
	if jobDto.Address == "" {
		jobDto.Address = "n/a"
	}
	if jobDto.Suburb == "" {
		jobDto.Suburb = "n/a"
	}
	if jobDto.PostCode == "" {
		jobDto.PostCode = "n/a"
	}
	if jobDto.City == "" {
		jobDto.City = "n/a"
	}
	return jobDto
}

func PostJob(jobDto JobDto, db *sql.DB) (JobDto, error) {
	job, err := jobDtoToJob(jobDto)
	if err != nil {
		return JobDto{}, err
	}

	errors := job.Validate()
	if errors.IsSuccessful == false {
		jobDto = mapErrorsToJobDto(errors, jobDto)
		jobDto = fillJobDtoProperties(jobDto)
		log.Println(jobDto)
		return jobDto, nil
	} else {
		job = fillJobProperties(job)
		job, err = infrastructure.PostJob(job, db)
		if err != nil {
			return JobDto{}, err
		}
		jobDto.SuccessMsg = "Job submitted successfully."
		return jobDto, nil
	}
}

func PutJob(jobDto JobDto, db *sql.DB) (JobDto, error) {
	job, err := jobDtoToJob(jobDto)
	if err != nil {
		return JobDto{}, err
	}

	errors := job.Validate()
	if errors.IsSuccessful == false {
		jobDto = mapErrorsToJobDto(errors, jobDto)
		jobDto = fillJobDtoProperties(jobDto)
		return jobDto, nil
	} else {
		job = fillJobProperties(job)
		job, err = infrastructure.PutJob(job.ID, job, db)
		jobDto.SuccessMsg = "Job successfully updated."
		return jobDto, nil
	}
}
