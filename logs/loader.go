package logs

import "csv-query/db"

type Loader interface {
	File(string, *db.UrlDB) []db.LogRecord
}
