// localizable project main.go
package main

import (
	"flag" //命令行选项解析器
	"fmt"
	"localizable/readFile"
	"strings"
)

var inFile = flag.String("i", " ", "input file path")
var outFile = flag.String("o", " ", "output file path")
var str string

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Error input,please write i2a or a2i")
		return
	}

	str = flag.Arg(0)
	if strings.EqualFold(str, "a2i") {
		readFile.ParseXML(*inFile, *outFile)
	} else if strings.EqualFold(str, "i2a") {
		readFile.ReadIOS(*inFile, *outFile)
	} else {
		fmt.Println("error, invalid para!")
	}
}
