package csv

import (
	"csv-query/db"
	"database/sql"
	"encoding/csv"
	"github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"time"
)

func createLogList(data [][]string, database *sql.DB) []db.LogRecord {
	var logList []db.LogRecord
	txn, err := database.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := txn.Prepare(pq.CopyIn(os.Getenv("DB_TABLE"), "t_stamp", "username", "upload", "file_size"))
	if err != nil {
		log.Fatal(err)
	}
	for i, line := range data {
		if i < 1 {
			continue
		}
		var rec db.LogRecord
		for j, field := range line {
			if j == 0 {
				timestamp, err := time.Parse(time.UnixDate, field)
				if err != nil {
					panic(err)
				}
				rec.Timestamp = timestamp
			} else if j == 1 {
				rec.Username = field
			} else if j == 2 {
				if field == "upload" {
					rec.Upload = true
				} else {
					rec.Upload = false
				}
			} else if j == 3 {
				size, err := strconv.Atoi(field)
				if err != nil {
					panic(err)
				}
				rec.Size = size
			}
		}
		_, err = stmt.Exec(rec.Timestamp, rec.Username, rec.Upload, rec.Size)
		if err != nil {
			log.Fatal(err)
		}

		logList = append(logList, rec)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Println(err)
	}
	err = stmt.Close()
	if err != nil {
		log.Println(err)
	}
	err = txn.Commit()
	if err != nil {
		log.Println(err)
	}
	return logList
}

func File(file string, database *sql.DB) []db.LogRecord {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	logList := createLogList(data, database)

	return logList
}
