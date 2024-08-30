package tests

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"
	"net/http"
	"os"
	"testing"
)

func TestAPI(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

type APITestSuite struct {
	suite.Suite
	serverURL string

	e      *httpexpect.Expect
	users  []User
	tokens []string
}

func (su *APITestSuite) SetupSuite() {
	su.serverURL = os.Getenv("SERVER_URL")
	if su.serverURL == "" {
		su.serverURL = "http://localhost:8080"
	}
	su.Require().NotEmpty(su.serverURL)

	su.e = httpexpect.Default(su.T(), su.serverURL)

	su.users = Get4TestUsers()
	su.Require().Len(su.users, 4)

	su.tokens = make([]string, 0, 4)
	for _, u := range su.users {
		su.tokens = append(su.tokens, su.GetToken(u.Name, u.Password))
	}
}

func (su *APITestSuite) TestLogin() {
	u := su.users[0] // valid user

	su.Run("successful", func() {
		obj := su.Login(u.Name, u.Password).
			Status(http.StatusOK).JSON().
			Object()

		obj.Keys().ContainsOnly("access_token")
		obj.Value("access_token").String().NotEmpty()
	})

	errorCases := []struct {
		Name     string
		Username string
		Password string
		Status   int
	}{
		{
			Name:     "not exist",
			Username: "not_exists_1231231231", Password: "password",
			Status: http.StatusUnauthorized,
		},
		{
			Name:     "wrong password",
			Username: "test1", Password: "wrong_password",
			Status: http.StatusUnauthorized,
		},
		{
			Name:     "empty username",
			Username: "", Password: "password",
			Status: http.StatusBadRequest,
		},
		{
			Name:     "empty password",
			Username: "username", Password: "",
			Status: http.StatusBadRequest,
		},
	}

	for _, tc := range errorCases {
		su.Run(tc.Name, func() {
			su.Login(tc.Username, tc.Password).
				Status(tc.Status).JSON().
				Object().Keys().ContainsOnly("error")
		})
	}
}

func (su *APITestSuite) TestCreateNote() {
	validToken := su.tokens[0]

	su.Run("successful", func() {
		su.CreateNote("правильное правописание", validToken).
			Status(http.StatusCreated).JSON().
			Object().Keys().ContainsOnly("id", "content")
	})

	su.Run("wrong spelling", func() {
		su.CreateNote("превет", validToken).
			Status(http.StatusBadRequest).JSON().
			Object().Keys().ContainsOnly("note_content", "error", "spell_errors")
	})

	errorCases := []struct {
		Name    string
		Content string
		Token   string
		Status  int
	}{
		{
			Name:    "empty content",
			Content: "", Token: validToken,
			Status: http.StatusBadRequest,
		},
		{
			Name:    "invalid token",
			Content: "some content", Token: "123",
			Status: http.StatusBadRequest,
		},
	}
	for _, tc := range errorCases {
		su.Run(tc.Name, func() {
			su.CreateNote(tc.Content, tc.Token).
				Status(tc.Status).JSON().
				Object().Keys().ContainsOnly("error")
		})
	}
}

func (su *APITestSuite) TestGetNotes() {
	su.Run("invalid token", func() {
		su.GetNotes("123").
			Status(http.StatusBadRequest).JSON().
			Object().Keys().ContainsOnly("error")
	})

	token1 := su.tokens[1] // token that has user1
	token2 := su.tokens[2] // token that has user2

	newCounter := func(token string) *int {
		obj := su.GetNotes(token).
			Status(http.StatusOK).
			JSON().Object()

		obj.Keys().ContainsOnly("notes")
		counter := int(obj.Value("notes").Array().Length().Raw())
		return &counter
	}

	bufSize := 16
	buf := make([]byte, bufSize)

	n, err := rand.Read(buf[:])
	su.Require().NoError(err)
	su.Require().Equal(bufSize, n)

	noteSuffix := base64.URLEncoding.EncodeToString(buf[:]) // для того, чтобы note не повторялись

	counter1 := newCounter(token1)
	counter2 := newCounter(token2)

	testCases := []struct {
		Token   string
		Counter *int
		Note    string

		IsError bool
	}{
		{Token: token1, Counter: counter1, Note: "note1", IsError: false},
		{Token: token2, Counter: counter2, Note: "note2", IsError: false},
		{Token: token1, Counter: counter1, Note: "превет", IsError: true},
	}

	for _, tc := range testCases {
		tc.Note = fmt.Sprintf("[TESTGEN] %v %v", tc.Note, noteSuffix)

		su.CreateNote(tc.Note, tc.Token).Raw()

		obj := su.GetNotes(tc.Token).
			Status(http.StatusOK).
			JSON().Object()

		obj.Keys().ContainsOnly("notes")
		arr := obj.Value("notes").Array()

		findFunc := func(index int, value *httpexpect.Value) bool {
			return tc.Note == value.Object().Value("content").String().Raw()
		}

		if tc.IsError {
			arr.NotFind(findFunc)
		} else {
			arr.Find(findFunc)
			*tc.Counter = *tc.Counter + 1
		}

		arr.Length().IsEqual(*tc.Counter)
	}
}
