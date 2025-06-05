package infrastructure

import (
	"database/sql"

	"github.com/jtalev/chat_gpg/application/services/jobnotes"
)

type JobnotesRepo struct {
	Db *sql.DB
}

func (j *JobnotesRepo) GetJobnotesByJobId(jobId int) ([]jobnotes.Note, error) {
	q := `
	select * from note where job_id = ?
	`
	rows, err := j.Db.Query(q, jobId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []jobnotes.Note{}
	n := jobnotes.Note{}
	for rows.Next() {
		if err := rows.Scan(
			&n.Uuid,
			&n.JobId,
			&n.NoteType,
			&n.Note,
			&n.CreatedAt,
			&n.ModifiedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, nil
}
