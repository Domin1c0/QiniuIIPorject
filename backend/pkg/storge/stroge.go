package storage

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

type Storage struct {
	engine *xorm.Engine
	logger *zerolog.Logger
}

func NewStorage(engineType, engineUrl string, logger *zerolog.Logger) (*Storage, error) {
	if logger == nil {
		logger = &log.Logger
	}

	eng, err := xorm.NewEngine(parseDriverName(engineType), engineUrl)
	if err != nil {
		return nil, err
	}
	eng.SetLogger(&xormLogger{logger: logger})

	if err := eng.Ping(); err != nil {
		return nil, err
	}

	// sync db
	if err := eng.Sync(&User{}, &UserSession{}, &Session{}, &Message{}); err != nil {
		return nil, err
	}

	return &Storage{logger: logger, engine: eng}, nil
}

func parseDriverName(engineType string) string {
	if engineType == "postgresql" {
		return "pgx"
	}
	// if engineType == "mariadb" {
	// 	return "mysql"
	// }
	return engineType
}

func (db *Storage) Init() error {
	if err := db.engine.Sync(&User{}, &UserSession{}, &Session{}, &Message{}); err != nil {
		return err
	}
	// TODO foreign key
	return nil
}
