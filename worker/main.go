package main

import (
	"worker/internal"
)

func main() {
	chMessages := make(chan string, 50)
	worker := buildIn()
	go worker.ConsumeVotesList(chMessages)
	worker.PersistVotes(chMessages)
}

func buildIn() internal.Worker {
	redis := internal.NewRedisClient()
	db := internal.NewDBInstance()
	worker := internal.NewWorker(redis, db)
	return worker
}
