package influxdb

import (
	"context"
	"log"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/Ramazon1227/BeatSync/config"
	// "github.com/Ramazon1227/BeatSync/pkg/logger"
	"github.com/Ramazon1227/BeatSync/storage"
	// influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type Store struct {
	db   *influxdb3.Client
	user storage.UserRepoI
	analyze storage.AnalyzeRepoI
}

func (s *Store) CloseDB() {
	s.db.Close()
}

func NewInfluxDB(ctx context.Context, cfg *config.Config) (storage.StorageI,error) {
	client,err := influxdb3.New(influxdb3.ClientConfig{
		Host:  cfg.InfluxURL,
		Token: cfg.InfluxToken,
	  })
	if err != nil {
		log.Fatal("Error creating InfluxDB client: ", err)
		return nil, err
	}

	return &Store{
		db:    client,
		user : NewUserRepoI(client),
	}, nil
}


func (s *Store) User() storage.UserRepoI{
	if s.user == nil {
		s.user = NewUserRepoI(s.db)
	}

	return s.user
}
