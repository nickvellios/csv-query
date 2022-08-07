package db

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"os"
	"time"
)

type UrlDB struct {
	Db *sql.DB
}

func (udb *UrlDB) Open() error {
	dbDriver := os.Getenv("DB_DRIVER")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, pass, dbName)
	var err error
	udb.Db, err = sql.Open(dbDriver, dbInfo)
	if err != nil {
		return err
	}
	err = udb.CreateDefaultTable()
	if err != nil {
		return err
	}
	udb.setAllowedFields()
	return err
}

type LogRecord struct {
	Timestamp time.Time
	Username  string
	Upload    bool
	Size      int
}

var AllowedFields []string
var FieldTypeMap map[string]string

func (udb *UrlDB) setAllowedFields() {
	rows, err := udb.Db.Query(fmt.Sprintf("SELECT * FROM %s LIMIT 1", pq.QuoteIdentifier(os.Getenv("DB_TABLE"))))
	if err != nil {
		panic(err)
	}
	FieldTypeMap = make(map[string]string)
	colTypes, _ := rows.ColumnTypes()
	for _, c := range colTypes {
		AllowedFields = append(AllowedFields, c.Name())
		FieldTypeMap[c.Name()] = c.DatabaseTypeName()
	}
}

func (udb *UrlDB) CreateDefaultTable() error {
	_, err := udb.Db.Query(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
								id SERIAL PRIMARY KEY,
								username VARCHAR(256),
								upload BOOLEAN,
								file_size INTEGER,
								t_stamp TIMESTAMP,
								UNIQUE(username, upload, file_size, t_stamp)
							);`, os.Getenv("DB_TABLE")))
	return err
}
