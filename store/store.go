package store

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Store struct {
	DbReader *sqlx.DB
	DbWriter *sqlx.DB
	log      *logrus.Logger
}

func New(log *logrus.Logger, dsnReader, dsnWriter string) (*Store, error) {
	var s Store
	var dbReader *sqlx.DB
	var dbWriter *sqlx.DB
	var errReader error
	var errWriter error
	for i := 0; i < 30; i++ {
		dbReader, errReader = sqlx.Connect("pgx", dsnReader)
		if errReader != nil {
			return &s, errReader
		}

		dbWriter, errWriter = sqlx.Connect("pgx", dsnWriter)
		if errWriter != nil {
			return &s, errWriter
		}

		if errReader == nil && errWriter == nil {
			break
		}
		time.Sleep(time.Second)
		log.Info("retrying to connect to database")
	}

	// 4 instances in prod with ~1365 open connection capacity on RDS dbReader.r5.large
	// setting to 300 means 1200 max open connections to RDS so we should never
	// hit too many open connections error
	dbReader.SetMaxOpenConns(300)

	s.DbReader = dbReader
	s.DbWriter = dbWriter
	s.log = log

	if errReader != nil {
		return &s, errReader
	}

	if errWriter != nil {
		return &s, errWriter
	}
	return &s, nil
}

func (s *Store) HealthCheck(ctx context.Context) error {
	return s.DbReader.PingContext(ctx)
}
