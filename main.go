package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/CoryEvans2324/SystemsDesignAppAPI/database"
	"github.com/CoryEvans2324/SystemsDesignAppAPI/models"
	"github.com/CoryEvans2324/SystemsDesignAppAPI/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", routes.Index)

	tracksRouter := r.PathPrefix("/tracks").PathPrefix("/tracks").Subrouter()
	tracksRouter.HandleFunc("/upload", routes.UploadTracks).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	database.CreateDatabase(fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"127.0.0.1",
		5432,
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		"postgres",
	))

	database.DB.AutoMigrate(models.Track{})

	// Graceful shutdown from https://github.com/gorilla/mux#graceful-shutdown

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	log.Printf("Server is running at %s", srv.Addr)

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("shutting down")
	srv.Shutdown(ctx)

	os.Exit(0)
}
