package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("ファイル名を指定してください。")
		os.Exit(1)
	}
	fileName := args[0]
	fmt.Println(fileName)
}
