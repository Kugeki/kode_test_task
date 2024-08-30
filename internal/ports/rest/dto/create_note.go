package dto

import "github.com/Kugeki/kode_test_task/internal/domain"

type CreateNoteReq struct {
	Content string `json:"content" minLength:"1"`
}

func (r *CreateNoteReq) ToDomain() *domain.Note {
	return &domain.Note{Content: r.Content}
}

type CreateNoteResp struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

func (r *CreateNoteResp) FromDomain(note *domain.Note) {
	r.ID = note.ID
	r.Content = note.Content
}

type SpellError struct {
	Code int      `json:"code"`
	Pos  int      `json:"pos"`
	Row  int      `json:"row"`
	Col  int      `json:"col"`
	Len  int      `json:"len"`
	Word string   `json:"word"`
	S    []string `json:"s"`
}

type NoteSpellErrorResp struct {
	NoteContent string `json:"note_content"`
	Error       string `json:"error"`

	SpellErrors []SpellError `json:"spell_errors"`
}

var SpellErrorStr = "note: spell error"

func (r *NoteSpellErrorResp) FromDomain(note *domain.Note, results domain.SpellResults) {
	r.NoteContent = note.Content
	r.Error = SpellErrorStr

	r.SpellErrors = make([]SpellError, 0, len(results))
	for _, v := range results {
		spellError := SpellError{
			Code: v.Code,
			Pos:  v.Pos,
			Row:  v.Row,
			Col:  v.Col,
			Len:  v.Len,
			Word: v.Word,
			S:    v.S,
		}
		r.SpellErrors = append(r.SpellErrors, spellError)
	}
}
