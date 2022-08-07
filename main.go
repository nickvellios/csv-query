package main

import (
	"csv-query/db"
	"csv-query/examples"
	"csv-query/logs/csv"
	"fmt"
	"github.com/c-bata/go-prompt"
)

func exampleChoices(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "EXIT", Description: "Exit program"},
		{Text: "jeff22 Upload Count On 4/15/20", Description: "Return number of files uploaded by jeff22 on 4/15/2020."},
		{Text: "jeff22 Show Uploads", Description: "List all uploads by user jeff22."},
		{Text: "jeff22/rosannaM Upload Count On 4/15/20", Description: "Return number of files uploaded by either jeff22 or rosannaM on 4/15/2020."},
		{Text: "file_size Sum By Date Object & Size", Description: "Return sum of all files < 55kB and > 50kB uploaded/downloaded on 4/14/2020 using a time object."},
		{Text: "file_size Sum By Date String & Size", Description: "Return sum of all files < 55kB and > 50kB uploaded/downloaded on 4/14/2020 using a string parameter."},
		{Text: "Uploads > 50kB", Description: "List all uploads greater than 50kB in size."},
		{Text: "Average file_size", Description: "Return the average file size uploaded or downlodaed"},
		{Text: "username Count", Description: "Return count of distinct usernames."},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func filein(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "./test_data/server_log.csv", Description: "Probably the file you're looking for"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func main() {
	dbConn := db.SetupEnvironment().Db
	defer func() {
		err := dbConn.Close()
		if err != nil {
			panic(err)
		}
	}()

	fmt.Println("Please select CSV file:")
	in := prompt.Input("➜ ", filein)
	csv.File(in, dbConn)
	fmt.Println("Successfully loaded file: " + in)

	for true {
		switch opt := prompt.Input("(tab through options) ➜ ", exampleChoices); opt {
		case "jeff22 Show Uploads":
			examples.ExampleAllJeffUploads(dbConn)
		case "Uploads > 50kB":
			examples.ExampleUploadsLargerThanFiftyKB(dbConn)
		case "file_size Sum By Date Object & Size":
			examples.ExampleDateObjectFileSizeSum(dbConn)
		case "file_size Sum By Date String & Size":
			examples.ExampleDateStringFileSizeSum(dbConn)
		case "jeff22 Upload Count On 4/15/20":
			examples.ExampleJeffUploadCount(dbConn)
		case "jeff22/rosannaM Upload Count On 4/15/20":
			examples.ExampleJeffOrRosannaUploadCount(dbConn)
		case "username Count":
			examples.ExampleCountDistinctUsers(dbConn)
		case "Average file_size":
			examples.ExampleAverageFileSize(dbConn)
		case "EXIT":
			fmt.Println("Goodbye")
			return
		}
	}
}
