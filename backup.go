package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// 保存上次全量备份路径的文件
const lastBackupPathFile = "./last_full_backup_path.txt"

// downloadAndSaveAllArtifacts 全量备份
func downloadAndSaveAllArtifacts(baseURL, auth string) error {
	startTime := time.Now()
	fmt.Printf("Start time: %s\n", startTime.Format("2006-01-02 15:04:05.000000000"))

	// 调用 fetchAllArtifactsWithTypes 获取 URI 列表
	artifactURIs, err := fetchAllArtifactsWithTypes(baseURL, auth)
	if err != nil {
		fmt.Printf("Error fetching artifacts: %v\n", err)
		return err
	}

	// 获取 non_unknown_arch_uris 类型的 URI 列表
	nonUnknownArchURIs, ok := artifactURIs["non_unknown_arch_uris"]
	if !ok {
		fmt.Println("No non_unknown_arch_uris found.")
		return err
	}

	// 创建一个以时间戳命名的保存目录，包含 "full" 标识
	timestamp := time.Now().Format("2006-01-02_15-04-05.000000000")
	savePath := filepath.Join(".", "artifacts", "full_"+timestamp)
	err = os.MkdirAll(savePath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create save directory: %v", err)
	}

	// 创建一个清单文件
	listFilePath := filepath.Join(savePath, "all_uri_list.txt")
	err = saveURIsToFile(listFilePath, nonUnknownArchURIs)
	if err != nil {
		return fmt.Errorf("failed to save URI list: %v", err)
	}

	// 保存最新备份路径
	err = saveLastBackupPath(savePath)
	if err != nil {
		return fmt.Errorf("failed to save last backup path: %v", err)
	}

	// 使用带缓冲的 channel 来限制并发 goroutine 数量
	concurrencyLimit := 5 // 并发数量限制
	semaphore := make(chan struct{}, concurrencyLimit)
	var wg sync.WaitGroup

	for _, uri := range nonUnknownArchURIs {
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

// downloadAndSaveDeltaArtifactsWithDiffList 差量备份，并保存差异清单
func downloadAndSaveDeltaArtifacts(baseURL, auth string) error {
	startTime := time.Now()
	fmt.Printf("Start time: %s\n", startTime.Format("2006-01-02 15:04:05.000000000"))

	// 调用 fetchAllArtifactsWithTypes 获取 URI 列表
	artifactURIs, err := fetchAllArtifactsWithTypes(baseURL, auth)
	if err != nil {
		fmt.Printf("Error fetching artifacts: %v\n", err)
		return err
	}

	// 获取 non_unknown_arch_uris 类型的 URI 列表
	nonUnknownArchURIs, ok := artifactURIs["non_unknown_arch_uris"]
	if !ok {
		fmt.Println("No non_unknown_arch_uris found.")
		return err
	}

	// 获取上次全量备份的路径
	lastBackupPath, err := getLastBackupPath()
	if err != nil {
		return err
	}

	// 从上次全量备份的清单文件中读取上次备份的 URI 列表
	previousListFile := filepath.Join(lastBackupPath, "all_uri_list.txt")
	previousURIs, err := readURIsFromFile(previousListFile)
	if err != nil {
		return err
	}

	// 找出新的或变更的 URI
	newOrChangedURIs := findNewOrChangedURIs(nonUnknownArchURIs, previousURIs)
	if len(newOrChangedURIs) == 0 {
		fmt.Println("No new or changed artifacts to download.")
		return nil
	}

	// 创建以时间戳命名的保存目录，包含 "delta" 标识
	timestamp := time.Now().Format("2006-01-02_15-04-05.000000000")
	savePath := filepath.Join(".", "artifacts", "delta_"+timestamp)
	err = os.MkdirAll(savePath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create save directory: %v", err)
	}

	// 创建差异清单文件
	diffListFilePath := filepath.Join(savePath, "diff_list.txt")
	err = saveURIsToFile(diffListFilePath, newOrChangedURIs)
	if err != nil {
		return fmt.Errorf("failed to save diff URI list: %v", err)
	}

	// 创建一个清单文件
	listFilePath := filepath.Join(savePath, "all_uri_list.txt")
	err = saveURIsToFile(listFilePath, nonUnknownArchURIs)
	if err != nil {
		return fmt.Errorf("failed to save URI list: %v", err)
	}

	// 使用带缓冲的 channel 来限制并发 goroutine 数量
	concurrencyLimit := 5 // 并发数量限制
	semaphore := make(chan struct{}, concurrencyLimit)
	var wg sync.WaitGroup

	for _, uri := range newOrChangedURIs {
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

// findNewOrChangedURIs 查找新增或更改的制品 URI
func findNewOrChangedURIs(currentURIs, previousURIs []string) []string {
	previousURISet := make(map[string]struct{}, len(previousURIs))
	for _, uri := range previousURIs {
		previousURISet[uri] = struct{}{}
	}

	var newOrChangedURIs []string
	for _, uri := range currentURIs {
		if _, found := previousURISet[uri]; !found {
			newOrChangedURIs = append(newOrChangedURIs, uri)
		}
	}

	return newOrChangedURIs
}

// readURIsFromFile 从文件读取 URI 列表
func readURIsFromFile(filePath string) ([]string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	uris := strings.Split(string(data), "\n")
	if len(uris) > 0 && uris[len(uris)-1] == "" {
		uris = uris[:len(uris)-1]
	}
	return uris, nil
}

// saveURIsToFile 将 URI 列表保存到文件
func saveURIsToFile(filePath string, uris []string) error {
	data := strings.Join(uris, "\n")
	return ioutil.WriteFile(filePath, []byte(data), 0644)
}

// saveLastBackupPath 保存上次备份路径到固定文件
func saveLastBackupPath(path string) error {
	return ioutil.WriteFile(lastBackupPathFile, []byte(path), 0644)
}

// getLastBackupPath 获取上次备份路径
func getLastBackupPath() (string, error) {
	data, err := ioutil.ReadFile(lastBackupPathFile)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
