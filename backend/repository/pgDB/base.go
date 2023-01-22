package pgDB

import (
	"database/sql"
	"fmt"

	"github.com/jerry0420/queue-system/backend/logging"
	_ "github.com/lib/pq"
)

func GetDb(username string, password string, dbLocation string, logger logging.LoggerTool) *sql.DB {
	dbConnectionString := fmt.Sprintf("postgres://%s:%s@%s",
		username,
		password,
		dbLocation,
	)

	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		logger.FATALf("db connection fail %v", err)
	}

	err = db.Ping()
	if err != nil {
		logger.FATALf("db ping fail %v", err)
	}
	return db
}
