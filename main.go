package main

import (
	"flag"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

const (
	JPG  = ".jpg"
	JPEG = ".jpeg"
	PNG  = ".png"
)

func main() {
	flag.Parse()
	args := flag.Args()
	// 引数チェック
	if len(args) < 1 {
		log.Println("ファイル名を指定してください")
		os.Exit(1)
	}
	inputFileName := args[0]
	TARGET_EXTENSIONS := []string{JPG, JPEG, PNG}
	ext := filepath.Ext(inputFileName)
	// 拡張子有無チェック
	if ext == "" || ext == "." {
		log.Println("ファイル名に拡張子が存在しません")
		os.Exit(1)
	}
	// ファイル名拡張子チェック
	isTargetExt := slices.ContainsFunc(TARGET_EXTENSIONS, func(s string) bool {
		return strings.EqualFold(ext, s)
	})
	if !isTargetExt {
		log.Printf("拡張子 %s は対応しておりません\n", ext)
		os.Exit(1)
	}

	// ファイルを開く
	f, err := os.Open(inputFileName)
	if err != nil {
		log.Println("ファイルを開くことができませんでした")
		os.Exit(1)
	}
	defer f.Close()

	// ファイルタイプチェック
	buf := make([]byte, 512)
	if _, err := io.ReadAtLeast(f, buf, 512); err != nil {
		log.Println("ファイルを開くことができませんでした")
		os.Exit(1)
	}
	mimeType := http.DetectContentType(buf)
	fmt.Println(mimeType)
	if mimeType != "image/jpeg" && mimeType != "image/png" {
		log.Println("jpgまたはpngファイルを指定してください")
		os.Exit(1)
	}
}
