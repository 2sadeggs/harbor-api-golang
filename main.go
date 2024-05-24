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
	action := flag.String("action", "", "Action to perform: ping, health, statistics, projects, repositories, artifacts, non-unknown-artifacts")
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
	case "alltypeuris":
		NonUnknownArchUris, singles, multies, multiesWithChilds, multiesWithChildsAndUnknownArchs, err :=
			fetchAllArtifactURIs(baseURL, auth)
		if err != nil {
			fmt.Printf("Error fetching artifact URIs: %v\n", err)
			return
		}

		fmt.Println("NonUnknownArch URIs:")
		for _, NonUnknownArchUri := range NonUnknownArchUris {
			fmt.Println(NonUnknownArchUri)
		}

		fmt.Println("Single URIs:")
		for _, single := range singles {
			fmt.Println(single)
		}

		fmt.Println("Muities URIs:")
		for _, multi := range multies {
			fmt.Println(multi)
		}
		fmt.Println("MultiesWithChilds URIs:")
		for _, multiesWithChild := range multiesWithChilds {
			fmt.Println(multiesWithChild)
		}
		fmt.Println("MultiesWithChildsAndUnknownArch URIs:")
		for _, multiesWithChildsAndUnknownArch := range multiesWithChildsAndUnknownArchs {
			fmt.Println(multiesWithChildsAndUnknownArch)
		}
	case "formaturis":
		URIs, err := fetchAllArtifactURIsNonUnknownArch(baseURL, auth)
		if err != nil {
			fmt.Printf("Error fetching artifact URIs: %v\n", err)
			return
		}

		fmt.Println("Artifact URIs:")
		for _, uri := range URIs {
			fmt.Println(uri)
		}
	default:
		fmt.Println("Invalid action. Please choose one of: ping, health, statistics, projects, repositories, artifacts, non-unknown-artifacts")
	}

	/*	URIs, err := fetchAllArtifactURIsNonUnknownArch(baseURL, auth)
		if err != nil {
			fmt.Printf("Error fetching artifact URIs: %v\n", err)
			return
		}

		fmt.Println("Artifact URIs:")
		for _, uri := range URIs {
			fmt.Println(uri)
		}*/

	/*	URIs, err := fetchAllArtifactURIsV3(scheme, harborHost, auth)
		if err != nil {
			fmt.Printf("Error fetching artifact URIs: %v\n", err)
			return
		}

		fmt.Println("Artifact URIs:")
		for _, uri := range URIs {
			fmt.Println(uri)
		}*/

	/*
		NonUnknownUris, singles, multies, multiesWithchilds, multiesWithchildsAndUnknownArchs, err := fetchAllArtifactURIsV2(scheme, harborHost, auth)
		if err != nil {
			fmt.Printf("Error fetching artifact URIs: %v\n", err)
			return
		}

		fmt.Println("NonUnknown URIs:")
		for _, NonUnknownUri := range NonUnknownUris {
			fmt.Println(NonUnknownUri)
		}

		fmt.Println("Single URIs:")
		for _, single := range singles {
			fmt.Println(single)
		}

		fmt.Println("Muities URIs:")
		for _, multi := range multies {
			fmt.Println(multi)
		}
		fmt.Println("MultiesWithchilds URIs:")
		for _, multiesWithchild := range multiesWithchilds {
			fmt.Println(multiesWithchild)
		}
		fmt.Println("MultiesWithchildsAndUnknownArch URIs:")
		for _, multiesWithchildsAndUnknownArch := range multiesWithchildsAndUnknownArchs {
			fmt.Println(multiesWithchildsAndUnknownArch)
		}
	*/
	/*
		artifactURIs, err := fetchAllArtifactURIs(baseURL, auth)
		if err != nil {
			fmt.Printf("Error fetching artifact URIs: %v\n", err)
			return
		}

		fmt.Println("Artifact URIs:")
		for _, uri := range artifactURIs {
			fmt.Println(uri)
		}
	*/
	/*artifactDigests, err := fetchAllArtifactsWithChildDigests(baseURL, auth)
	if err != nil {
		fmt.Printf("Error fetching artifact digests: %v\n", err)
		return
	}

	fmt.Println("Single-architecture artifact digests:")
	for _, digest := range artifactDigests["single_architecture"].([]string) {
		fmt.Printf("Digest: %s\n", digest)
	}

	fmt.Println("Multi-architecture artifact digests:")
	for _, multiArch := range artifactDigests["multi_architecture"].([]map[string]interface{}) {
		fmt.Printf("Digest: %s\n", multiArch["digest"])
		fmt.Println("Child Digests:")
		for _, childDigest := range multiArch["childDigests"].([]string) {
			fmt.Printf("  - %s\n", childDigest)
		}
	}*/

	/*artifactDigests, err := fetchAllArtifactsWithTypes(baseURL, auth)
	if err != nil {
		fmt.Printf("Error fetching artifact digests: %v\n", err)
		return
	}

	fmt.Println("Single-architecture artifact digests:")
	for _, digest := range artifactDigests["single_architecture"] {
		fmt.Printf("Digest: %s\n", digest)
	}

	fmt.Println("Multi-architecture artifact digests:")
	for _, digest := range artifactDigests["multi_architecture"] {
		fmt.Printf("Digest: %s\n", digest)
	}*/

	/*	// Fetch all projects
		projects, err := fetchAllProjects(baseURL, auth)
		if err != nil {
			fmt.Printf("Error fetching projects: %v\n", err)
			return
		}
		printProjects(projects)*/

	/*	// Fetch all repositories
		repositories, err := fetchAllRepositories(baseURL, auth)
		if err != nil {
			fmt.Printf("Error fetching repositories: %v\n", err)
			return
		}
		printRepositories(repositories)*/

	/*	// Fetch all artifacts
		artifacts, err := fetchAllArtifacts(baseURL, auth)
		if err != nil {
			fmt.Printf("Error fetching artifacts: %v\n", err)
			return
		}
		printArtifacts(artifacts)*/
}
