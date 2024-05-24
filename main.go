package main

import (
	"fmt"
	"os"
)

const (
	DefaultPageSize = 10
	MaxPageSize     = 100
)

func main() {
	// 从环境变量中获取 harbor host 和认证信息
	scheme := os.Getenv("HARBOR_SCHEME")
	harborHost := os.Getenv("HARBOR_HOST")
	auth := os.Getenv("HARBOR_AUTH")
	if scheme == "" || harborHost == "" || auth == "" {
		fmt.Println("Error: HARBOR_SCHEME, HARBOR_HOST or HARBOR_AUTH environment variables are not set.")
		return
	}

	//baseURL := fmt.Sprintf("%s://%s/api/v2.0", scheme, harborHost)

	URIs, err := fetchAllArtifactURIsV3(scheme, harborHost, auth)
	if err != nil {
		fmt.Printf("Error fetching artifact URIs: %v\n", err)
		return
	}

	fmt.Println("Artifact URIs:")
	for _, uri := range URIs {
		fmt.Println(uri)
	}

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
