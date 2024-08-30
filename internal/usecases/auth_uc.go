package usecases

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/Kugeki/kode_test_task/internal/domain"
	"golang.org/x/crypto/argon2"
	"reflect"
)

type UserRepository interface {
	CreateUser(ctx context.Context, u *domain.User) error
	GetUser(ctx context.Context, username string) (*domain.User, error)
}

type AuthUC struct {
	userRepo UserRepository
}

func NewAuthUC(userRepo UserRepository) *AuthUC {
	return &AuthUC{userRepo: userRepo}
}

// SaltSize recommended is 16: https://datatracker.ietf.org/doc/html/rfc9106#name-argon2-inputs-and-outputs
var SaltSize = 16

const (
	Argon2Time    = 1
	Argon2Memory  = 128
	Argon2Threads = 1
	Argon2KeyLen  = 32
)

func (uc *AuthUC) CreateUser(ctx context.Context, u *domain.User, password string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	salt := make([]byte, SaltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}

	pwHash := argon2.IDKey([]byte(password), salt, Argon2Time, Argon2Memory, Argon2Threads, Argon2KeyLen)

	pwHashBase64 := base64.StdEncoding.EncodeToString(pwHash)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)

	u.Password = domain.Password{
		HashBase64:    pwHashBase64,
		Argon2Version: argon2.Version,
		Argon2Type:    domain.Argon2idType,
		SaltBase64:    saltBase64,
		Time:          Argon2Time,
		Memory:        Argon2Memory,
		Threads:       Argon2Threads,
		KeyLen:        Argon2KeyLen,
	}

	err = uc.userRepo.CreateUser(ctx, u)
	if err != nil {
		return err
	}

	return nil
}

func (uc *AuthUC) VerifyLogin(ctx context.Context, username string, password string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	u, err := uc.userRepo.GetUser(ctx, username)
	if errors.Is(err, domain.ErrNotFound) {
		return ErrWrongPassword // не даём чекать на существование пользователя
	}
	if err != nil {
		return err
	}

	pw := u.Password

	wantHash, err := base64.StdEncoding.DecodeString(pw.HashBase64)
	if err != nil {
		return err
	}

	salt, err := base64.StdEncoding.DecodeString(pw.SaltBase64)
	if err != nil {
		return err
	}

	gotHash := argon2.IDKey([]byte(password), salt, pw.Time, pw.Memory, pw.Threads, pw.KeyLen)

	if !reflect.DeepEqual(wantHash, gotHash) {
		return ErrWrongPassword
	}

	return nil
}
