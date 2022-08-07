package examples

import (
	"csv-query/db"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

func ExampleAllJeffUploads(dbConn *sql.DB) {
	q := db.NewQuery(dbConn)
	q = q.AndEQ("username", "jeff22")
	q = q.AndEQ("upload", true)
	rows := q.Rows()
	for _, row := range rows {
		fmt.Printf("ExampleAllJeffUploads: Timestamp: %s | Username: %s | Upload: %s | File Size: %dkB\n", row.Timestamp.Format(time.UnixDate), row.Username, strconv.FormatBool(row.Upload), row.Size)
	}
}

func ExampleUploadsLargerThanFiftyKB(dbConn *sql.DB) {
	q := db.NewQuery(dbConn)
	q = q.AndGT("file_size", 50)
	q = q.AndEQ("upload", true)
	rows := q.Rows()
	for _, row := range rows {
		fmt.Printf("ExampleUploadsLargerThanFiftyKB: Timestamp: %s | Username: %s | Upload: %s | File Size: %dkB\n", row.Timestamp.Format(time.UnixDate), row.Username, strconv.FormatBool(row.Upload), row.Size)
	}
}

func ExampleDateObjectFileSizeSum(dbConn *sql.DB) {
	q := db.NewQuery(dbConn)
	q = q.AndGTE("file_size", 50)
	q = q.AndLTE("file_size", 55)
	q = q.AndEQ("t_stamp", time.Date(2020, 04, 14, 0, 0, 0, 0, time.UTC), true)
	summedFilesize := q.Sum("file_size")
	fmt.Printf("ExampleDateObjectFileSizeSum: %d\n", summedFilesize)
}

func ExampleDateStringFileSizeSum(dbConn *sql.DB) {
	q := db.NewQuery(dbConn)
	q = q.AndGTE("file_size", 50)
	q = q.AndLTE("file_size", 55)
	q = q.AndEQ("t_stamp", "2020-04-14", true)
	summedFilesize := q.Sum("file_size")
	fmt.Printf("ExampleDateStringFileSizeSum: %d\n", summedFilesize)
}

func ExampleJeffUploadCount(dbConn *sql.DB) {
	q := db.NewQuery(dbConn)
	q = q.AndEQ("username", "jeff22")
	q = q.AndEQ("upload", true)
	q = q.AndEQ("t_stamp", "2020-04-15", true)
	uploadCount := q.Count("id")
	fmt.Printf("ExampleJeffUploadCount: %d\n", uploadCount)
}

func ExampleJeffOrRosannaUploadCount(dbConn *sql.DB) {
	q := db.NewQuery(dbConn)

	q = q.AndEQ("upload", true)
	q = q.AndEQ("t_stamp", "2020-04-15", true)
	names := []string{"jeff22", "rosannaM"}
	q = q.AndIN("username", names)
	uploadCount := q.Count("id")
	fmt.Printf("ExampleJeffOrRosannaUploadCount: %d\n", uploadCount)
}

func ExampleCountDistinctUsers(dbConn *sql.DB) {
	q := db.NewQuery(dbConn)
	distinctUsers := q.CountDistinct("username")
	fmt.Printf("ExampleCountDistinctUsers: %d\n", distinctUsers)
}

func ExampleAverageFileSize(dbConn *sql.DB) {
	q := db.NewQuery(dbConn)
	avgSize := q.Avg("file_size")
	fmt.Printf("ExampleAverageFileSize: %d\n", avgSize)
}
