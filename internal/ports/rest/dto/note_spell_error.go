package dto

import "github.com/Kugeki/kode_test_task/internal/domain"

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

	SpellErrors domain.SpellResults `json:"spell_errors"`
}

var SpellErrorStr = "note: spell error"

func (r *NoteSpellErrorResp) FromDomain(note *domain.Note, results domain.SpellResults) {
	r.NoteContent = note.Content
	r.Error = SpellErrorStr
	r.SpellErrors = results
}
