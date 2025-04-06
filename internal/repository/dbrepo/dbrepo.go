package dbrepo

import (
	"database/sql"
	"github.com/sangketkit01/bookings/internal/config"
	"github.com/sangketkit01/bookings/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}


// NewPostgresRepo creates a new postgres repository
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DataBaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}


// NewTestingRepo creates a new testing postgres repository
func NewTestingRepo(a *config.AppConfig) repository.DataBaseRepo {
	return &testDBRepo{
		App: a,
	}
}
