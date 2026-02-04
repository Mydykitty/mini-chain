package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMergeFile(t *testing.T) {
	// 输出文件
	outFile := "merged.txt"
	out, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// 遍历当前目录及子目录
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 忽略目录
		if info.IsDir() {
			return nil
		}

		// 只处理 .go 文件，并排除 _test.go
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// 打开文件
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		// 写入文件名分隔
		out.WriteString(fmt.Sprintf("===== %s =====\n", path))

		// 拷贝文件内容
		_, err = io.Copy(out, f)
		if err != nil {
			return err
		}

		// 换行
		out.WriteString("\n\n")
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("已汇总到", outFile)
}
