package influxdb

import (
	"context"

	"github.com/Ramazon1227/BeatSync/config"
	"github.com/Ramazon1227/BeatSync/storage"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type Store struct {
	db   *influxdb2.Client
	user storage.UserRepoI
}

func (s *Store) CloseDB() {
	(*s.db).Close()
}

func NewInfluxDB(ctx context.Context, cfg config.Config) (storage.StorageI,error) {
	client := influxdb2.NewClient(cfg.InfluxURL, cfg.InfluxToken)

	return &Store{
		db:    &client,
		user : NewUserRepoI(&client),
	}, nil
}


func (s *Store) User() storage.UserRepoI{
	if s.user == nil {
		s.user = NewUserRepoI(s.db)
	}

	return s.user
}
