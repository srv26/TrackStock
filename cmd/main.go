package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/sourav/TrackStock/pkg/database/sql"
	"github.com/sourav/TrackStock/pkg/routers"
)

func main() {
	dbStorage, err := sql.SetupSqlStorage()
	if err != nil {
		panic(err)
	}
	defer dbStorage.(*sql.Storage).Db.Close()
	sm := mux.NewRouter()
	routers.RegisterTrackStockRoutes(sm, dbStorage)
	http.Handle("/", sm)

	// create a new server
	s := http.Server{
		Addr:         ":8080",           // configure the bind address
		Handler:      sm,                // set the default handler
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		log.Println("Starting server on port 8090")

		err := s.ListenAndServe()
		if err != nil {
			log.Fatalf("Error starting server error %v", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
