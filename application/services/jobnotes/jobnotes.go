package jobnotes

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type NoteRepo interface {
	GetNotesByJobId(jobId int) ([]Note, error)
	PostNote(note Note) error
	PutNote(note Note) error
	DeleteNote(uuid string) error
}

type Jobnotes struct {
	JobId int

	Paintnote paintnote
	Tasknote  tasknote
	Imagenote imagenote

	Paintnotes []paintnote
	Tasknotes  []tasknote
	Imagenotes []imagenote

	Note       Note
	NoteErrors interface{}

	Repo NoteRepo
}

type Note struct {
	Uuid       string `json:"uuid"`
	JobId      int    `json:"job_id"`
	NoteType   string `json:"note_type"` // eg. 'paintnote', 'tasknote', 'imagenote'
	Note       string `json:"note"`      // going to be serialized json, deserialized to note type
	CreatedAt  string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
}

type paintnote struct {
	NoteUuid string `json:"note_uuid"`
	Brand    string `json:"brand"`
	Product  string `json:"product"`
	Colour   string `json:"colour"`
	Finish   string `json:"finish"`
	Area     string `json:"area"`
	Coats    int    `json:"coats"`
	Surfaces string `json:"surfaces"`
	Notes    string `json:"notes"`
}

type tasknote struct {
	NoteUuid    string `json:"note_uuid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	Notes       string `json:"notes"`
}

type imagenote struct {
	NoteUuid string `json:"note_uuid"`
	S3uuid   string `json:"s3uuid"`
	Caption  string `json:"caption"`
	Area     string `json:"area"`
	Notes    string `json:"notes"`
}

func unmarshalPaintnote(n Note) (paintnote, error) {
	var p paintnote
	err := json.Unmarshal([]byte(n.Note), &p)
	if err != nil {
		return p, err
	}
	return p, nil
}

func unmarshalTasknote(n Note) (tasknote, error) {
	var t tasknote
	err := json.Unmarshal([]byte(n.Note), &t)
	if err != nil {
		return t, err
	}
	return t, nil
}

func unmarshalImagenote(n Note) (imagenote, error) {
	var i imagenote
	err := json.Unmarshal([]byte(n.Note), &i)
	if err != nil {
		return i, err
	}
	return i, nil
}

func (j *Jobnotes) unmarshalNotes(jobnotes []Note) {
	type unmarshalerFunc func(Note) (interface{}, error)
	unmarshalers := map[string]unmarshalerFunc{
		"paint_note": func(n Note) (interface{}, error) { return unmarshalPaintnote(n) },
		"task_note":  func(n Note) (interface{}, error) { return unmarshalTasknote(n) },
		"image_note": func(n Note) (interface{}, error) { return unmarshalImagenote(n) },
	}

	for _, jn := range jobnotes {
		decoder, ok := unmarshalers[jn.NoteType]
		if !ok {
			log.Printf("note type %s not supported", jn.NoteType)
			continue
		}

		note, err := decoder(jn)
		if err != nil {
			log.Printf("error decoding %s: %v", jn.NoteType, err)
			continue
		}

		switch n := note.(type) {
		case paintnote:
			j.Paintnotes = append(j.Paintnotes, n)
		case tasknote:
			j.Tasknotes = append(j.Tasknotes, n)
		case imagenote:
			j.Imagenotes = append(j.Imagenotes, n)
		default:
			log.Printf("note type %s not supported", jn.NoteType)
		}
	}
}

func (j *Jobnotes) GetJobNotes(jobId int) error {
	notes, err := j.Repo.GetNotesByJobId(jobId)
	if err != nil {
		log.Printf("error getting notes for job %v: %v", jobId, err)
		return err
	}

	j.unmarshalNotes(notes)

	return nil
}

func (j *Jobnotes) validateNote(noteType string) bool {
	validatorFuncs := map[string]validatorFunc{
		"paint_note": func() (interface{}, bool) { return j.Paintnote.validate() },
		"task_note":  func() (interface{}, bool) { return j.Tasknote.validate() },
		"image_note": func() (interface{}, bool) { return j.Imagenote.validate() },
	}

	var isSuccess bool
	if validator, ok := validatorFuncs[noteType]; ok {
		j.NoteErrors, isSuccess = validator()
	}

	return isSuccess
}

func (j *Jobnotes) marshalNote(noteType, uuid string) error {
	switch noteType {
	case "paint_note":
		j.Paintnote.NoteUuid = uuid
		n, err := json.Marshal(j.Paintnote)
		if err != nil {
			return err
		}
		j.Note.Note = string(n)
	case "task_note":
		j.Tasknote.NoteUuid = uuid
		n, err := json.Marshal(j.Tasknote)
		if err != nil {
			return err
		}
		j.Note.Note = string(n)
	case "image_note":
		j.Imagenote.NoteUuid = uuid
		n, err := json.Marshal(j.Imagenote)
		if err != nil {
			return err
		}
		j.Note.Note = string(n)
	default:
		return fmt.Errorf("note type %s not supported", noteType)
	}
	return nil
}

func (j *Jobnotes) PostNote(noteType string) error {
	isSuccess := j.validateNote(noteType)
	if !isSuccess {
		// simply return, j.NoteErrors will be present
		return nil
	} else {
		uuid := uuid.NewString()
		err := j.marshalNote(noteType, uuid)
		if err != nil {
			return err
		}

		j.Note.Uuid = uuid

		err = j.Repo.PostNote(j.Note)
		if err != nil {
			log.Printf("error posting note %v: %v", j.Note, err)
			return err
		}
		return nil
	}
}

func (j *Jobnotes) PutNote(noteType string) error {
	isSuccess := j.validateNote(noteType)
	if !isSuccess {
		// simply return, j.NoteErrors will be present
		return nil
	} else {
		err := j.marshalNote(noteType, j.Note.Uuid)
		if err != nil {
			return err
		}

		err = j.Repo.PutNote(j.Note)
		if err != nil {
			log.Printf("error posting note %v: %v", j.Note, err)
			return err
		}
		return nil
	}
}

func (j *Jobnotes) DeleteNote(uuid string) error {
	err := j.Repo.DeleteNote(uuid)
	if err != nil {
		return err
	}
	return nil
}
