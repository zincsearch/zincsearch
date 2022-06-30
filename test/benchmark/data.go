package benchmark

import (
	"fmt"
	"os"
)

func getOlympicsFile() *os.File {
	f, err := os.Open("../../tmp/olympics.ndjson")
	if err != nil {
		fmt.Println("# !!! Please download olympics.ndjosn first")
		fmt.Println("# mkdir -p tmp")
		fmt.Println("# wget https://github.com/zinclabs/zinc/releases/download/v0.2.4/olympics.ndjson.tar.gz -O tmp/olympics.ndjson.tar.gz")
		fmt.Println("# tar zxf tmp/olympics.ndjson.tar.gz -C tmp/")
		panic(err)
	}
	return f
}
