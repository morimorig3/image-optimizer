package main

import (
	"flag"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

const (
	JPG  = "jpg"
	JPEG = "jpeg"
	PNG  = "png"
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
		withOutDot := strings.Replace(ext, ".", "", -1)
		return strings.EqualFold(withOutDot, s)
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

	// mimeType取得
	mt, err := getMimeType(f)
	if err != nil {
		log.Println("ファイルを開くことができませんでした")
		os.Exit(1)
	}
	if mt != "image/jpeg" && mt != "image/png" {
		log.Println("jpgまたはpngファイルを指定してください")
		os.Exit(1)
	}

	f, _ = os.Open(inputFileName)
	image, format, err := image.Decode(f)
	if err != nil {
		log.Println("ファイルを開くことができませんでした")
		os.Exit(1)
	}

	outputFileName := createOutputFileName(inputFileName, "_optimized")
	outputWriter, err := os.Create(outputFileName)
	if err != nil {
		log.Println("予期せぬエラーが発生しました")
		os.Exit(1)
	}
	defer outputWriter.Close()

	switch format {
	case JPG:
		err := optimizeJpg(outputWriter, image, 85)
		if err != nil {
			log.Println("予期せぬエラーが発生しました")
			os.Exit(1)
		}
	case JPEG:
		err := optimizeJpg(outputWriter, image, 85)
		if err != nil {
			log.Println("予期せぬエラーが発生しました")
			os.Exit(1)
		}
	case PNG:
		err := optimizePng(outputWriter, image, png.BestCompression)
		if err != nil {
			log.Println("予期せぬエラーが発生しました")
			os.Exit(1)
		}
	default:
		log.Println("予期せぬエラーが発生しました")
		os.Exit(1)
	}
	log.Println(format, "画像の圧縮が完了しました")
}

func getMimeType(f *os.File) (string, error) {
	buf := make([]byte, 512)
	// mimeType取得には最初の512byteあれば良い
	if _, err := io.ReadAtLeast(f, buf, 512); err != nil {
		return "", err
	}
	mimeType := http.DetectContentType(buf)
	return mimeType, nil
}

func optimizeJpg(w io.Writer, m image.Image, q int) error {
	err := jpeg.Encode(w, m, &jpeg.Options{
		Quality: q,
	})
	if err != nil {
		return err
	}
	return nil
}

func optimizePng(w io.Writer, m image.Image, q png.CompressionLevel) error {
	encoder := png.Encoder{
		CompressionLevel: q,
	}
	err := encoder.Encode(w, m)
	if err != nil {
		return err
	}
	return nil
}

func createOutputFileName(s string, suffix string) string {
	fileName := filepath.Base(s[:len(s)-len(filepath.Ext(s))])
	return fileName + suffix + filepath.Ext(s)
}
