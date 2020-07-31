package application

/*
"success": {true} or {false} //로그인 여부
    "message": null,
    "errors": null,
	"data": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1OThkZGI2MzIyYWMxMDExZTA3MDJjYjAiLCJ1c2VybmFtZSI6InRlc3QxIiwibmFtZSI6InRlc3QxIiwiZW1haWwiOiIiLCJpYXQiOjE1MDQ3MzI2NzcsImV4cCI6MTUwNDgxOTA3N30.4eG2zGpSeY2XezKB4Djf6usy7DdygIybR1VKUBj-ScE"
*/

//Response struct
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Errors  string `json:"errors"`
	Data    string `json:"data"`
}

//LoginResponse struct
type LoginResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	Errors       string `json:"errors"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
