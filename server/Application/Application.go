package Application

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
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

	app.e.Logger.Fatal(app.e.Start(":9530")) // localhost:1323

}
func test(c echo.Context) error {

	return c.String(http.StatusOK, "유별나 김선진")
}
func (app Application) AddAPI() {

	//추후에 Router class로 뺄 예정
	auth := app.e.Group("/auth")
	users := app.e.Group("/users")
	codingtogethers := app.e.Group("/codingtogethers")
	auth.POST("/login", app.login)

	auth.GET("/test", test)
	auth.GET("/duplication/:user_id", app.checkDuplication)

	users.POST("/", app.createUser)

	config := middleware.JWTConfig{
		Skipper: func(c echo.Context) bool {
			if c.Path() == "/users" {
				return true
			}
			return false
		},
		//ContextKey: "data",
		SigningKey: []byte("soolfam"),
	}
	users.Use(middleware.JWTWithConfig(config))
	codingtogethers.Use(middleware.JWTWithConfig(config))

	users.GET("/", app.showUsersAll)
	users.GET("/test/:user_id", app.showUser)
	users.GET("/me", app.showMySelf)

	//coding together
	codingtogethers.GET("/", app.showCodingTogether)
	codingtogethers.POST("/", app.createCodingTogether)
}

//POST
func (app Application) login(c echo.Context) error {

	user_id := c.FormValue("user_id")
	user_pw := c.FormValue("user_pw")

	fmt.Println(user_id, " ", user_pw)
	var user_idx int
	var ret int
	var nick_name string
	rows, err := app.db.Query("SELECT count(*), user_idx,user_nickname FROM user where user_id ='" + user_id + "' and user_pw='" + user_pw + "'")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() //반드시 닫는다 (지연하여 닫기)

	rows.Next()
	rows.Scan(&ret, &user_idx, &nick_name)

	if ret == 1 {

		//AccessToken 생성
		token := jwt.New(jwt.SigningMethodHS256)
		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["user_idx"] = strconv.Itoa(user_idx)
		fmt.Println(claims["user_idx"])
		claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

		accessToken, err := token.SignedString([]byte("soolfam"))
		if err != nil {
			return err
		}

		token = jwt.New(jwt.SigningMethodHS256)
		claims = token.Claims.(jwt.MapClaims)
		claims["user_idx"] = strconv.Itoa(user_idx)
		claims["real_ip"] = c.RealIP()
		claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()

		refreshToken, _ := token.SignedString([]byte("soolfam"))

		response := LoginResponse{true, "Login Successed", "", accessToken, refreshToken}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusOK, json)
	} else {

		response := LoginResponse{false, "Login Failure", "Not correct ID or PW", "", ""}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusUnauthorized, json)
	}

}

//GET
func (app Application) checkDuplication(c echo.Context) error {

	user_id, _ := url.QueryUnescape(c.Param("user_id"))
	rows, err := app.db.Query("SELECT count(*) FROM user where user_id ='" + user_id + "'")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() //반드시 닫는다 (지연하여 닫기)

	var ret int
	rows.Next()
	rows.Scan(&ret)

	if ret == 0 {
		response := Response{true, "Not Duplicate ID", "", ""}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusOK, json)
	} else {

		response := Response{false, "", "Duplicate ID", ""}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusConflict, json)
	}

}

//POST
func (app Application) createUser(c echo.Context) error {

	user_id := c.FormValue("user_id")
	user_pw := c.FormValue("user_pw")
	user_nickname := c.FormValue("user_nickname")

	sqlStr := fmt.Sprintf("INSERT INTO user(user_id,user_pw,user_nickname) VALUES ('%s', '%s', '%s')", user_id, user_pw, user_nickname)

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
		return c.JSONBlob(http.StatusInternalServerError, json)
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

//GET
func (app Application) showMySelf(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	user_idx := claims["user_idx"].(string)

	fmt.Println("type  = ", reflect.TypeOf(user_idx))
	fmt.Println("user idx = ", user_idx)
	var ret int
	var nick_name string
	rows, err := app.db.Query("SELECT count(*), user_nickname FROM user where user_idx ='" + user_idx + "'")

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
func (app Application) showCodingTogether(c echo.Context) error {

	var codingTogether_idx int
	var codingTogether_name string
	var codingTogether_img_url sql.NullString
	var codingTogether_create_time string
	var codingTogether_Orgnizer_name string
	var codingTogether_member_count int

	rows, err := app.db.Query("SELECT * FROM codingtogether.codingtogether_lookup_view;")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() //반드시 닫는다 (지연하여 닫기)

	var Result []interface{}

	for rows.Next() {
		err := rows.Scan(&codingTogether_idx, &codingTogether_name, &codingTogether_img_url, &codingTogether_create_time, &codingTogether_Orgnizer_name, &codingTogether_member_count)

		data := make(map[string]interface{})

		data["codingTogether_idx"] = codingTogether_idx
		data["codingTogether_name"] = codingTogether_name
		data["codingTogether_img_url"] = codingTogether_img_url
		data["codingTogether_create_time"] = codingTogether_create_time
		data["codingTogether_Orgnizer_name"] = codingTogether_Orgnizer_name
		data["codingTogether_member_count"] = codingTogether_member_count

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

//POST
func (app Application) createCodingTogether(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	user_idx := claims["user_idx"].(string)
	codingtogether_name := c.FormValue("codingtogether_name")
	codingtogether_contents := c.FormValue("codingtogether_contents")

	sqlStr := fmt.Sprintf("CALL create_codingtogether('%s','%s','%s')", user_idx, codingtogether_name, codingtogether_contents)

	result, err := app.db.Exec(sqlStr)

	if err != nil {
		fmt.Println(err)
	}

	nRow, err := result.RowsAffected()

	if nRow > 0 {
		response := Response{true, "Create Success", "", ""}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusOK, json)
	} else {

		response := Response{false, "Create Failure", "Create Failure", ""}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusInternalServerError, json)
	}

}
