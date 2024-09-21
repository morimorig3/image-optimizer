package main

import (
	"errors"
	"flag"
	"fmt"
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
		exitWithError(errors.New("ファイル名を指定してください"))
	}
	inputFileName := args[0]
	TARGET_EXTENSIONS := []string{JPG, JPEG, PNG}
	ext := filepath.Ext(inputFileName)
	// 拡張子有無チェック
	if ext == "" || ext == "." {
		exitWithError(errors.New("ファイル名に拡張子が存在しません"))
	}
	// ファイル名拡張子チェック
	isTargetExt := slices.ContainsFunc(TARGET_EXTENSIONS, func(s string) bool {
		withOutDot := strings.Replace(ext, ".", "", -1)
		return strings.EqualFold(withOutDot, s)
	})
	if !isTargetExt {
		exitWithError(errors.New(fmt.Sprintf("拡張子 %s は対応しておりません\n", ext)))
	}

	// ファイルを開く
	f, err := os.Open(inputFileName)
	if err != nil {
		exitWithError(err)
	}
	defer f.Close()

	// mimeType取得
	mt, err := getMimeType(f)
	if err != nil {
		exitWithError(err)
	}
	if mt != "image/jpeg" && mt != "image/png" {
		exitWithError(errors.New("jpgまたはpngファイルを指定してください"))
	}

	f, _ = os.Open(inputFileName)
	image, format, err := image.Decode(f)
	if err != nil {
		exitWithError(err)
	}

	outputFileName := createOutputFileName(inputFileName, "_optimized")
	outputWriter, err := os.Create(outputFileName)
	if err != nil {
		exitWithError(err)
	}
	defer outputWriter.Close()

	switch format {
	case JPG:
		err := optimizeJpg(outputWriter, image, 85)
		if err != nil {
			exitWithError(err)
		}
	case JPEG:
		err := optimizeJpg(outputWriter, image, 85)
		if err != nil {

			exitWithError(err)
		}
	case PNG:
		err := optimizePng(outputWriter, image, png.BestCompression)
		if err != nil {
			exitWithError(err)
		}
	default:
		exitWithError(errors.New("予期せぬエラーが発生しました"))
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

func exitWithError(err error) {
	log.Println(err)
	os.Exit(1)
}
