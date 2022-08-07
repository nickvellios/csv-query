# csv-query

Take-home assignment for `[REDACTED]`, Inc.

### Given a log file of events, csv-query provides the ability to gather metrics on user activity.

Log file is in CSV format with the following structure:


| timestamp                    | username | operation | size (kB) |
|------------------------------|----------|-----------|-----------|
| Sun Apr 12 22:10:38 UTC 2020 | erlichB  | upload    | 45        |
| Sun Apr 12 22:12:30 UTC 2020 | jianY    | download  | 76        |
| Sun Apr 12 22:18:01 UTC 2020 | gilfoyle | upload    | 56812     |
| Sun Apr 12 22:43:54 UTC 2020 | richardH | upload    | 11        |
| Sun Apr 12 22:43:59 UTC 2020 | richardH | upload    | 1         |
| Sun Apr 12 22:44:06 UTC 2020 | richardH | upload    | 9         |
| Sun Apr 12 22:48:37 UTC 2020 | dinesh   | download  | 127       |
| Sun Apr 12 23:02:02 UTC 2020 | jianY    | download  | 547       |

---
## Overview

Given the requirements, the most intuitive solution was to pre-process the CSV file into a data structure allowing for flexible aggregation, summation, and the ability to filter result set by any row or column combination.

Fortunately, there exists multiple open-source solutions to at least part of the problem.  In this case, I selected Postgres, but most RBDMS should plug and play with only minor changes to the query compiler.

Query syntax resembles that of an ORM, although it isn't quite fleshed out enough to be considered an ORM.  If I had an additional two days though...

This could be considered more of a proof of concept than a solution.  I'll provide more context which should shed light on that below.

---
## How to use demo

1. Configure environment variables in: `csv-query/db/environment.go`
2. Build and run:
   ```
   ➜ go build
   ➜ go install
   ➜ csv-query
   ```
3. Follow the prompts, entering a path to the log file:
   ```
   Please select CSV file:
   ➜ ./test_data/server_log.csv

   ```
   If your database settings are configured properly, the logs will be imported into your database under a table name specified in the environment variable `DB_TABLE`
4. `TAB` through the various examples to query the dataset.
---
## How to query data manually

There is currently not a way to use the demo console to perform custom queries.  However, you can add some additional test functions to `csv-query/examples/examples.go` and link them to the prompt in `csv-query/main.go`'s `exampleChoices` and `main` functions.

### Usage basics

Build a new `Combined Query` object by passing in the current database connection:
```
query := db.NewQuery(dbConn)
```
Filter result set using query functions.
```
query = query.AndEQ("username", "nick")
```
Each operation returns an updated combined query object with all prior filters included.
```
query = query.AndEQ("username", "nick")
query = query.AndGTE("file_size", 50)
```
These can also be chained similar to Django querysets.

```
query = query.AndEQ("username", "nick").AndGTE("file_size", 50)
```
To return a list of applicable rows, call the `Rows` function on the combined query to generate a `[]LogRecord` object array that can be iterated over.
```
rows := query.Rows()
for _, row := range rows {
    fmt.Printf("Timestamp: %s | Username: %s | Upload: %s | File Size: %dkB\n", row.Timestamp.Format(time.UnixDate), row.Username, strconv.FormatBool(row.Upload), row.Size)
}

Output:
> Timestamp: Sun Apr 26 02:42:52 +0000 2020 | Username: nick | Upload: true | File Size: 68kB
> Timestamp: Sun Apr 26 12:28:01 +0000 2020 | Username: nick | Upload: false | File Size: 66kB
> Timestamp: Sun Apr 26 15:12:00 +0000 2020 | Username: nick | Upload: true | File Size: 73kB
```
### There are also aggregation functions to allow pulling summarization data.

Count distinct field values:
```
distinctUsers := query.CountDistinct("username")
```
View average of a field in a combined query.  For example, average file size:
```
avgSize := query.Avg("file_size")
```
See how many times a user has uploaded a file:
```
query = query.AndEQ("upload", true)
query = query.AndEQ("username", "nick")
uploadCount := query.Count("id")
```
See how many times a list of users has uploaded a file:
```
query = query.AndEQ("upload", true)
names := []string{"nick", "gumboTheWonderPuppy"}
query = query.AndIN("username", names)
uploadCount := query.Count("id")
```
Find the sum of all file sizes of files uploaded with size greater than or equal to 50kB and less than or equal to 55kB:
```
query = query.AndGTE("file_size", 50)
query = query.AndLTE("file_size", 55)
summedFilesize := query.Sum("file_size")
```
Filter results by date using a date string, ignoring timestamp:
```
query = query.AndEQ("t_stamp", "2020-04-14", true)  // All filter functions allow an optional trailing parameter to Cast timestamps to dates.
```
Filter results by date using a date/time object, ignoring timestamp:
```
query = query.AndEQ("t_stamp", time.Date(2020, 04, 14, 0, 0, 0, 0, time.UTC), true)  // All filter functions allow an optional trailing parameter to Cast timestamps to dates.
```
Aside from the `AndIN` functions, you also have access to `OR` operations, although currently there is no way to parenthesize/contain their logic, so only useful for simple queries.
The Following will filter results where `file_size > 50kB OR username = gumboTheWonderPuppy`
```
query = query.AndGTE("file_size", 50)
query = query.OrEQ("username", "gumboTheWonderPuppy")
```

## *** IMPORTANT NOTE ***
The first filter whether it's preficed by `And` or `Or` will compile into a `WHERE` clause.  For example:
```
query = query.AndGTE("file_size", 50)
query = query.OrEQ("username", "gumboTheWonderPuppy")
```
Compiles into:
```
SELECT * FROM `my_table` WHERE `file_size` >= 50 OR `username` = "gumboTheWonderPuppy";
```
---
## Filter Functions

```
func (cq *combinedQuery) AndIN(field string, value any, castDate ...bool) *combinedQuery { ... }

func (cq *combinedQuery) AndNotIN(field string, value any, castDate ...bool) *combinedQuery { ... }

func (cq *combinedQuery) AndEQ(field string, value interface{}, castDate ...bool) *combinedQuery { ... }

func (cq *combinedQuery) AndGT(field string, value interface{}, castDate ...bool) *combinedQuery { ... }

func (cq *combinedQuery) AndGTE(field string, value interface{}, castDate ...bool) *combinedQuery { ... }

func (cq *combinedQuery) AndLT(field string, value interface{}, castDate ...bool) *combinedQuery { ... }

func (cq *combinedQuery) AndLTE(field string, value interface{}, castDate ...bool) *combinedQuery { ... }

func (cq *combinedQuery) AndIsNot(field string, value interface{}, castDate ...bool) *combinedQuery { ... }

func (cq *combinedQuery) AndNE(field string, value interface{}, castDate ...bool) *combinedQuery { ... }

func (cq *combinedQuery) OrIN(field string, value any, castDate ...bool) *combinedQuery { ... }

func (cq *combinedQuery) OrNotIn(field string, value any, castDate ...bool) *combinedQuery { ... }

func (cq *combinedQuery) OrEQ(field string, value interface{}, castDate ...bool) *combinedQuery { ... }

func (cq *combinedQuery) OrGT(field string, value interface{}, castDate ...bool) *combinedQuery { ... }

func (cq *combinedQuery) OrGTE(field string, value interface{}, castDate ...bool) *combinedQuery { ... }

func (cq *combinedQuery) OrLT(field string, value interface{}, castDate ...bool) *combinedQuery { ... }

func (cq *combinedQuery) OrLTE(field string, value interface{}, castDate ...bool) *combinedQuery { ... }

func (cq *combinedQuery) OrIsNot(field string, value interface{}, castDate ...bool) *combinedQuery { ... }

func (cq *combinedQuery) OrNE(field string, value interface{}, castDate ...bool) *combinedQuery { ... }
```
---
## Aggregate Functions

```
func (cq *combinedQuery) Sum(field string) int { ... }

func (cq *combinedQuery) Count(field string) int { ... }

func (cq *combinedQuery) CountDistinct(field string) int { ... }

func (cq *combinedQuery) Avg(field string) uint8 { ... }
```
---
## List Results Functions

```
func (cq *combinedQuery) Rows() []LogRecord { ... }
```
---
## The Good

For simpler data queries over fixed log structures, this offers an abstraction away from using SQL directly without a heavy ORM.

---
## The Bad ...err not so good

### Complex queries are not supported.  That means no:

- `JOIN`
- `SUBQUERY`
- `GROUP`
- `HAVING`
- Parenthesized conditions e.g. `WHERE x = 1 OR (y = 'p' AND x = 2)`
- Aliasing, casting (other than datetime to dates), `CASE`/`WHEN`/`THEN` or pretty much _anything_ fun.

### Other limitations

- Input data currently bound to a fixed file structure and type.  Some work was started to support flexibility in that regard.  See time permitting section below.
- No way to add a new query without recompiling.
- No API.
---
## Time permitting...

My entire town conveniently lost power for the first half of the first of the two allotted days to complete the assignment and cell data doesn't work here.  That would be fine, except I hadn't used Go in almost 5 years which made this an especially fun challenge.

- Cleanup file structure.  Went in a couple directions here and it shows.
- A way to use the console to use predictive text to build a query.
- Dynamic data structure/schema.
- Support for multiple data storage backends (MySQL, NoSQL, File)
- Expanded support for complex queries.
- REST API.
- RPC API.
- Support for multiple file formats including piping tailed logs directly in.
- Support for accessing only certain fields when listing results.
- Data transformation.
- Cleanup of types.  Explore the new Go generics and see how they could be applied to aleviate type coersion.
- Unit & Integration tests.
- More intuitive naming of query functions so the first applied filter's condition isn't replaced with `WHERE` clause.