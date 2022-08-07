package db

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"os"
	"strings"
	"time"
)

// clause statement mapping
const (
	WHERE string = "WHERE"
	AND          = "AND"
	OR           = "OR"
)

// conditional statement mapping
const (
	EQ    string = "="
	GT           = ">"
	GTE          = ">="
	LT           = "<"
	LTE          = "<="
	NE           = "<>"
	IN           = "IN"
	NOTIN        = "NOT IN"
	ISNT         = "IS NOT"
)

type queryFilter struct {
	clause   string
	field    string
	operator string
	castDate bool
	value    any
}

type combinedQuery struct {
	ids     []int
	filters []queryFilter
	dB      *sql.DB
}

func NewQuery(db *sql.DB) *combinedQuery {
	return &combinedQuery{
		ids:     []int{},
		filters: []queryFilter{},
		dB:      db,
	}
}

func (cq *combinedQuery) add(qf queryFilter) *combinedQuery {
	valid := false
	for _, val := range AllowedFields {
		if val == qf.field {
			valid = true
		}
	}
	if valid == false {
		panic(fmt.Sprintf("%s not in AllowedFields", qf.field))
	}
	cq.filters = append(cq.filters, qf)
	return cq
}

func (cq *combinedQuery) constructFilter(clause string, field string, value any, operator string, castDate ...bool) queryFilter {
	cd := false
	if len(castDate) > 0 {
		cd = castDate[0]
	}
	qf := queryFilter{
		clause:   clause,
		field:    field,
		operator: operator,
		castDate: cd,
		value:    value,
	}
	return qf
}

func (cq *combinedQuery) and(field string, value any, operator string, castDate ...bool) queryFilter {
	return cq.constructFilter(AND, field, value, operator, castDate...)
}

func (cq *combinedQuery) or(field string, value any, operator string, castDate ...bool) queryFilter {
	return cq.constructFilter(OR, field, value, operator, castDate...)
}

func (cq *combinedQuery) AndIN(field string, value any, castDate ...bool) *combinedQuery {
	qf := cq.and(field, value, IN, castDate...)
	return cq.add(qf)
}

func (cq *combinedQuery) AndNotIN(field string, value any, castDate ...bool) *combinedQuery {
	qf := cq.and(field, value, NOTIN, castDate...)
	return cq.add(qf)
}

func (cq *combinedQuery) AndEQ(field string, value interface{}, castDate ...bool) *combinedQuery {
	qf := cq.and(field, value, EQ, castDate...)
	return cq.add(qf)
}

func (cq *combinedQuery) AndGT(field string, value interface{}, castDate ...bool) *combinedQuery {
	qf := cq.and(field, value, GT, castDate...)
	return cq.add(qf)
}

func (cq *combinedQuery) AndGTE(field string, value interface{}, castDate ...bool) *combinedQuery {
	qf := cq.and(field, value, GTE, castDate...)
	return cq.add(qf)
}

func (cq *combinedQuery) AndLT(field string, value interface{}, castDate ...bool) *combinedQuery {
	qf := cq.and(field, value, LT, castDate...)
	return cq.add(qf)
}

func (cq *combinedQuery) AndLTE(field string, value interface{}, castDate ...bool) *combinedQuery {
	qf := cq.and(field, value, LTE, castDate...)
	return cq.add(qf)
}

func (cq *combinedQuery) AndIsNot(field string, value interface{}, castDate ...bool) *combinedQuery {
	qf := cq.and(field, value, ISNT, castDate...)
	return cq.add(qf)
}

func (cq *combinedQuery) AndNE(field string, value interface{}, castDate ...bool) *combinedQuery {
	qf := cq.and(field, value, NE, castDate...)
	return cq.add(qf)
}

func (cq *combinedQuery) OrIN(field string, value any, castDate ...bool) *combinedQuery {
	qf := cq.or(field, value, IN, castDate...)
	return cq.add(qf)
}

func (cq *combinedQuery) OrNotIn(field string, value any, castDate ...bool) *combinedQuery {
	qf := cq.or(field, value, NOTIN, castDate...)
	return cq.add(qf)
}

func (cq *combinedQuery) OrEQ(field string, value interface{}, castDate ...bool) *combinedQuery {
	qf := cq.or(field, value, EQ, castDate...)
	return cq.add(qf)
}

func (cq *combinedQuery) OrGT(field string, value interface{}, castDate ...bool) *combinedQuery {
	qf := cq.or(field, value, GT, castDate...)
	return cq.add(qf)
}

func (cq *combinedQuery) OrGTE(field string, value interface{}, castDate ...bool) *combinedQuery {
	qf := cq.or(field, value, GTE, castDate...)
	return cq.add(qf)
}

func (cq *combinedQuery) OrLT(field string, value interface{}, castDate ...bool) *combinedQuery {
	qf := cq.or(field, value, LT, castDate...)
	return cq.add(qf)
}

func (cq *combinedQuery) OrLTE(field string, value interface{}, castDate ...bool) *combinedQuery {
	qf := cq.or(field, value, LTE, castDate...)
	return cq.add(qf)
}

func (cq *combinedQuery) OrIsNot(field string, value interface{}, castDate ...bool) *combinedQuery {
	qf := cq.or(field, value, ISNT, castDate...)
	return cq.add(qf)
}

func (cq *combinedQuery) OrNE(field string, value interface{}, castDate ...bool) *combinedQuery {
	qf := cq.or(field, value, NE, castDate...)
	return cq.add(qf)
}

func (cq *combinedQuery) preCompile() ([]string, []any) {
	statements := make([]string, len(cq.filters))
	args := make([]any, 0)
	paramCounter := 0
	for i, f := range cq.filters {
		statement := ""
		clause := f.clause
		if i == 0 {
			clause = WHERE
		}
		if FieldTypeMap[f.field] == "TIMESTAMP" && f.castDate == true {
			dateStr := f.value
			if _, ok := f.value.(string); !ok {
				dateStr = fmt.Sprintf("%s", f.value.(time.Time).Format("2006-01-02"))
			}
			statement = fmt.Sprintf("%s%s %s::DATE %s '%s' AND $%d = 1", clause, statement, f.field, f.operator, dateStr, paramCounter+1)
			args = append(args, "1") // @todo Figure out later. Hack to make this query work. Adds an always true condition: 1 = 1
		} else {
			if f.operator == IN || f.operator == NOTIN {
				values := f.value.([]string)
				argLen := len(values)
				subStrArr := make([]string, 0)

				for i := 0; i < argLen; i++ {
					paramCounter++
					subStrArr = append(subStrArr, fmt.Sprintf("$%d", paramCounter))

					args = append(args, values[i])
				}
				subStr := strings.Join(subStrArr, ",")
				statement = fmt.Sprintf(" %s%s %s %s (%s) ", clause, statement, f.field, f.operator, subStr)
			} else {
				statement = fmt.Sprintf("%s%s %s %s $%d ", clause, statement, f.field, f.operator, paramCounter+1)
				args = append(args, f.value)
			}
		}
		paramCounter++
		statements = append(statements, statement)

	}
	return statements, args
}

func (cq *combinedQuery) Rows() []LogRecord {
	statements, args := cq.preCompile()
	queryStatement := fmt.Sprintf("SELECT t_stamp, username, upload, file_size FROM %s %s", pq.QuoteIdentifier(os.Getenv("DB_TABLE")), strings.Join(statements, ""))
	rows, err := cq.dB.Query(queryStatement, args...)
	if err != nil {
		fmt.Println(err.Error())
	}

	var logRecords []LogRecord
	for rows.Next() {
		logRecord := LogRecord{}
		err = rows.Scan(&logRecord.Timestamp, &logRecord.Username, &logRecord.Upload, &logRecord.Size)
		if err != nil {
			fmt.Println(err.Error())
		}
		logRecords = append(logRecords, logRecord)
	}

	return logRecords
}

func (cq *combinedQuery) Sum(field string) int {
	return cq.aggregate(fmt.Sprintf("SUM(%s)", field)).(int)
}

func (cq *combinedQuery) Count(field string) int {
	return cq.aggregate(fmt.Sprintf("COUNT(%s)", field)).(int)
}

func (cq *combinedQuery) CountDistinct(field string) int {
	return cq.aggregate(fmt.Sprintf("COUNT(DISTINCT %s)", field)).(int)
}

func (cq *combinedQuery) Avg(field string) uint8 {
	return cq.aggregate(fmt.Sprintf("AVG(%s)", field), true).(uint8)
}

func (cq *combinedQuery) aggregate(function string, decimal ...bool) any {
	statements, args := cq.preCompile()
	queryStatement := fmt.Sprintf("SELECT %s AS aggregated_value FROM %s %s", function, pq.QuoteIdentifier(os.Getenv("DB_TABLE")), strings.Join(statements, ""))
	rows, err := cq.dB.Query(queryStatement, args...)
	if err != nil {
		panic(err)
	}
	isDec := false
	if len(decimal) > 0 {
		isDec = decimal[0]
	}
	if isDec == false {
		aggregatedVal, err := getIntAggregate(rows)
		if err != nil {
			panic(err)
		}
		return aggregatedVal
	} else {
		aggregatedVal, err := getDecAggregate(rows)
		if err != nil {
			panic(err)
		}
		return aggregatedVal
	}
}

func getIntAggregate(rows *sql.Rows) (int, error) {
	for rows.Next() {
		var aggregatedVal sql.NullInt32
		err := rows.Scan(&aggregatedVal)
		if err != nil {
			return 0, err
		}
		if aggregatedVal.Valid {
			return int(aggregatedVal.Int32), nil
		}
	}

	return 0, nil
}

func getDecAggregate(rows *sql.Rows) (uint8, error) {
	for rows.Next() {
		var aggregatedVal []uint8
		err := rows.Scan(&aggregatedVal)
		if err != nil {
			return 0, err
		}
		if len(aggregatedVal) < 1 {
			return 0, nil
		}
		return aggregatedVal[0], nil
	}

	return 0, nil
}
