package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	DefaultPageSize = 10
	MaxPageSize     = 100
)

func main() {
	// 从环境变量中获取 harbor host 和认证信息
	baseURL := os.Getenv("HARBOR_BASEURL")
	auth := os.Getenv("HARBOR_AUTH")
	if baseURL == "" || auth == "" {
		fmt.Println("Error: HARBOR_BASEURL or HARBOR_AUTH environment variables are not set.")
		return
	}

	// 定义命令行选项
	action := flag.String("action", "", "Action to perform:"+
		"ping , health , statistics , projects , repositories , artifacts , uris , download , save")
	flag.Parse()

	switch *action {
	case "ping":
		// 检查 Harbor 是否可用
		isAlive, err := CheckHarborPing(baseURL)
		if err != nil {
			fmt.Println("Error checking Harbor availability:", err)
			return
		}
		if isAlive {
			fmt.Println("Harbor is alive.")
		} else {
			fmt.Println("Harbor is not alive.")
		}
	case "health":
		// 检查 Harbor API 健康状态
		isHealthy, err := CheckHarborHealth(baseURL)
		if err != nil {
			fmt.Println("Error checking Harbor health:", err)
			return
		}
		if isHealthy {
			fmt.Println("Harbor API is healthy.")
		} else {
			fmt.Println("Harbor API is not healthy.")
		}
	case "statistics":
		// 获取 Harbor 统计信息
		stats, err := GetHarborStatistics(baseURL, auth)
		if err != nil {
			fmt.Println("Error getting Harbor statistics:", err)
			return
		}
		// 输出统计信息
		PrintHarborStatistics(stats)
	case "projects":
		projects, err := fetchAllProjects(baseURL, auth)
		if err != nil {
			fmt.Printf("Error fetching projects: %v\n", err)
			return
		}
		printProjects(projects)
	case "repositories":
		repositories, err := fetchAllRepositories(baseURL, auth)
		if err != nil {
			fmt.Printf("Error fetching repositories: %v\n", err)
			return
		}
		printRepositories(repositories)
	case "artifacts":
		artifacts, err := fetchAllArtifacts(baseURL, auth)
		if err != nil {
			fmt.Printf("Error fetching artifacts: %v\n", err)
			return
		}
		printArtifacts(artifacts)
	case "uris":
		// 获取所有 URI 列表
		singleArchURIs, multiArchURIs, multiArchWithChildURIs, allURIs, nonUnknownArchURIs, unknownArchURIs, err := fetchAllURIs(baseURL, auth)
		if err != nil {
			fmt.Printf("Error fetching URIs: %v\n", err)
			return
		}
		// 打印所有 URI 列表
		printAllURIs(singleArchURIs, multiArchURIs, multiArchWithChildURIs, allURIs, nonUnknownArchURIs, unknownArchURIs)
	case "download":
		err := downloadArtifacts(baseURL, auth)
		if err != nil {
			fmt.Printf("Error downloading and saving artifacts: %v\n", err)
			return
		}
	case "save":
		err := downloadAndSaveArtifacts(baseURL, auth)
		if err != nil {
			fmt.Printf("Error downloading and saving artifacts: %v\n", err)
		}
	default:
		fmt.Println("Invalid action. Please choose one of: " +
			"ping , health , statistics , projects , repositories , artifacts , uris , " +
			"download , save")
	}
}
