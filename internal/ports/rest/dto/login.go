package dto

type LoginReq struct {
	Username string `json:"username" minLength:"1"`
	Password string `json:"password" minLength:"1"`
}

type LoginResp struct {
	AccessToken string `json:"access_token"`
}
