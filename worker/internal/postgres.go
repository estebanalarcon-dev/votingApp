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
	SaveVote(vote *VoteData) error
}

type postgresDb struct {
	db *bun.DB
	//conn   *bun.Conn
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

func (p postgresDb) SaveVote(vote *VoteData) error {
	_, err := p.db.NewInsert().Model(vote).Exec(context.Background())
	return err
}
