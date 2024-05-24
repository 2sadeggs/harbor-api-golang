package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// 用于下载并保存所有制品的函数
func downloadAndSaveArtifacts(baseURL, auth string) error {
	startTime := time.Now()
	fmt.Printf("Start time: %s\n", startTime.Format("2006-01-02 15:04:05.000000000"))

	uris, err := fetchNonUnknownArchURIs(baseURL, auth)
	if err != nil {
		return err
	}

	// 创建一个以时间戳命名的保存目录
	timestamp := time.Now().Format("2006-01-02_15-04-05.000000000")
	savePath := filepath.Join(".", "artifacts", timestamp)
	err = os.MkdirAll(savePath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create save directory: %v", err)
	}

	// 使用带缓冲的 channel 来限制并发 goroutine 数量
	concurrencyLimit := 5 // 并发数量限制
	semaphore := make(chan struct{}, concurrencyLimit)
	var wg sync.WaitGroup

	for _, uri := range uris {
		wg.Add(1)
		go func(uri string) {
			defer wg.Done()

			semaphore <- struct{}{}        // 获取信号量
			defer func() { <-semaphore }() // 释放信号量

			fmt.Printf("Downloading artifact: %s\n", uri)
			cmd := exec.Command("docker", "pull", uri)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("failed to download artifact %s: %v\nOutput: %s", uri, err, string(output))
				return
			}
			fmt.Printf("Successfully downloaded artifact: %s\n", uri)

			// Save the artifact to a file
			fileName := fmt.Sprintf("%s.tar", uriToFileName(uri))
			filePath := filepath.Join(savePath, fileName)

			fmt.Printf("Saving artifact to: %s\n", filePath)
			saveCmd := exec.Command("docker", "save", "-o", filePath, uri)
			saveOutput, saveErr := saveCmd.CombinedOutput()
			if saveErr != nil {
				fmt.Printf("failed to save artifact %s: %v\nOutput: %s", uri, saveErr, string(saveOutput))
				return
			}
			fmt.Printf("Successfully saved artifact: %s\n", filePath)
		}(uri)
	}

	wg.Wait()

	endTime := time.Now()
	fmt.Printf("End time: %s\n", endTime.Format("2006-01-02 15:04:05.000000000"))
	fmt.Printf("Duration: %s\n", endTime.Sub(startTime))

	return nil
}

// 用于下载所有制品
func downloadArtifacts(baseURL, auth string) error {
	//_, _, _, uris, _, _, err := fetchAllURIs(baseURL, auth)
	uris, err := fetchNonUnknownArchURIs(baseURL, auth)
	if err != nil {
		return err
	}

	for _, uri := range uris {
		fmt.Printf("Downloading artifact: %s\n", uri)
		cmd := exec.Command("docker", "pull", uri)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to download artifact %s: %v\nOutput: %s", uri, err, string(output))
		}
		fmt.Printf("Successfully downloaded artifact: %s\n", uri)
	}

	return nil
}

/*
// 用于下载并保存所有制品的函数
func downloadAndSaveArtifacts(baseURL, auth string) error {
	startTime := time.Now()
	fmt.Printf("Start time: %s\n", startTime.Format("2006-01-02 15:04:05.000000000"))

	uris, err := fetchNonUnknownArchURIs(baseURL, auth)
	if err != nil {
		return err
	}

	// 创建一个以时间戳命名的保存目录
	timestamp := time.Now().Format("2006-01-02_15-04-05.000000000")
	savePath := filepath.Join(".", "artifacts", timestamp)
	err = os.MkdirAll(savePath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create save directory: %v", err)
	}

	for _, uri := range uris {
		fmt.Printf("Downloading artifact: %s\n", uri)
		cmd := exec.Command("docker", "pull", uri)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to download artifact %s: %v\nOutput: %s", uri, err, string(output))
		}
		fmt.Printf("Successfully downloaded artifact: %s\n", uri)

		// Save the artifact to a file
		fileName := fmt.Sprintf("%s.tar", uriToFileName(uri))
		filePath := filepath.Join(savePath, fileName)

		fmt.Printf("Saving artifact to: %s\n", filePath)
		saveCmd := exec.Command("docker", "save", "-o", filePath, uri)
		saveOutput, saveErr := saveCmd.CombinedOutput()
		if saveErr != nil {
			return fmt.Errorf("failed to save artifact %s: %v\nOutput: %s", uri, saveErr, string(saveOutput))
		}
		fmt.Printf("Successfully saved artifact: %s\n", filePath)
	}

	endTime := time.Now()
	fmt.Printf("End time: %s\n", endTime.Format("2006-01-02 15:04:05.000000000"))
	fmt.Printf("Duration: %s\n", endTime.Sub(startTime))

	return nil
}*/

// 将 URI 转换为文件名的辅助函数
func uriToFileName(uri string) string {
	// 替换不适合文件名的字符
	uri = strings.ReplaceAll(uri, "/", "_")
	uri = strings.ReplaceAll(uri, "@", "__")
	uri = strings.ReplaceAll(uri, ":", "___")
	return uri
}
