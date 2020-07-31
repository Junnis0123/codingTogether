package application

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//Application is main Application.
type Application struct {
	db *sql.DB
	e  *echo.Echo
}

//Skeleton code
func (app Application) Skeleton(c echo.Context) error {

	return c.String(http.StatusOK, "Hi")
}

//New is Application New Method
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

//AddAPI is api add method
//to-do refactoring
func (app Application) AddAPI() {

	//추후에 Router class로 뺄 예정
	auth := app.e.Group("/auth")
	users := app.e.Group("/users")
	codingTogethers := app.e.Group("/codingTogethers")
	auth.POST("/login", app.login)

	auth.GET("/test", test)
	auth.GET("/duplication/:userID", app.checkDuplication)

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
	codingTogethers.Use(middleware.JWTWithConfig(config))

	users.GET("/", app.showUsersAll)
	users.GET("/test/:userID", app.showUser)
	users.GET("/me", app.showMySelf)

	//coding together
	codingTogethers.GET("/", app.showCodingTogether)
	codingTogethers.POST("/", app.createCodingTogether)
	codingTogethers.GET("/me", app.showCodingTogetherMySelf)
	codingTogethers.GET("/:codingTogetherIdx", app.showCodingTogetherContents)
}

//login
func (app Application) login(c echo.Context) error {

	userID := c.FormValue("userID")
	userPW := c.FormValue("userPW")

	var userIdx int
	var ret int
	var nickName string
	rows, err := app.db.Query("SELECT count(*), user_idx,user_nickname FROM user where user_id ='" + userID + "' and user_pw='" + userPW + "'")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() //반드시 닫는다 (지연하여 닫기)

	rows.Next()
	rows.Scan(&ret, &userIdx, &nickName)

	if ret == 1 {

		//AccessToken 생성
		token := jwt.New(jwt.SigningMethodHS256)
		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["userIdx"] = strconv.Itoa(userIdx)
		claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

		accessToken, err := token.SignedString([]byte("soolfam"))
		if err != nil {
			return err
		}

		token = jwt.New(jwt.SigningMethodHS256)
		claims = token.Claims.(jwt.MapClaims)
		claims["userIdx"] = strconv.Itoa(userIdx)
		claims["realIP"] = c.RealIP()
		claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()

		refreshToken, _ := token.SignedString([]byte("soolfam"))

		response := LoginResponse{true, "로그인 성공", "", accessToken, refreshToken}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusOK, json)

	}
	response := LoginResponse{false, "로그인 실패", "Not correct ID or PW", "", ""}
	json, _ := json.Marshal(response)
	return c.JSONBlob(http.StatusUnauthorized, json)

}

//GET
func (app Application) checkDuplication(c echo.Context) error {

	userID, _ := url.QueryUnescape(c.Param("userID"))
	rows, err := app.db.Query("SELECT count(*) FROM user where user_id ='" + userID + "'")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() //반드시 닫는다 (지연하여 닫기)

	var ret int
	rows.Next()
	rows.Scan(&ret)

	if ret == 0 {
		response := Response{true, "사용가능한 ID입니다.", "", ""}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusOK, json)
	}

	response := Response{false, "중복된 ID입니다", "Duplicate ID", ""}
	json, _ := json.Marshal(response)
	return c.JSONBlob(http.StatusConflict, json)

}

//POST
func (app Application) createUser(c echo.Context) error {

	userID := c.FormValue("userID")
	userPW := c.FormValue("userPW")
	userNickname := c.FormValue("userNickname")

	sqlStr := fmt.Sprintf("INSERT INTO user(user_id,user_pw,user_nickname) VALUES ('%s', '%s', '%s')", userID, userPW, userNickname)

	result, err := app.db.Exec(sqlStr)

	if err != nil {
		fmt.Println(err)
	}

	nRow, err := result.RowsAffected()

	if nRow == 1 {
		response := Response{true, "회원 가입 완료", "", ""}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusOK, json)
	}

	response := Response{false, "회원 가입 실패입니다", "Not correct ID", ""}
	json, _ := json.Marshal(response)
	return c.JSONBlob(http.StatusInternalServerError, json)

}

//GET
func (app Application) showUser(c echo.Context) error {

	userID, _ := url.QueryUnescape(c.Param("userID"))

	var ret int
	var userNickname string
	rows, err := app.db.Query("SELECT count(*), user_nickname FROM user where user_id ='" + userID + "'")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() //반드시 닫는다 (지연하여 닫기)

	rows.Next()
	rows.Scan(&ret, &userNickname)

	if ret == 1 {

		response := Response{true, "", "", userNickname}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusOK, json)
	}

	response := Response{false, "", "Not correct ID", ""}
	json, _ := json.Marshal(response)
	return c.JSONBlob(http.StatusNoContent, json)

}

//GET
func (app Application) showUsersAll(c echo.Context) error {

	var userID string
	var userNickname string

	rows, err := app.db.Query("SELECT user_id, user_nickname FROM user")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() //반드시 닫는다 (지연하여 닫기)

	var Result []interface{}

	for rows.Next() {

		err := rows.Scan(&userID, &userNickname)

		data := make(map[string]interface{})

		data["userID"] = userID
		data["userNickname"] = userNickname

		if err != nil {
			log.Fatal(err)
		}
		Result = append(Result, data)

	}

	datas, _ := json.Marshal(Result)
	response := Response{true, "전체 조회 완료", "", string(datas)}
	json, _ := json.Marshal(response)

	return c.JSONBlob(http.StatusOK, json)
}

//GET
func (app Application) showMySelf(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIdx := claims["userIdx"].(string)

	var ret int
	var userNickname string
	rows, err := app.db.Query("SELECT count(*), user_nickname FROM user where user_idx ='" + userIdx + "'")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() //반드시 닫는다 (지연하여 닫기)

	rows.Next()
	rows.Scan(&ret, &userNickname)

	if ret == 1 {

		response := Response{true, "조회 성공", "", userNickname}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusOK, json)
	}

	response := Response{false, "조회 실패", "Not correct user idx", ""}
	json, _ := json.Marshal(response)
	return c.JSONBlob(http.StatusNoContent, json)

}

//GET
func (app Application) showCodingTogether(c echo.Context) error {

	var codingTogetherIdx int
	var codingTogetherName string
	var codingTogetherImgURL string
	var codingTogetherCreateTime string
	var codingTogetherOrgnizerName string
	var codingTogetherMemberCount int

	rows, err := app.db.Query("SELECT * FROM codingtogether.codingtogether_lookup_view;")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() //반드시 닫는다 (지연하여 닫기)

	var Result []interface{}

	for rows.Next() {
		err := rows.Scan(&codingTogetherIdx, &codingTogetherName, &codingTogetherImgURL, &codingTogetherCreateTime, &codingTogetherOrgnizerName, &codingTogetherMemberCount)

		data := make(map[string]interface{})

		data["codingTogetherIdx"] = codingTogetherIdx
		data["codingTogetherName"] = codingTogetherName
		data["codingTogetherImgURL"] = codingTogetherImgURL
		data["codingTogetherCreateTime"] = codingTogetherCreateTime
		data["codingTogetherOrgnizerName"] = codingTogetherOrgnizerName
		data["codingTogetherMemberCount"] = codingTogetherMemberCount

		if err != nil {
			log.Fatal(err)
		}
		Result = append(Result, data)

	}

	datas, _ := json.Marshal(Result)

	response := Response{true, "조회 성공", "", string(datas)}

	json, _ := json.Marshal(response)

	return c.JSONBlob(http.StatusOK, json)
}

//POST
func (app Application) createCodingTogether(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	userIdx := claims["userIdx"].(string)
	codingTogetherName := c.FormValue("codingTogetherName")
	codingTogetherContents := c.FormValue("codingTogetherContents")

	sqlStr := fmt.Sprintf("CALL create_codingtogether('%s','%s','%s')", userIdx, codingTogetherName, codingTogetherContents)

	result, err := app.db.Exec(sqlStr)

	if err != nil {
		fmt.Println(err)
	}

	nRow, err := result.RowsAffected()

	if nRow > 0 {
		response := Response{true, "모각코 생성 성공", "", ""}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusOK, json)
	} else {

		response := Response{false, "모각코 생성 실패", "CodingTogether Create Failure", ""}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusInternalServerError, json)
	}

}

//GET
func (app Application) showCodingTogetherMySelf(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIdx := claims["userIdx"].(string)

	var codingTogetherIdx int
	var codingTogetherName string
	var codingTogetherImgURL string
	var codingTogetherCreateTime string
	var codingTogetherOrgnizerName string
	var codingTogetherMemberCount int
	var codingTogetherUserIdx int

	rows, err := app.db.Query("SELECT * FROM codingtogether.codingtogether_lookup_view_all where user_idx = '" + userIdx + "'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() //반드시 닫는다 (지연하여 닫기)

	var Result []interface{}

	for rows.Next() {
		err := rows.Scan(&codingTogetherIdx, &codingTogetherName, &codingTogetherImgURL, &codingTogetherCreateTime, &codingTogetherOrgnizerName, &codingTogetherMemberCount,
			&codingTogetherUserIdx)

		data := make(map[string]interface{})

		data["codingTogetherIdx"] = codingTogetherIdx
		data["codingTogetherName"] = codingTogetherName
		data["codingTogetherImgURL"] = codingTogetherImgURL
		data["codingTogetherCreateTime"] = codingTogetherCreateTime
		data["codingTogetherOrgnizerName"] = codingTogetherOrgnizerName
		data["codingTogetherMemberCount"] = codingTogetherMemberCount

		if err != nil {
			log.Fatal(err)
		}
		Result = append(Result, data)

	}

	datas, _ := json.Marshal(Result)
	response := Response{true, "참가 모각코 조회 섣공", "", string(datas)}
	json, _ := json.Marshal(response)

	return c.JSONBlob(http.StatusOK, json)

}

func (app Application) showCodingTogetherContents(c echo.Context) error {

	codingTogetherIdx, _ := url.QueryUnescape(c.Param("codingTogetherIdx"))
	fmt.Println(codingTogetherIdx)
	var codingTogetherContents string

	rows, err := app.db.Query("SELECT codingtogether_contents FROM codingtogether where codingTogether_idx =" + codingTogetherIdx)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() //반드시 닫는다 (지연하여 닫기)

	isSuccess := rows.Next()

	if isSuccess {

		rows.Scan(&codingTogetherContents)

		data := make(map[string]interface{})

		data["codingTogetherContents"] = codingTogetherContents

		datas, _ := json.Marshal(data)
		response := Response{true, "조회 성공", "", string(datas)}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusOK, json)

	}

	response := Response{false, "조회 실패", "Not correct codingTogether idx", ""}
	json, _ := json.Marshal(response)
	return c.JSONBlob(http.StatusNoContent, json)
	/*
		row := app.db.QueryRow("SELECT codingtogether_contents FROM codingtogether where codingTogether_idx =" + codingTogetherIdx)

		err := row.Scan(&codingTogetherContents)

		fmt.Println(codingTogetherContents)
		if err != nil {

			response := Response{false, "조회 실패", "Not correct user idx", ""}
			json, _ := json.Marshal(response)
			return c.JSONBlob(http.StatusNoContent, json)

		}

		response := Response{true, "조회 성공", "", codingTogetherContents}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusOK, json)
	*/
}
