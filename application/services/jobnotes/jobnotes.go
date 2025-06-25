package jobnotes

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/google/uuid"
	models "github.com/jtalev/chat_gpg/domain/models"
	"github.com/jtalev/chat_gpg/infrastructure/img"
)

type NoteRepo interface {
	GetNotesByJobId(jobId int) ([]Note, error)
	PostNote(note Note) error
	PutNote(note Note) error
	DeleteNote(uuid string) error
}

type Note struct {
	Uuid       string `json:"uuid"`
	JobId      int    `json:"job_id"`
	NoteType   string `json:"note_type"` // eg. 'paintnote', 'tasknote', 'imagenote'
	Note       string `json:"note"`      // going to be serialized json, deserialized to note type
	CreatedAt  string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
}

type Paintnote struct {
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

type Tasknote struct {
	NoteUuid    string `json:"note_uuid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	Notes       string `json:"notes"`
}

type Imagenote struct {
	NoteUuid    string `json:"note_uuid"`
	S3uuid      string `json:"s3uuid"`
	ImageBase64 string `json:"image_base64"`
	Caption     string `json:"caption"`
	Area        string `json:"area"`
	Notes       string `json:"notes"`
	S3Url       template.HTML
}

type jobnoteSummary struct {
	Id             int
	Name           string
	Address        string
	PaintnoteCount int
	TasknoteCount  int
	ImagenoteCount int
}

type jobnoteViewData struct {
	JobCount  int
	Summaries []jobnoteSummary
}

type jobSummary struct {
	ID             int
	Name           string
	Address        string
	PaintnoteCount int
	TasknoteCount  int
	ImagenoteCount int
}

type PaintnoteFormData struct {
	FormType  string
	JobId     int
	Paintnote Paintnote
	Errors    paintnoteerrors
}

type TasknoteFormData struct {
	FormType string
	JobId    int
	Tasknote Tasknote
	Errors   tasknoteerrors
}

type ImagenoteFormData struct {
	FormType  string
	JobId     int
	Imagenote Imagenote
	Errors    imagenoteerrors
}

type Jobnotes struct {
	JobnoteViewData   jobnoteViewData
	PaintnoteFormData PaintnoteFormData
	TasknoteFormData  TasknoteFormData
	ImagenoteFormData ImagenoteFormData

	JobId      int
	JobSummary jobSummary

	Paintnote Paintnote
	Tasknote  Tasknote
	Imagenote Imagenote

	Paintnotes []Paintnote
	Tasknotes  []Tasknote
	Imagenotes []Imagenote

	Note Note

	NoteErrors interface{}

	Repo NoteRepo
}

func (j *Jobnotes) InitialJobnoteViewData(jobs []models.Job) error {
	var data jobnoteViewData
	data.JobCount = len(jobs)

	for _, job := range jobs {
		err := j.GetJobNotes(job.ID)
		if err != nil {
			return err
		}
		data.Summaries = append(data.Summaries, jobnoteSummary{
			Id:             job.ID,
			Name:           job.Name,
			Address:        fmt.Sprintf("%d %s", job.Number, job.Address),
			PaintnoteCount: len(j.Paintnotes),
			TasknoteCount:  len(j.Tasknotes),
			ImagenoteCount: len(j.Imagenotes),
		})
	}

	j.JobnoteViewData = data

	return nil
}

func unmarshalPaintnote(n Note) (Paintnote, error) {
	var p Paintnote
	err := json.Unmarshal([]byte(n.Note), &p)
	if err != nil {
		return p, err
	}
	return p, nil
}

func unmarshalTasknote(n Note) (Tasknote, error) {
	var t Tasknote
	err := json.Unmarshal([]byte(n.Note), &t)
	if err != nil {
		return t, err
	}
	return t, nil
}

func unmarshalImagenote(n Note) (Imagenote, error) {
	var i Imagenote
	err := json.Unmarshal([]byte(n.Note), &i)
	if err != nil {
		return i, err
	}
	return i, nil
}

func (j *Jobnotes) unmarshalNotes(jobnotes []Note) {
	j.Paintnotes = []Paintnote{}
	j.Tasknotes = []Tasknote{}
	j.Imagenotes = []Imagenote{}

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
		case Paintnote:
			j.Paintnotes = append(j.Paintnotes, n)
		case Tasknote:
			j.Tasknotes = append(j.Tasknotes, n)
		case Imagenote:
			j.Imagenotes = append(j.Imagenotes, n)
		default:
			log.Printf("note type %s not supported", jn.NoteType)
		}
	}
}

func getImgUrl(imgNote *Imagenote, imgStore *img.ImgStore) error {
	url, err := imgStore.GetImgUrl(imgNote.S3uuid, "paint-note")
	if err != nil {
		return err
	}

	imgNote.S3Url = template.HTML(url)
	return nil
}

func (j *Jobnotes) GetJobNotes(jobId int) error {
	notes, err := j.Repo.GetNotesByJobId(jobId)
	if err != nil {
		log.Printf("error getting notes for job %v: %v", jobId, err)
		return err
	}

	j.unmarshalNotes(notes)

	imgStore := img.InitImgStore()
	for i := range j.Imagenotes {
		err = getImgUrl(&j.Imagenotes[i], &imgStore)
		if err != nil {
			return err
		}
	}

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
		switch errors := j.NoteErrors.(type) {
		case paintnoteerrors:
			j.PaintnoteFormData.Errors = errors
		case tasknoteerrors:
			j.TasknoteFormData.Errors = errors
		case imagenoteerrors:
			j.ImagenoteFormData.Errors = errors
		default:
			log.Println("note type %s note supported", noteType)
			return false
		}
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

func writeImg(imgBase64 string) error {
	decoded, err := base64.StdEncoding.DecodeString(imgBase64)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("../ui/static/jobnotes/temp_img.jpg", decoded, 0644)
	if err != nil {
		return err
	}

	return nil
}

func storeImg(imgBase64, uuid string) error {
	err := writeImg(imgBase64)
	if err != nil {
		return err
	}

	imgStore := img.InitImgStore()
	path := filepath.Join("..", "ui", "static", "jobnotes", "temp_img.jpg")
	err = imgStore.Store(path, uuid, "paint-note")
	if err != nil {
		return err
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

		if noteType == "image_note" {
			j.Imagenote.S3uuid = uuid
			err := storeImg(j.Imagenote.ImageBase64, uuid)
			if err != nil {
				return err
			}
		}

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
