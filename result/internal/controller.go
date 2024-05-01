package internal

import (
	"github.com/uptrace/bunrouter"
	"net/http"
)

type ResultController interface {
	GetResultHandler(w http.ResponseWriter, req bunrouter.Request) error
}

type resultController struct {
	db Database
}

func NewResultController(db Database) ResultController {
	return &resultController{db: db}
}

func (r resultController) GetResultHandler(w http.ResponseWriter, req bunrouter.Request) error {
	results, err := r.db.GetResultGroupByVote()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return bunrouter.JSON(w, err.Error())
	}

	if len(results) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return bunrouter.JSON(w, results)
	}

	w.WriteHeader(http.StatusOK)
	return bunrouter.JSON(w, results)
}
