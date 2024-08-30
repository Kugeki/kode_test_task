package pgstore

import (
	"context"
	"errors"
	"github.com/Kugeki/kode_test_task/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type UserRepoPG struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewUserRepo(db *pgxpool.Pool, log *slog.Logger) *UserRepoPG {
	return &UserRepoPG{db: db, log: log.With(slog.String("repository", "user"))}
}

func (r *UserRepoPG) CreateUser(ctx context.Context, u *domain.User) error {
	q1 := `insert into passwords
		(hash_base64, argon2_version, argon2_type, salt_base64, time, memory, threads, keylen) 
		values($1, $2, $3, $4, $5, $6, $7, $8) returning id`

	q2 := "insert into users(name, password_id) values($1, $2)"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		if err := tx.Rollback(ctx); err != nil {
			r.log.Error("transaction rollback error", slog.Any("error", err))
		}
	}(tx, ctx)

	p := u.Password

	var passwordID int
	err = r.db.QueryRow(ctx, q1, p.HashBase64, p.Argon2Version, p.Argon2Type,
		p.SaltBase64, p.Time, p.Memory, p.Threads, p.KeyLen).Scan(&passwordID)
	if err != nil {
		return DomainCreateError(err)
	}

	ct, err := r.db.Exec(ctx, q2, u.Name, passwordID)
	if err != nil {
		return DomainCreateError(err)
	}
	if ct.RowsAffected() <= 0 {
		return errors.New("user isn't inserted")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepoPG) GetUser(ctx context.Context, username string) (*domain.User, error) {
	q := `select p.hash_base64, p.argon2_version, p.argon2_type, 
       		p.salt_base64, p.time, p.memory, p.threads, p.keylen 
		from users u
			join passwords p on p.id = u.password_id 
         where u.name = $1`

	p := domain.Password{}
	err := r.db.QueryRow(ctx, q, username).
		Scan(&p.HashBase64, &p.Argon2Version, &p.Argon2Type, &p.SaltBase64, &p.Time, &p.Memory, &p.Threads, &p.KeyLen)
	if err != nil {
		return nil, DomainGetError(err)
	}

	u := &domain.User{
		Name:     username,
		Password: p,
	}

	return u, nil
}
