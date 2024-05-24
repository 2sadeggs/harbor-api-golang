package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Health check API
// The endpoint returns the health stauts of the system.
// CheckHarborHealth 检查 Harbor API 健康状态，只有返回状态为 "healthy" 时才表示 Harbor API 健康
func CheckHarborHealth(baseURL string) (bool, error) {
	// 路径
	path := "/health"

	// 构造请求
	req, err := http.NewRequest("GET", baseURL+path, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false, err
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return false, err
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", resp.StatusCode)
		return false, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return false, err
	}

	// 解析 JSON 响应
	var healthStatus HealthStatus
	err = json.Unmarshal(body, &healthStatus)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return false, err
	}

	// 输出健康状态
	fmt.Println("Harbor API Health Status:")
	fmt.Println("Overall Status:", healthStatus.Status)
	fmt.Println("Components:")
	for _, comp := range healthStatus.Components {
		fmt.Printf("  %s: %s\n", comp.Name, comp.Status)
	}

	// 检查健康状态是否为 "healthy"
	if healthStatus.Status != "healthy" {
		fmt.Println("Harbor API is not healthy.")
		return false, nil
	}

	fmt.Println("Harbor API is healthy.")
	return true, nil
}
