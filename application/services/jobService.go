package application

import (
	"database/sql"
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

	number, err := strconv.Atoi(jobDto.Number)
	if err != nil {
		job.Number = number
		jobDto.NumberErr = "Numbers only"
		return job, nil
	}
	job.Number = number

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

func PostJob(jobDto JobDto, db *sql.DB) (JobDto, error) {
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
		return jobDto, nil
	}
}
