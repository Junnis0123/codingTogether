package Application

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Application struct {
	db *sql.DB
	e  *echo.Echo
}

//Skeleton ode
func (app Application) Skeletion(c echo.Context) error {

	return c.String(http.StatusOK, "Hi")
}

func (app Application) New(connectionInfoFileName string) {

	connectonInfo, _ := ioutil.ReadFile(connectionInfoFileName)

	db, err := sql.Open("mysql", string(connectonInfo))

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app.db = db

	//ECHO 기본 설정
	app.e = echo.New()

	app.e.Use(middleware.Recover())

	app.e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	app.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	//Echo API 추가
	app.AddAPI()

	app.e.Logger.Fatal(app.e.Start(":80")) // localhost:1323

}
func test(c echo.Context) error {

	return c.String(http.StatusOK, "유별나 김선진")
}
func (app Application) AddAPI() {

	//추후에 Router class로 뺄 예정

	auth := app.e.Group("/auth")
	users := app.e.Group("/users")

	auth.POST("/login", app.login)

	auth.GET("/test", test)

	users.POST("/", app.createUser)
	users.GET("/", app.showUsersAll)
	users.GET("/:user_id", app.showUser)

}

//POST
func (app Application) login(c echo.Context) error {

	user_id := c.FormValue("user_id")
	user_pw := c.FormValue("user_pw")

	fmt.Println(user_id)
	fmt.Println(user_pw)

	var ret int
	var nick_name string
	rows, err := app.db.Query("SELECT count(*), user_nickname FROM user where user_id ='" + user_id + "' and user_pw='" + user_pw + "'")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() //반드시 닫는다 (지연하여 닫기)

	rows.Next()
	rows.Scan(&ret, &nick_name)

	if ret == 1 {

		response := Response{true, "김선진 칭찬해", "", ""}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusOK, json)
	} else {

		response := Response{false, "엥 로그인 실패여", "Not correct ID or PW", ""}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusUnauthorized, json)
	}

}

//POST
func (app Application) createUser(c echo.Context) error {

	user_id := c.FormValue("user_id")
	user_pw := c.FormValue("user_pw")
	user_nickname := c.FormValue("user_nickname")

	sqlStr := fmt.Sprintf("INSERT INTO user(user_id,user_pw,user_nickname) VALUES ('%s', '%s', '%s')", user_id, user_pw, user_nickname)

	fmt.Println(sqlStr)
	result, err := app.db.Exec(sqlStr)

	if err != nil {
		fmt.Println(err)
	}

	nRow, err := result.RowsAffected()

	if nRow == 1 {
		response := Response{true, "회원 가입 완료", "", ""}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusOK, json)
	} else {

		response := Response{false, "회원 가입 실패입니다", "Not correct ID", ""}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusUnauthorized, json)
	}

}

//GET
func (app Application) showUser(c echo.Context) error {

	user_id, _ := url.QueryUnescape(c.Param("user_id"))

	var ret int
	var nick_name string
	rows, err := app.db.Query("SELECT count(*), user_nickname FROM user where user_id ='" + user_id + "'")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() //반드시 닫는다 (지연하여 닫기)

	rows.Next()
	rows.Scan(&ret, &nick_name)

	if ret == 1 {

		response := Response{true, "", "", nick_name}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusOK, json)
	} else {

		response := Response{false, "", "Not correct ID", ""}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusNoContent, json)
	}

}

//GET
func (app Application) showUsersAll(c echo.Context) error {

	var user_id string
	var user_nickname string

	rows, err := app.db.Query("SELECT user_id, user_nickname FROM user")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() //반드시 닫는다 (지연하여 닫기)

	var Result []interface{}

	for rows.Next() {
		err := rows.Scan(&user_id, &user_nickname)

		data := make(map[string]interface{})

		data["user_id"] = user_id
		data["user_nickname"] = user_nickname

		if err != nil {
			log.Fatal(err)
		}
		Result = append(Result, data)

	}

	datas, _ := json.Marshal(Result)
	response := Response{true, "", "", string(datas)}
	json, _ := json.Marshal(response)

	return c.JSONBlob(http.StatusOK, json)
}
