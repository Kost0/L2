package main

import (
	"L2/18/handlers"
	"L2/18/middleware"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	port := flag.String("port", "8080", "port to listen on")

	flag.Parse()

	envPort := os.Getenv("PORT")

	if envPort != "" {
		port = &envPort
	}

	address := fmt.Sprintf(":%s", *port)

	r := chi.NewRouter()

	srv := &http.Server{
		Addr:    address,
		Handler: r,
	}

	r.Use(middleware.Middleware)

	r.Post("/delete_event", handlers.DeleteEvent)

	r.Post("/update_event", handlers.UpdateEvent)

	r.Post("/create_event", handlers.CreateEvent)

	r.Get("/events_for_day", handlers.GetEventForDay)

	r.Get("/events_for_week", handlers.GetEventForWeek)

	r.Get("/events_for_month", handlers.GetEventForMonth)

	log.Println("Starting HTTP server...")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Error starting HTTP server: %v\n", err)
	}
}
