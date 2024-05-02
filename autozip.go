package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// 压缩文件夹函数
func compressFolder(folderName string) (string, error) {
	// 创建压缩包文件
	zipFile, err := os.Create(fmt.Sprintf("%s.zip", folderName))
	if err != nil {
		return "", fmt.Errorf("创建压缩包文件失败：%w", err)
	}
	defer zipFile.Close()

	// 创建压缩包写入器
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 遍历文件夹中的所有文件
	err = filepath.Walk(folderName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过文件夹
		if info.IsDir() {
			return nil
		}

		// 创建压缩包文件头
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("创建压缩包文件头失败：%w", err)
		}

		// 设置压缩包文件名称
		header.Name = filepath.Join(folderName, path[len(folderName)+1:])

		// 创建压缩包文件写入器
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("创建压缩包文件写入器失败：%w", err)
		}

		// 打开原始文件
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("打开原始文件失败：%w", err)
		}
		defer file.Close()

		// 复制文件内容到压缩包
		_, err = io.Copy(writer, file)
		if err != nil {
			return fmt.Errorf("复制文件内容到压缩包失败：%w", err)
		}

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("遍历文件夹失败：%w", err)
	}

	return fmt.Sprintf("压缩完成：%s.zip", folderName), nil
}

func main() {
	// 获取当前目录下的所有文件夹
	folders, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}

	// 创建一个 WaitGroup
	var wg sync.WaitGroup

	// 遍历每个文件夹
	for _, folder := range folders {
		if folder.IsDir() {
			// 增加 WaitGroup 计数
			wg.Add(1)

			// 启动一个 goroutine 压缩文件夹
			go func(folderName string) {
				defer wg.Done() // 完成压缩后减少 WaitGroup 计数

				// 压缩文件夹
				result, err := compressFolder(folderName)
				if err != nil {
					// 如果压缩失败，打印错误信息
					fmt.Printf("压缩文件夹 %s 失败：%s\n", folderName, err.Error())
				} else {
					// 如果压缩成功，打印压缩结果
					fmt.Println(result)
				}
			}(folder.Name())
		}
	}

	// 等待所有压缩任务完成
	wg.Wait()

	fmt.Println("所有文件夹压缩完成")
}
