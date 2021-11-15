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
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", routes.Index)

	tracksRouter := r.PathPrefix("/tracks").Subrouter()
	tracksRouter.HandleFunc("/", routes.GetTracks).Methods("GET")
	tracksRouter.HandleFunc("/upload", routes.UploadTracks).Methods("POST")

	headersOk := handlers.AllowedHeaders([]string{"*"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	srv := &http.Server{
		Handler:      handlers.CORS(originsOk, headersOk, methodsOk)(r),
		Addr:         "0.0.0.0:80",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	databaseURL := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	// log.Println(databaseURL)
	database.CreateDatabase(databaseURL)

	database.DB.AutoMigrate(models.Track{})
	database.DB.Exec(`
	CREATE OR REPLACE FUNCTION calculate_distance(lat1 float, lon1 float, lat2 float, lon2 float, units varchar)
	RETURNS float AS $dist$
		DECLARE
			dist float = 0;
			radlat1 float;
			radlat2 float;
			theta float;
			radtheta float;
		BEGIN
			IF lat1 = lat2 OR lon1 = lon2
				THEN RETURN dist;
			ELSE
				radlat1 = pi() * lat1 / 180;
				radlat2 = pi() * lat2 / 180;
				theta = lon1 - lon2;
				radtheta = pi() * theta / 180;
				dist = sin(radlat1) * sin(radlat2) + cos(radlat1) * cos(radlat2) * cos(radtheta);
	
				IF dist > 1 THEN dist = 1; END IF;
	
				dist = acos(dist);
				dist = dist * 180 / pi();
				dist = dist * 60 * 1.1515;
	
				IF units = 'K' THEN dist = dist * 1.609344; END IF;
				IF units = 'N' THEN dist = dist * 0.8684; END IF;
	
				RETURN dist;
			END IF;
		END;
	$dist$ LANGUAGE plpgsql;
`)

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
