package jobnotes

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type JobnotesRepo interface {
	GetNotesByJobId(jobId int) ([]Note, error)
	PostNote(note Note) error
	PutNote(note Note) error
	DeleteNote(uuid string) error
}

type Jobnotes struct {
	NoteUuid string
	JobId    int

	Paintnote paintnote
	Tasknote  tasknote
	Imagenote imagenote

	Paintnotes []paintnote
	Tasknotes  []tasknote
	Imagenotes []imagenote

	NoteErrors interface{}

	Repo JobnotesRepo
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

func (j *Jobnotes) marshalNote(noteType string) (string, error) {
	note := ""
	switch noteType {
	case "paint_note":
		n, err := json.Marshal(j.Paintnote)
		if err != nil {
			return note, err
		}
		note = string(n)
	case "task_note":
		n, err := json.Marshal(j.Tasknote)
		if err != nil {
			return note, err
		}
		note = string(n)
	case "image_note":
		n, err := json.Marshal(j.Imagenote)
		if err != nil {
			return note, err
		}
		note = string(n)
	default:
		return "", fmt.Errorf("note type %s not supported", noteType)
	}
	return note, nil
}

func (j *Jobnotes) PostNote(noteType string) error {
	isSuccess := j.validateNote(noteType)
	if !isSuccess {
		// simply return, j.NoteErrors will be present
		return nil
	} else {
		n, err := j.marshalNote(noteType)
		if err != nil {
			return err
		}

		uuid := uuid.NewString()
		note := Note{
			Uuid:     uuid,
			JobId:    j.JobId,
			NoteType: noteType,
			Note:     n,
		}
		err = j.Repo.PostNote(note)
		if err != nil {
			log.Printf("error posting note %v: %v", note, err)
			return err
		}
		return nil
	}
}

func (j *Jobnotes) getNoteUuid(noteType string) (string, error) {
	switch noteType {
	case "paint_note":
		return j.Paintnote.NoteUuid, nil
	case "task_note":
		return j.Tasknote.NoteUuid, nil
	case "image_note":
		return j.Imagenote.NoteUuid, nil
	default:
		return "", fmt.Errorf("note type %s not supported", noteType)
	}
}

func (j *Jobnotes) PutNote(noteType string) error {
	isSuccess := j.validateNote(noteType)
	if !isSuccess {
		// simply return, j.NoteErrors will be present
		return nil
	} else {
		n, err := j.marshalNote(noteType)
		if err != nil {
			return err
		}

		noteUuid, err := j.getNoteUuid(noteType)
		if err != nil {
			return err
		}
		note := Note{
			Uuid: noteUuid,
			Note: n,
		}
		err = j.Repo.PutNote(note)
		if err != nil {
			log.Printf("error posting note %v: %v", note, err)
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
