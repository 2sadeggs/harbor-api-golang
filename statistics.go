package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// GetHarborStatistics 获取 Harbor 统计信息
func GetHarborStatistics(baseURL, auth string) (*HarborStatistics, error) {
	// 路径
	path := "/statistics"

	// 构造请求
	req, err := http.NewRequest("GET", baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// 添加 Basic Auth 头部
	req.Header.Set("Authorization", "Basic "+auth)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	// 解析 JSON 响应
	var stats HarborStatistics
	err = json.Unmarshal(body, &stats)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %v", err)
	}

	return &stats, nil
}

// 打印Harbor统计信息
func PrintHarborStatistics(stats *HarborStatistics) {
	fmt.Println("Harbor Statistics:")
	fmt.Printf("Private Project Count: %d\n", stats.PrivateProjectCount)
	fmt.Printf("Private Repo Count: %d\n", stats.PrivateRepoCount)
	fmt.Printf("Public Project Count: %d\n", stats.PublicProjectCount)
	fmt.Printf("Public Repo Count: %d\n", stats.PublicRepoCount)
	fmt.Printf("Total Project Count: %d\n", stats.TotalProjectCount)
	fmt.Printf("Total Repo Count: %d\n", stats.TotalRepoCount)
}
