package main

import (
	"context"
	"github.com/MykhailoKondrat/go-micro/handlers"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	l := log.New(os.Stdout, "propduct-api", log.LstdFlags)
	//hh := handlers.NewHello(l)
	//gb := handlers.NewGoodBuy(l)
	ph := handlers.NewProducts(l)
	sm := mux.NewRouter()
	ops := middleware.RedocOpts{
		SpecURL: "./swagger.yaml",
	}
	sh := middleware.Redoc(ops, nil)

	getRouter := sm.Methods("GET").Subrouter()
	getRouter.HandleFunc("/products", ph.GetProducts)
	getRouter.Handle("/docs", sh)

	putRouter := sm.Methods("PUT").Subrouter()

	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)
	putRouter.Use(ph.MiddlewareProductValidation)

	postRouter := sm.Methods("POST").Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)

	deleteRouter := sm.Methods("DELETE").Subrouter()

	deleteRouter.HandleFunc("/products/{id:[0-9]+}", ph.DeleteProduct)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	go func() {
		err := s.ListenAndServe()

		if err != nil {
			l.Fatal(err)
		}
	}()
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
