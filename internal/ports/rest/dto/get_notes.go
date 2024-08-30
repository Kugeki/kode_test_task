package dto

import "github.com/Kugeki/kode_test_task/internal/domain"

type NoteDto struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

type GetNotesResp struct {
	Notes []NoteDto `json:"notes"`
}

func (r *GetNotesResp) FromDomain(notes []*domain.Note) {
	r.Notes = make([]NoteDto, 0, len(notes))
	for _, n := range notes {
		r.Notes = append(r.Notes, NoteDto{
			ID:      n.ID,
			Content: n.Content,
		})
	}
}
