package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// CheckHarborPing 检查 Harbor 是否可用，只有返回 "pong" 时才表示 Harbor 存活
// Ping Harbor to check if it's alive.
// This API simply replies a pong to indicate the process to handle API is up, disregarding the health status of dependent components.
func CheckHarborPing(baseURL string) (bool, error) {
	// 路径
	path := "/ping"

	// 构造请求
	req, err := http.NewRequest("GET", baseURL+path, nil)
	if err != nil {
		return false, fmt.Errorf("error creating request: %v", err)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("error reading response: %v", err)
	}

	// 输出调试信息
	fmt.Println("Ping response:", string(body))

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", resp.StatusCode)
		return false, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	// 判断是否返回 Pong
	if string(body) != "Pong" {
		fmt.Println("Harbor is not alive.")
		return false, nil
	}

	return true, nil
}
