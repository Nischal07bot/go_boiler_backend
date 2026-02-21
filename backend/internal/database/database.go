package database

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Database struct {
	Pool *pgxpool.Pool//this will be exported to be used in other packages do convention capital letters
	log  *zerolog.Logger
}

const DatabasePingTimeout=10
