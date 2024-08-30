package tests

import (
	"fmt"
	"github.com/gavv/httpexpect/v2"
	"net/http"
)

type User struct {
	Name     string
	Password string
}

// эти пользователи должны присутствовать на сервере (т.к. регистрация в задании не обязательна)
func Get4TestUsers() []User {
	return []User{
		{Name: "test1", Password: "password1"},
		{Name: "test2", Password: "password2"},
		{Name: "test3", Password: "password3"},
		{Name: "test4", Password: "password4"},
	}
}

func (su *APITestSuite) Login(username, password string) *httpexpect.Response {
	loginReq := map[string]interface{}{
		"username": username,
		"password": password,
	}

	return su.e.POST("/users/login/").
		WithJSON(loginReq).
		Expect()
}

func (su *APITestSuite) GetToken(username, password string) string {
	return su.Login(username, password).
		Status(http.StatusOK).
		JSON().Object().ContainsKey("access_token").
		Value("access_token").String().Raw()
}

func (su *APITestSuite) CreateNote(content string, token string) *httpexpect.Response {
	createNoteReq := map[string]interface{}{
		"content": content,
	}

	return su.e.POST("/notes/create/").
		WithJSON(createNoteReq).
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect()
}

func (su *APITestSuite) GetNotes(token string) *httpexpect.Response {
	return su.e.GET("/notes/").
		WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect()
}
