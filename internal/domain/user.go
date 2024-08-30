package domain

type User struct {
	Name     string
	Password Password
}

type Password struct {
	HashBase64 string

	Argon2Version int
	Argon2Type    Argon2Type
	SaltBase64    string
	Time          uint32
	Memory        uint32
	Threads       uint8
	KeyLen        uint32
}

type Argon2Type int

// https://datatracker.ietf.org/doc/html/rfc9106#name-argon2-inputs-and-outputs
const (
	Argon2dType Argon2Type = iota
	Argon2iType
	Argon2idType
)
