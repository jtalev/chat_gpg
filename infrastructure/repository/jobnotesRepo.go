package infrastructure

import (
	"database/sql"

	"github.com/jtalev/chat_gpg/application/services/jobnotes"
)

type JobnotesRepo struct {
	Db *sql.DB
}

func (j *JobnotesRepo) GetNoteByUuid(uuid string) (jobnotes.Note, error) {
	q := `
	select * from note where uuid = ?
	`
	rows, err := j.Db.Query(q, uuid)
	if err != nil {
		return jobnotes.Note{}, err
	}
	defer rows.Close()

	n := jobnotes.Note{}
	for rows.Next() {
		if err := rows.Scan(
			&n.Uuid,
			&n.JobId,
			&n.NoteType,
			&n.Note,
			&n.IsArchived,
			&n.CreatedAt,
			&n.ModifiedAt,
		); err != nil {
			return n, err
		}
	}
	return n, nil
}

func (j *JobnotesRepo) GetNotesByJobId(jobId int) ([]jobnotes.Note, error) {
	q := `
	select * from note where job_id = ? and is_archived = 0;
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
			&n.IsArchived,
			&n.CreatedAt,
			&n.ModifiedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, nil
}

func (j *JobnotesRepo) PostNote(note jobnotes.Note) error {
	q := `
	INSERT INTO note (uuid, job_id, note_type, note)
	VALUES ($1, $2, $3, $4);
	`

	_, err := j.Db.Exec(q, note.Uuid, note.JobId, note.NoteType, note.Note)
	if err != nil {
		return err
	}

	return nil
}

func (j *JobnotesRepo) PutNote(note jobnotes.Note) error {
	q := `
	UPDATE note
	SET note = $1
	WHERE uuid = $2;
	`

	_, err := j.Db.Exec(q, note.Note, note.Uuid)
	if err != nil {
		return err
	}
	return nil
}

func (j *JobnotesRepo) DeleteNote(uuid string) error {
	q := `
	DELETE FROM note WHERE uuid = ?;
	`
	_, err := j.Db.Exec(q, uuid)
	if err != nil {
		return err
	}
	return nil
}

func (j *JobnotesRepo) ArchiveNote(uuid string, isArchived bool) error {
	q := `
	UPDATE note
	SET is_archived = $1
	WHERE uuid = $2;
	`

	isArchivedInt := 0
	if isArchived {
		isArchivedInt = 1
	}

	_, err := j.Db.Exec(q, isArchivedInt, uuid)
	if err != nil {
		return err
	}
	return nil
}
