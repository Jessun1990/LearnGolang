package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func main() {
	loc, _ := time.LoadLocation("Local")
	StartAt, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-05-23 20:30:00", loc)

	jsonStu, _ := json.Marshal([]string{StartAt.Format(time.RFC3339)})

	fmt.Println(string(jsonStu), StartAt.Format(time.RFC3339), time.Now().Format(time.RFC3339))
}
