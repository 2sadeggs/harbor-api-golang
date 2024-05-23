package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
)

const (
	DefaultPageSize = 10
	MaxPageSize     = 100
)

// Fetch all projects
func fetchAllProjects(apiEndpoint, auth string) ([]Project, error) {
	var projects []Project
	page := 1

	for {
		url := fmt.Sprintf("%s/projects?page=%d&page_size=%d", apiEndpoint, page, MaxPageSize)
		body, err := getRequest(url, auth)
		if err != nil {
			return nil, err
		}

		var pageProjects []Project
		err = json.Unmarshal(body, &pageProjects)
		if err != nil {
			return nil, err
		}

		// Check if the result is empty
		if len(pageProjects) == 0 {
			break
		}

		projects = append(projects, pageProjects...)
		page++
	}

	return projects, nil
}

// Fetch all repositories for all projects
func fetchAllRepositories(apiEndpoint, auth string) ([]Repository, error) {
	var allRepositories []Repository
	projects, err := fetchAllProjects(apiEndpoint, auth)
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		var repositories []Repository
		page := 1

		for {
			url := fmt.Sprintf("%s/projects/%s/repositories?page=%d&page_size=%d", apiEndpoint, project.Name, page, MaxPageSize)
			body, err := getRequest(url, auth)
			if err != nil {
				return nil, err
			}

			var pageRepositories []Repository
			err = json.Unmarshal(body, &pageRepositories)
			if err != nil {
				return nil, err
			}

			if len(pageRepositories) == 0 {
				break
			}

			repositories = append(repositories, pageRepositories...)
			page++
		}

		allRepositories = append(allRepositories, repositories...)
	}

	return allRepositories, nil
}

// Fetch all artifacts for all repositories
func fetchAllArtifacts(apiEndpoint, auth string) ([]Artifact, error) {
	var allArtifacts []Artifact
	repositories, err := fetchAllRepositories(apiEndpoint, auth)
	if err != nil {
		return nil, err
	}

	for _, repository := range repositories {
		var artifacts []Artifact
		page := 1

		repoNameParts := strings.SplitN(repository.Name, "/", 2)
		if len(repoNameParts) != 2 {
			return nil, fmt.Errorf("invalid repository name format: %s", repository.Name)
		}
		encodedRepoName := url.PathEscape(repoNameParts[1])
		doubleEncodedRepoName := url.PathEscape(encodedRepoName)

		for {
			url := fmt.Sprintf("%s/projects/%s/repositories/%s/artifacts?page=%d&page_size=%d",
				apiEndpoint, repoNameParts[0], doubleEncodedRepoName, page, MaxPageSize)
			body, err := getRequest(url, auth)
			if err != nil {
				return nil, err
			}

			var pageArtifacts []Artifact
			err = json.Unmarshal(body, &pageArtifacts)
			if err != nil {
				return nil, err
			}

			if len(pageArtifacts) == 0 {
				break
			}

			artifacts = append(artifacts, pageArtifacts...)
			page++
		}

		allArtifacts = append(allArtifacts, artifacts...)
	}

	return allArtifacts, nil
}

// Print projects in a formatted way
func printProjects(projects []Project) {
	fmt.Println("Projects:")
	for _, project := range projects {
		fmt.Printf("  - Name: %s\n", project.Name)
		fmt.Printf("    ID: %d\n", project.ProjectID)
		fmt.Printf("    Owner: %s\n", project.OwnerName)
		fmt.Printf("    Repo Count: %d\n", project.RepoCount)
		fmt.Printf("    Creation Time: %s\n", project.CreationTime)
		fmt.Println()
	}
}

// Print repositories in a formatted way
func printRepositories(repositories []Repository) {
	fmt.Println("Repositories:")
	for _, repository := range repositories {
		fmt.Printf("  - Name: %s\n", repository.Name)
		fmt.Printf("    ID: %d\n", repository.ID)
		fmt.Printf("    Project ID: %d\n", repository.ProjectID)
		fmt.Printf("    Artifact Count: %d\n", repository.ArtifactCount)
		fmt.Printf("    Pull Count: %d\n", repository.PullCount)
		fmt.Printf("    Creation Time: %s\n", repository.CreationTime)
		fmt.Println()
	}
}

// Print artifacts in a formatted way
func printArtifacts(artifacts []Artifact) {
	fmt.Println("Artifacts:")
	for _, artifact := range artifacts {
		fmt.Printf("  - Digest: %s\n", artifact.Digest)
		fmt.Printf("    ID: %d\n", artifact.ID)
		fmt.Printf("    Project ID: %d\n", artifact.ProjectID)
		fmt.Printf("    Repository ID: %d\n", artifact.RepositoryID)
		fmt.Printf("    Media Type: %s\n", artifact.MediaType)
		fmt.Printf("    Size: %d\n", artifact.Size)
		fmt.Printf("    Push Time: %s\n", artifact.PushTime)
		fmt.Println()
	}
}

func main() {
	// 从环境变量中获取 harbor host 和认证信息
	scheme := os.Getenv("HARBOR_SCHEME")
	harborHost := os.Getenv("HARBOR_HOST")
	auth := os.Getenv("HARBOR_AUTH")
	if scheme == "" || harborHost == "" || auth == "" {
		fmt.Println("Error: HARBOR_SCHEME, HARBOR_HOST or HARBOR_AUTH environment variables are not set.")
		return
	}

	apiEndpoint := fmt.Sprintf("%s://%s/api/v2.0", scheme, harborHost)

	/*	// Fetch all projects
		projects, err := fetchAllProjects(apiEndpoint, auth)
		if err != nil {
			fmt.Printf("Error fetching projects: %v\n", err)
			return
		}
		printProjects(projects)*/

	/*	// Fetch all repositories
		repositories, err := fetchAllRepositories(apiEndpoint, auth)
		if err != nil {
			fmt.Printf("Error fetching repositories: %v\n", err)
			return
		}
		printRepositories(repositories)*/

	// Fetch all artifacts
	artifacts, err := fetchAllArtifacts(apiEndpoint, auth)
	if err != nil {
		fmt.Printf("Error fetching artifacts: %v\n", err)
		return
	}
	printArtifacts(artifacts)
}
