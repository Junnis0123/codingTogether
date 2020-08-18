package main

import (
	Application "codingtogether/application"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	App := Application.Application{}
	App.New(os.Args[1])
}

/*
func main() {

	//DB loading
	const connectQuery = (`duckbo:@Testuser2995@tcp(139.150.64.36:3306)/codingtogether`)
	db, _ = sql.Open("mysql", connectQuery)

	defer db.Close()

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	auth := e.Group("/auth")

	auth.POST("/login", login)
	auth.GET("/test", test)

	e.Logger.Fatal(e.Start(":80"))
}
*/
