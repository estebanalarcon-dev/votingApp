package main

import (
	"github.com/uptrace/bunrouter"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"vote/internal"
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

func NewRouter(voteController internal.VoteController) http.Handler {
	router := bunrouter.New()
	router.POST("/vote/:vote", voteController.VoteHandler)
	router.GET("/hello", func(w http.ResponseWriter, req bunrouter.Request) error {
		return bunrouter.JSON(w, "hello")
	})
	return router
}

func buildIn() http.Handler {
	redis := internal.NewRedisClient()
	voteController := internal.NewVoteController(redis)
	return NewRouter(voteController)
}
