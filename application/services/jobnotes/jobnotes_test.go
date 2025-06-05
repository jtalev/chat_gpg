package jobnotes

import (
	"reflect"
	"testing"
)

func TestDecodePaintnote(t *testing.T) {
	tests := []struct {
		note     Note
		expected interface{}
	}{
		{
			Note{
				Note: `{"note_uuid":"123abc","brand":"Haymes","product":"Expressions","colour":"Natural White","finish":"Low Sheen","area":"Living Room","coats":2,"surfaces":"Walls","notes":""}`,
			},
			paintnote{
				NoteUuid: "123abc",
				Brand:    "Haymes",
				Product:  "Expressions",
				Colour:   "Natural White",
				Finish:   "Low Sheen",
				Area:     "Living Room",
				Coats:    2,
				Surfaces: "Walls",
				Notes:    "",
			},
		},
	}

	for _, tt := range tests {
		p, err := decodePaintnote(tt.note)
		if err != nil {
			t.Fatal(err)
		}

		if p != tt.expected {
			t.Errorf("error decoding valid paintnote. expected = %v, got = %v", tt.expected, p)
		}
	}
}

// func TestDecodeTasknote(t *testing.T) {
// 	tests := []struct {
// 		note     Note
// 		expected interface{}
// 	}{}
// }

// func TestDecodeImagenote(t *testing.T) {
// 	tests := []struct {
// 		note     Note
// 		expected interface{}
// 	}{}
// }

func TestDecodeJobnotes(t *testing.T) {
	tests := []struct {
		jobnotes []Note
		expected Jobnotes
	}{
		{
			[]Note{
				{
					NoteType: "paint_note",
					Note:     `{"note_uuid":"123abc","brand":"Haymes","product":"Expressions","colour":"Natural White","finish":"Low Sheen","area":"Living Room","coats":2,"surfaces":"Walls","notes":""}`,
				},
				{
					NoteType: "task_note",
					Note:     `{"note_uuid":"123abc","title":"Paint Walls","description":"Paint all walls","status":"complete","priority":"low","notes":""}`,
				},
				{
					NoteType: "image_note",
					Note:     `{"note_uuid":"123abc","s3uuid":"123abc","caption":"hole","area":"Living room","notes":""}`,
				},
			},
			Jobnotes{
				Paintnotes: []paintnote{
					{
						NoteUuid: "123abc",
						Brand:    "Haymes",
						Product:  "Expressions",
						Colour:   "Natural White",
						Finish:   "Low Sheen",
						Area:     "Living Room",
						Coats:    2,
						Surfaces: "Walls",
					},
				},
				Tasknotes: []tasknote{
					{
						NoteUuid:    "123abc",
						Title:       "Paint Walls",
						Description: "Paint all walls",
						Status:      "complete",
						Priority:    "low",
					},
				},
				Imagenotes: []imagenote{
					{
						NoteUuid: "123abc",
						S3uuid:   "123abc",
						Caption:  "hole",
						Area:     "Living room",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		j := Jobnotes{}
		j.decodeJobNotes(tt.jobnotes)
		if !reflect.DeepEqual(j.Paintnotes, tt.expected.Paintnotes) ||
			!reflect.DeepEqual(j.Tasknotes, tt.expected.Tasknotes) ||
			!reflect.DeepEqual(j.Imagenotes, tt.expected.Imagenotes) {
			t.Errorf("decodeJobNotes() = %v, expected %v", j, tt.expected)
		}
	}
}
