package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

const (
	JPG = ".jpg"
	PNG = ".png"
)

func main() {
	flag.Parse()
	args := flag.Args()
	// 引数チェック
	if len(args) < 1 {
		fmt.Println("ファイル名を指定してください")
		os.Exit(1)
	}
	fileName := args[0]
	TARGET_EXTENSIONS := []string{JPG, PNG}
	ext := filepath.Ext(fileName)
	// 拡張子有無チェック
	if ext == "" || ext == "." {
		fmt.Println("ファイル名に拡張子が存在しません")
		os.Exit(1)
	}
	// ファイル名拡張子チェック
	isTargetExt := slices.ContainsFunc(TARGET_EXTENSIONS, func(s string) bool {
		return strings.EqualFold(ext, s)
	})
	if !isTargetExt {
		fmt.Printf("拡張子 %s は対応しておりません\n", ext)
		os.Exit(1)
	}
	fmt.Println(fileName)
}
