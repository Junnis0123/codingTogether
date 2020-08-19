package application

import (
	"bytes"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	customeMiddleware "codingtogether/application/middleware"

	response "codingtogether/application/response"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	//"github.com/labstack/echo/v4/middleware"
)

//Application is main Application.
type Application struct {
	db     *sql.DB
	e      *echo.Echo
	shaKey string
	jwtKey string
}

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), " <@>")
}

//Hashing
func (app Application) sha512Str(str string) string {

	sha := sha512.New()
	sha.Write([]byte(str))
	sha.Write([]byte(app.shaKey))

	return hex.EncodeToString(sha.Sum(nil))
}

//send mail
func (app Application) sendAuthMail(id string, mailAddr string, authKey string) {

	const templ = `이 메일은 가입 인증을 위한 메일로 답장을 보내지 마세요.
	아래 링크를 클릭하시면 가입 인증이 이루어지게 됩니다.
	https://duckbo.site:9530/auth/mail?key={{.Key}}
	
`
	t := template.New("Person template")
	t, err := t.Parse(templ)
	if err != nil {
		panic(err)
	}

	key := struct {
		Key string
	}{Key: authKey}

	var tpl bytes.Buffer
	err = t.Execute(&tpl, key)

	body := tpl.String()
	c, _ := smtp.Dial("localhost:25")
	c.Mail("noreply@duckbo.site")
	c.Rcpt(mailAddr)

	wc, _ := c.Data()
	defer wc.Close()

	msg := "From: " + "noreply@duckbo.site" + "\n" +
		"To: " + mailAddr + "\n" +
		"Subject: 모각코 가입 인증 메일입니다.\n\n" +
		body

	fmt.Fprintf(wc, msg)

	err = c.Quit()

}

//Skeleton code
func (app Application) Skeleton(c echo.Context) error {

	return c.String(http.StatusOK, "Hi")
}

//New is Application New Method
func (app Application) New(connectionInfoFileName string) {

	connectonInfo, _ := ioutil.ReadFile(connectionInfoFileName)
	shakey, _ := ioutil.ReadFile(".env.shakey")
	jwtkey, _ := ioutil.ReadFile(".env.jwtkey")
	app.shaKey = string(shakey)
	app.jwtKey = string(jwtkey)

	db, err := sql.Open("mysql", string(connectonInfo))

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app.db = db

	//ECHO 기본 설정
	app.e = echo.New()

	app.e.Use(echomiddleware.Recover())

	app.e.Use(echomiddleware.LoggerWithConfig(echomiddleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	app.e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	//static 추가
	app.e.Static("/", "static") // Echo.Static(path, root string)

	//Echo API 추가
	app.AddAPI()

	if os.Args[2] == "local" {
		app.e.Logger.Fatal(app.e.Start(":9530")) // localhost:1323
	} else {
		app.e.Logger.Fatal(app.e.StartTLS(":9530", "/etc/letsencrypt/live/duckbo.site/cert.pem", "/etc/letsencrypt/live/duckbo.site/privkey.pem")) // localhost:1323
	}

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
	auth.GET("/mail", app.authMail)
	auth.POST("/mail", app.reAuthMail)

	users.POST("/", app.createUser)

	config := customeMiddleware.JWTConfig{
		Skipper: func(c echo.Context) bool {
			if c.Path() == "/users" {
				return true
			}
			return false
		},
		//ContextKey: "data",
		SigningKey: []byte(app.jwtKey),
	}

	users.Use(customeMiddleware.JWTWithConfig(config))
	codingTogethers.Use(customeMiddleware.JWTWithConfig(config))

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

	hashUserPW := app.sha512Str(userPW)

	var userIdx int
	var ret int
	var nickName string
	var userAuth int
	rows, err := app.db.Query("SELECT count(*), user_idx,user_nickname, user_auth FROM user where user_id ='" + userID + "' and user_pw='" + hashUserPW + "'")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() //반드시 닫는다 (지연하여 닫기)

	rows.Next()
	rows.Scan(&ret, &userIdx, &nickName, &userAuth)

	if ret == 1 && userAuth != 0 {

		//AccessToken 생성
		token := jwt.New(jwt.SigningMethodHS256)
		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["userIdx"] = strconv.Itoa(userIdx)
		claims["exp"] = time.Now().Add(time.Minute * 30 * 2 * 24).Unix()

		accessToken, err := token.SignedString([]byte(app.jwtKey))
		if err != nil {
			return err
		}

		token = jwt.New(jwt.SigningMethodHS256)
		claims = token.Claims.(jwt.MapClaims)
		claims["userIdx"] = strconv.Itoa(userIdx)
		claims["realIP"] = c.RealIP()
		claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()

		refreshToken, _ := token.SignedString([]byte(app.jwtKey))

		response := response.LoginResponse{Success: true, Message: "로그인 성공", Errors: "", AccessToken: accessToken, RefreshToken: refreshToken}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusOK, json)

	} else if userAuth == 0 {
		response := response.LoginResponse{false, "이메일 미인증", "Not Auth Email", "", ""}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusUnauthorized, json)

	}
	response := response.LoginResponse{false, "로그인 실패", "Not correct ID or PW", "", ""}
	json, _ := json.Marshal(response)
	return c.JSONBlob(http.StatusUnauthorized, json)

}

//GET
func (app Application) authMail(c echo.Context) error {

	key, _ := url.QueryUnescape(c.QueryParam("key"))

	var ret int

	rows, err := app.db.Query("SELECT auth_email('" + key + "');")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() //반드시 닫는다 (지연하여 닫기)

	rows.Next()
	rows.Scan(&ret)

	if ret == 1 {

		response := response.Response{true, "가입 완료", "", ""}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusOK, json)
	}

	response := response.Response{false, "가입 실패", "Auth Failed, ", ""}
	json, _ := json.Marshal(response)
	return c.JSONBlob(http.StatusNoContent, json)

}

//POST
func (app Application) reAuthMail(c echo.Context) error {

	userID := c.FormValue("userID")
	userEmail := c.FormValue("userEmail")

	//재인증 시도

	//메일 재전송
	authKey := app.sha512Str(userID)
	app.sendAuthMail(userID, userEmail, authKey)

	response := response.Response{true, "이메일 재전송", "", ""}
	json, _ := json.Marshal(response)

	return c.JSONBlob(http.StatusOK, json)

	//or

	//새로운 키값 생성(현재는 구현 X)

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
		response := response.Response{true, "사용가능한 ID입니다.", "", ""}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusOK, json)
	}

	response := response.Response{false, "중복된 ID입니다", "Duplicate ID", ""}
	json, _ := json.Marshal(response)
	return c.JSONBlob(http.StatusConflict, json)

}

//POST
func (app Application) createUser(c echo.Context) error {

	userID := c.FormValue("userID")
	userPW := c.FormValue("userPW")
	userNickname := c.FormValue("userNickname")
	userEmail := c.FormValue("userEmail")

	hashUserPW := app.sha512Str(userPW)

	sqlStr := fmt.Sprintf("INSERT INTO user(user_id,user_pw,user_nickname, user_email) VALUES ('%s', '%s', '%s', '%s')", userID, hashUserPW, userNickname, userEmail)

	result, err := app.db.Exec(sqlStr)

	if err != nil {
		fmt.Println(err)
	}

	nRow, err := result.RowsAffected()

	if nRow == 1 {

		//여기는 이제 auth 테이블에 집어 넣기
		authKey := app.sha512Str(userID)
		sqlStr = fmt.Sprintf("INSERT INTO user_auth_key(user_auth_key_value, user_auth_key_user_id) VALUES ('%s', '%s')", authKey, userID)
		app.db.Exec(sqlStr)

		app.sendAuthMail(userID, userEmail, authKey)

		response := response.Response{true, "회원 가입 완료", "", ""}
		json, _ := json.Marshal(response)

		return c.JSONBlob(http.StatusOK, json)

	}

	response := response.Response{false, "회원 가입 실패입니다", "Not correct ID", ""}
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

		response := response.Response{true, "", "", userNickname}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusOK, json)
	}

	response := response.Response{false, "", "Not correct ID", ""}
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
	response := response.Response{true, "전체 조회 완료", "", string(datas)}
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

		response := response.Response{true, "조회 성공", "", userNickname}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusOK, json)
	}

	response := response.Response{false, "조회 실패", "Not correct user idx", ""}
	json, _ := json.Marshal(response)
	return c.JSONBlob(http.StatusNoContent, json)

}

//GET
func (app Application) showCodingTogether(c echo.Context) error {

	var codingTogetherIdx int
	var codingTogetherName string
	var codingTogetherImgURL string
	var codingTogetherCreateTime string
	var codingTogetherStartTime string
	var codingTogetherEndTime string
	var codingTogetherOrgnizerName string
	var codingTogetherUserID string
	var codingTogetherPublic int
	var codingTogetherMemberCount int
	var codingTogetherUserIdx int

	rows, err := app.db.Query("SELECT * FROM codingtogether.codingtogether_lookup_view_all where public = 1;")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() //반드시 닫는다 (지연하여 닫기)

	var Result []interface{}

	for rows.Next() {
		err := rows.Scan(&codingTogetherIdx, &codingTogetherName, &codingTogetherImgURL, &codingTogetherCreateTime, &codingTogetherStartTime, &codingTogetherEndTime, &codingTogetherOrgnizerName, &codingTogetherUserID, &codingTogetherPublic, &codingTogetherMemberCount,
			&codingTogetherUserIdx)

		data := make(map[string]interface{})

		data["codingTogetherIdx"] = codingTogetherIdx
		data["codingTogetherName"] = codingTogetherName
		data["codingTogetherImgURL"] = codingTogetherImgURL
		data["codingTogetherCreateTime"] = codingTogetherCreateTime
		data["codingTogetherOrgnizerName"] = codingTogetherOrgnizerName
		data["codingTogetherUserID"] = codingTogetherUserID
		data["codingTogetherMemberCount"] = codingTogetherMemberCount
		data["codingTogetherPublic"] = codingTogetherPublic != 0

		if err != nil {
			log.Fatal(err)
		}
		Result = append(Result, data)

	}

	datas, _ := json.Marshal(Result)

	response := response.Response{true, "조회 성공", "", string(datas)}

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
	codingTogetherStartTime := c.FormValue("codingTogetherStartTime")
	codingTogetherEndTime := c.FormValue("codingTogetherEndTime")
	var codingTogetherPublic int
	if c.FormValue("codingTogetherPublic") == "true" {
		codingTogetherPublic = 1
	} else {
		codingTogetherPublic = 0
	}

	//file
	file, _ := c.FormFile("codingTogetherImgURL")
	src, err := file.Open()
	if err != nil {
		return err
	}

	defer src.Close()

	filePath := file.Filename

	fileName := filePath[:strings.LastIndex(filePath, ".")]
	fileExtension := filePath[strings.LastIndex(filePath, "."):]

	duplicate := 0
	for {
		filePath = app.sha512Str(fileName) + strconv.Itoa(duplicate) + fileExtension
		_, err := os.Stat("static/images/" + filePath)

		if os.IsNotExist(err) {
			break
		}
		duplicate++
	}

	dst, err := os.Create("static/images/" + filePath)
	if err != nil {
		panic(err)
	}

	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		panic(err)
	}

	sqlStr := fmt.Sprintf("CALL create_codingtogether('%s','%s','%s', '%s', '%s', '%s', '%d')", userIdx, codingTogetherName, codingTogetherContents, filePath,
		codingTogetherStartTime, codingTogetherEndTime, codingTogetherPublic)

	result, err := app.db.Exec(sqlStr)

	if err != nil {
		fmt.Println(err)

	}

	nRow, err := result.RowsAffected()

	if nRow > 0 {
		response := response.Response{true, "모각코 생성 성공", "", ""}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusOK, json)
	} else {

		response := response.Response{false, "모각코 생성 실패", "CodingTogether Create Failure", ""}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusInternalServerError, json)
	}

	return c.String(http.StatusInternalServerError, "")

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
	var codingTogetherStartTime string
	var codingTogetherEndTime string
	var codingTogetherOrgnizerName string
	var codingTogetherUserID string
	var codingTogetherPublic int
	var codingTogetherMemberCount int
	var codingTogetherUserIdx int

	rows, err := app.db.Query("SELECT * FROM codingtogether.codingtogether_lookup_view_all where user_idx = '" + userIdx + "'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() //반드시 닫는다 (지연하여 닫기)

	var Result []interface{}

	for rows.Next() {
		err := rows.Scan(&codingTogetherIdx, &codingTogetherName, &codingTogetherImgURL, &codingTogetherCreateTime, &codingTogetherStartTime, &codingTogetherEndTime, &codingTogetherOrgnizerName, &codingTogetherUserID, &codingTogetherPublic, &codingTogetherMemberCount,
			&codingTogetherUserIdx)

		data := make(map[string]interface{})

		data["codingTogetherIdx"] = codingTogetherIdx
		data["codingTogetherName"] = codingTogetherName
		data["codingTogetherImgURL"] = codingTogetherImgURL
		data["codingTogetherCreateTime"] = codingTogetherCreateTime
		data["codingTogetherOrgnizerName"] = codingTogetherOrgnizerName
		data["codingTogetherUserID"] = codingTogetherUserID
		data["codingTogetherMemberCount"] = codingTogetherMemberCount
		data["codingTogetherStartTime"] = codingTogetherStartTime
		data["codingTogetherEndTime"] = codingTogetherEndTime
		data["codingTogetherPublic"] = codingTogetherPublic != 0

		if err != nil {
			log.Fatal(err)
		}
		Result = append(Result, data)

	}

	datas, _ := json.Marshal(Result)
	response := response.Response{true, "참가 모각코 조회 섣공", "", string(datas)}
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

		//loading
		rows2, err := app.db.Query("SELECT user_id,user_nickname FROM codingtogether.user_lookup_with_codingtogether where codingTogether_idx =" + codingTogetherIdx)
		if err != nil {
			log.Fatal(err)
		}
		defer rows2.Close() //반드시 닫는다 (지연하여 닫기)

		var codingTogetherUsers []interface{}

		var userID string
		var userNickname string

		for rows2.Next() {
			err := rows2.Scan(&userID, &userNickname)

			data := make(map[string]interface{})

			data["userID"] = userID
			data["userNickname"] = userNickname

			if err != nil {
				log.Fatal(err)
			}
			codingTogetherUsers = append(codingTogetherUsers, data)

		}

		data["codingTogetherContents"] = codingTogetherContents
		data["codingTogetherUsers"] = codingTogetherUsers

		datas, _ := json.Marshal(data)
		response := response.Response{true, "조회 성공", "", string(datas)}
		json, _ := json.Marshal(response)
		return c.JSONBlob(http.StatusOK, json)

	}

	response := response.Response{false, "조회 실패", "Not correct codingTogether idx", ""}
	json, _ := json.Marshal(response)
	return c.JSONBlob(http.StatusNoContent, json)

}
