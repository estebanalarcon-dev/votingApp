package main

import (
	"github.com/uptrace/bunrouter"
	"log"
	"net/http"
	"os"
	"os/signal"
	"result/internal"
	"syscall"
)

func main() {
	router := buildIn()
	start(router)
}

func start(router http.Handler) {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
}

func buildIn() http.Handler {
	db := internal.NewDBInstance()
	voteController := internal.NewResultController(db)
	return NewRouter(voteController)
}

func NewRouter(resultController internal.ResultController) http.Handler {
	router := bunrouter.New()
	router.GET("/result", resultController.GetResultHandler)
	return router
}
