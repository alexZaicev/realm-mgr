package postgres

import (
	"github.com/alexZaicev/realm-mgr/internal/drivers/pgdb"
	"github.com/alexZaicev/realm-mgr/internal/drivers/uuidgenerator"
)

type DataStore struct {
	uuidgen uuidgenerator.Generator
	db      pgdb.Conn
}
