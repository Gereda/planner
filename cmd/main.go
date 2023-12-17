package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"planner/db"
	"planner/delivery/endpoints"
	"planner/delivery/rest"
)

func main() {
	db := db.InitDB()

	service := rest.NewService(db)

	defer db.Close()

	r := gin.Default()
	endpoints.EndPoints(r, service)

	r.Run("127.0.0.1:8000")
}
