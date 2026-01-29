package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sambeetpanda507/advance-search/controllers"
	"github.com/sambeetpanda507/advance-search/middlewares"
	"github.com/sambeetpanda507/advance-search/routers"
	"github.com/sambeetpanda507/advance-search/utils"
)

func main() {
	secrets := utils.GetSecrets()
	if len(secrets.PORT) == 0 {
		log.Fatal("PORT is missing")
	}

	port := fmt.Sprintf(":%s", secrets.PORT)
	mux := http.NewServeMux()

	// Connect to db
	db := utils.Connect()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance: ", err)
	}

	// Set up dependancies
	newsController := controllers.NewsController{DB: db}

	// Set up routes
	routers.Ping(mux)
	routers.NewsRoutes(mux, newsController)

	go func() {
		fmt.Printf("Starting server in %s\n", port)
		if err := http.ListenAndServe(port, middlewares.Logger(mux)); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := sqlDB.Close(); err != nil {
		log.Fatal("Falided to close db connection: ", err)
	}

	log.Println("Shutting down...")
	log.Println("Good Bye!!!")
}
