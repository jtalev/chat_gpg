package jobnotes

import (
	"encoding/json"
	"log"
)

type JobnotesRepo interface {
	GetJobnotesByJobId(jobId int) ([]Note, error)
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

func decodePaintnote(n Note) (paintnote, error) {
	var p paintnote
	err := json.Unmarshal([]byte(n.Note), &p)
	if err != nil {
		return p, err
	}
	return p, nil
}

func decodeTasknote(n Note) (tasknote, error) {
	var t tasknote
	err := json.Unmarshal([]byte(n.Note), &t)
	if err != nil {
		return t, err
	}
	return t, nil
}

func decodeImagenote(n Note) (imagenote, error) {
	var i imagenote
	err := json.Unmarshal([]byte(n.Note), &i)
	if err != nil {
		return i, err
	}
	return i, nil
}

func (j *Jobnotes) decodeJobNotes(jobnotes []Note) {
	type decoderFunc func(Note) (interface{}, error)
	decoders := map[string]decoderFunc{
		"paint_note": func(n Note) (interface{}, error) { return decodePaintnote(n) },
		"task_note":  func(n Note) (interface{}, error) { return decodeTasknote(n) },
		"image_note": func(n Note) (interface{}, error) { return decodeImagenote(n) },
	}

	for _, jn := range jobnotes {
		decoder, ok := decoders[jn.NoteType]
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
	notes, err := j.Repo.GetJobnotesByJobId(jobId)
	if err != nil {
		log.Printf("error getting notes for job %v: %v", jobId, err)
		return err
	}

	j.decodeJobNotes(notes)

	return nil
}
