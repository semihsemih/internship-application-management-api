package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/semihsemih/internship-application-management-api/db"
	"github.com/semihsemih/internship-application-management-api/router"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func init() {

	// Connect DB
	db.Init()
	client := db.GetClient()
	err := client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	}

}

func main() {

	// Router
	r := router.Init()

	fmt.Println("Server started at localhost:8000")
	http.ListenAndServe(":8000", r)
}
