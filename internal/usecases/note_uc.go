package usecases

import (
	"context"
	"github.com/Kugeki/kode_test_task/internal/domain"
)

type SpellingClient interface {
	CheckText(ctx context.Context, text string) (domain.SpellResults, error)
}

type NoteRepository interface {
	CreateNote(ctx context.Context, authorName string, n *domain.Note) error
	GetNotesForUser(ctx context.Context, username string) ([]*domain.Note, error)
}

type NoteUC struct {
	spellClient SpellingClient
	noteRepo    NoteRepository
}

func NewNoteUC(spellClient SpellingClient, noteRepo NoteRepository) *NoteUC {
	return &NoteUC{spellClient: spellClient, noteRepo: noteRepo}
}

func (uc *NoteUC) ValidateNoteSpelling(ctx context.Context, n *domain.Note) (domain.SpellResults, error) {
	return uc.spellClient.CheckText(ctx, n.Content)
}

func (uc *NoteUC) CreateNote(ctx context.Context, authorName string, n *domain.Note) error {
	return uc.noteRepo.CreateNote(ctx, authorName, n)
}

func (uc *NoteUC) GetNotesForUser(ctx context.Context, username string) ([]*domain.Note, error) {
	return uc.noteRepo.GetNotesForUser(ctx, username)
}
