package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()
	args := flag.Args()
	// 引数チェック	
	if len(args) < 1 {
		fmt.Println("ファイル名を指定してください。")
		os.Exit(1)
	}
	fileName := args[0]
	ext := filepath.Ext(fileName)
	// ファイル名拡張子チェック
	if ext != ".jpg" {
		fmt.Printf("拡張子%s は対応しておりません。\n", ext)
		os.Exit(1)
	}
	fmt.Println(fileName)
}
