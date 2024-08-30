package dto

import "github.com/Kugeki/kode_test_task/internal/domain"

type CreateNoteReq struct {
	Content string `json:"content"`
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
