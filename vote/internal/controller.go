package internal

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/uptrace/bunrouter"
	"net/http"
)

type VoteController interface {
	VoteHandler(w http.ResponseWriter, req bunrouter.Request) error
}

type voteController struct {
	redis RedisClient
}

func NewVoteController(redis RedisClient) VoteController {
	return &voteController{redis: redis}
}

type Vote struct {
	VoterId string `json:"voter_id"`
	Vote    string `json:"vote"`
}

func (v Vote) MarshalBinary() ([]byte, error) {
	return json.Marshal(v)
}

func (v voteController) VoteHandler(w http.ResponseWriter, req bunrouter.Request) error {
	voterId := uuid.New().String()
	vote := req.Params().ByName("vote")

	voteData := Vote{
		VoterId: voterId,
		Vote:    vote,
	}

	err := v.redis.RPush("votes", voteData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return bunrouter.JSON(w, err.Error())
	}
	w.WriteHeader(http.StatusCreated)
	return bunrouter.JSON(w, voteData)
}
