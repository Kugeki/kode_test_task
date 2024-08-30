package pgstore

import (
	"context"
	"github.com/Kugeki/kode_test_task/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NoteRepoPG struct {
	db *pgxpool.Pool
}

func NewNoteRepo(db *pgxpool.Pool) *NoteRepoPG {
	return &NoteRepoPG{db: db}
}

func (r *NoteRepoPG) CreateNote(ctx context.Context, authorName string, n *domain.Note) error {
	q := "insert into notes(author_name, content) values($1, $2) returning id"

	err := r.db.QueryRow(ctx, q, authorName, n.Content).Scan(&n.ID)
	if err != nil {
		return DomainCreateError(err)
	}

	return nil
}

func (r *NoteRepoPG) GetNotesForUser(ctx context.Context, username string) ([]*domain.Note, error) {
	q := "select n.id, n.content from notes n where n.author_name = $1"

	rows, err := r.db.Query(ctx, q, username)
	if err != nil {
		return nil, err
	}

	notes, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*domain.Note, error) {
		n := &domain.Note{}
		err := row.Scan(&n.ID, &n.Content)
		return n, err
	})
	if err != nil {
		return nil, err
	}

	return notes, nil
}
