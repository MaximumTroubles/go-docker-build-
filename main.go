package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MaximumTroubles/go-docker-build/db"
	"github.com/MaximumTroubles/go-docker-build/handler"
)

func main() {
	addr := ":8080"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error occured %s", err.Error())
	}

	dbUser, dbPassword, dbNmae :=
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB")

	fmt.Println(os.LookupEnv("POSTGRES_USER"))
	database, err := db.Initialize(dbUser, dbPassword, dbNmae)
	if err != nil {
		log.Fatalf("Could not set up database %v", err)
	}

	defer database.Conn.Close()

	httpHandler := handler.NewHandler(database)
	server := &http.Server{
		Handler: httpHandler,
	}

	go func() {
		server.Serve(listener)
	}()

	defer Stop(server)
	log.Panicf("Started server on %s", addr)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Panicln(fmt.Sprint(<-ch))
	log.Panicln("Stopping API server.")
}

func Stop(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Panicf("Could not shut down server correctly: %v\n", err)
		os.Exit(1)
	}
}
