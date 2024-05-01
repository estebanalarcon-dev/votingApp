package internal

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"time"
)

type Database interface {
	GetResultGroupByVote() ([]ResultByVote, error)
}

const queryGetResultGroupByVote = "SELECT vote, COUNT(id) AS count FROM votes GROUP BY vote"

type postgresDb struct {
	db *bun.DB
	//conn   *bun.Conn
}

type ResultByVote struct {
	Vote  string
	Count int64
}

func NewDBInstance() Database {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"db", 5432, "postgres", "postgres", "postgres")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//basic retrying connection
	for i := 0; i < 3; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return &postgresDb{db: bun.NewDB(db, pgdialect.New())}
}

func (p postgresDb) GetResultGroupByVote() ([]ResultByVote, error) {
	var results []ResultByVote
	err := p.db.NewRaw(queryGetResultGroupByVote).Scan(context.Background(), &results)
	return results, err
}
