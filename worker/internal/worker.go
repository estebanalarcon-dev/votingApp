package internal

import (
	"encoding/json"
	"fmt"
	"github.com/uptrace/bun"
	"log"
	"time"
)

type Worker interface {
	ConsumeVotesList(ch chan string) error
	PersistVotes(ch chan string)
}

type worker struct {
	redis RedisClient
	db    Database
}

func NewWorker(redis RedisClient, db Database) Worker {
	return worker{redis: redis, db: db}
}

type VoteData struct {
	bun.BaseModel `bun:"table:votes"`
	VoterId       string `json:"voter_id" bun:"id,pk,notnull"`
	Vote          string `json:"vote" bun:"vote,notnull"`
}

func (w worker) ConsumeVotesList(ch chan string) error {
	for {
		//only query each 1 ms
		time.Sleep(time.Millisecond)
		data, err := w.redis.RPop("votes")
		if err == nil {
			fmt.Println("data:", data)
			ch <- data
		} else if err != nil && err.Error() != "redis: nil" {
			fmt.Println(err)
		}
	}
}

func (w worker) PersistVotes(ch chan string) {
	//convert json into vote struct
	for dataJson := range ch {
		var vote VoteData
		json.Unmarshal([]byte(dataJson), &vote)
		fmt.Println(vote)
		err := w.db.SaveVote(&vote)
		if err != nil {
			log.Println(err)
		}
	}
}
