package infrastructure

import (
	"database/sql"
	"errors"
	"log"

	"github.com/jtalev/chat_gpg/domain/models"
)

func GetJobs(db *sql.DB) ([]domain.Job, error) {
	q := `
	select * from job order by name;
	`

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []domain.Job{}
	for rows.Next() {
		job := domain.Job{}
		if err := rows.Scan(
			&job.ID,
			&job.Name,
			&job.Number,
			&job.Address,
			&job.Suburb,
			&job.PostCode,
			&job.City,
			&job.IsComplete,
			&job.CreatedAt,
			&job.UpdatedAt,
		); err != nil {
			return nil, err
		}
		data = append(data, job)
	}
	return data, nil
}

func GetJobById(id int, db *sql.DB) (domain.Job, error) {
	q := `
	select * from job where id = ?;
	`

	rows, err := db.Query(q, id)
	if err != nil {
		return domain.Job{}, err
	}
	defer rows.Close()

	var job domain.Job
	if rows.Next() {
		if err := rows.Scan(
			&job.ID,
			&job.Name,
			&job.Number,
			&job.Address,
			&job.Suburb,
			&job.PostCode,
			&job.City,
			&job.IsComplete,
			&job.CreatedAt,
			&job.UpdatedAt,
		); err != nil {
			return job, err
		}
	} else {
		return job, errors.New("No job with provided id")
	}

	return job, nil
}

func GetJobByName(name string, db *sql.DB) (domain.Job, error) {
	q := `
	select * from job where name = ?;
	`

	rows, err := db.Query(q, name)
	if err != nil {
		return domain.Job{}, err
	}
	defer rows.Close()

	var job domain.Job
	if rows.Next() {
		if err := rows.Scan(
			&job.ID,
			&job.Name,
			&job.Number,
			&job.Address,
			&job.Suburb,
			&job.PostCode,
			&job.City,
			&job.IsComplete,
			&job.CreatedAt,
			&job.UpdatedAt,
		); err != nil {
			return job, err
		}
	} else {
		return job, errors.New("No job with provided id")
	}

	return job, nil
}

func PostJob(job domain.Job, db *sql.DB) (domain.Job, error) {
	q := `
	INSERT INTO job (name, number, address, suburb, post_code, city)
	VALUES ($1, $2, $3, $4, $5, $6);
	`

	_, err := db.Exec(q, job.Name, job.Number, job.Address, job.Suburb, job.PostCode, job.City)
	if err != nil {
		return domain.Job{}, err
	}

	newJob, err := GetJobByName(job.Name, db)
	if err != nil {
		return domain.Job{}, err
	}

	return newJob, nil
}

func PutJob(id int, job domain.Job, db *sql.DB) (domain.Job, error) {
	_, err := GetJobById(id, db)
	if err != nil {
		return domain.Job{}, err
	}

	q := `
	update job 
	set name = $1, 
	    number = $2, 
	    address = $3, 
	    suburb = $4, 
	    post_code = $5, 
	    city = $6, 
	    is_complete = $7,
	    updated_at = CURRENT_TIMESTAMP
	where id = $8;
	`

	_, err = db.Exec(q, job.Name, job.Number, job.Address, job.Suburb, job.PostCode, job.City, job.IsComplete, id)
	if err != nil {
		return domain.Job{}, err
	}
	newJob, err := GetJobByName(job.Name, db)
	if err != nil {
		return domain.Job{}, err
	}
	log.Println(newJob.IsComplete)
	return newJob, nil
}

func DeleteJob(id int, db *sql.DB) (domain.Job, error) {
	q := `
	delete from job
	where id = ?;
	`

	_, err := db.Exec(q, id)
	if err != nil {
		return domain.Job{}, err
	}

	job, err := GetJobById(id, db)
	if err == nil {
		return job, err
	}

	return job, nil
}
