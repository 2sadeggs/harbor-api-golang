package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	DefaultPageSize = 10
	MaxPageSize     = 100
)

// Data structures (omitted for brevity, use the ones from the initial example)

// Utility function to make GET requests and return the response body
func getRequest(url, auth string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Basic "+auth)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data, status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Check API availability
func checkAPI(apiEndpoint, auth string) error {
	url := fmt.Sprintf("%s/ping", apiEndpoint)
	_, err := getRequest(url, auth)
	return err
}

// Fetch statistics
func fetchStatistics(apiEndpoint, auth string) (*HarborStatistics, error) {
	url := fmt.Sprintf("%s/statistics", apiEndpoint)
	body, err := getRequest(url, auth)
	if err != nil {
		return nil, err
	}

	var stats HarborStatistics
	err = json.Unmarshal(body, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

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

		if len(pageProjects) == 0 {
			break
		}

		projects = append(projects, pageProjects...)
		page++
	}

	return projects, nil
}

// Fetch repositories for a project
func fetchRepositories(apiEndpoint, projectName, auth string) ([]Repository, error) {
	var repositories []Repository
	page := 1

	for {
		url := fmt.Sprintf("%s/projects/%s/repositories?page=%d&page_size=%d", apiEndpoint, projectName, page, MaxPageSize)
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

	return repositories, nil
}

// Fetch artifacts for a repository
func fetchArtifacts(apiEndpoint, projectName, repositoryName, auth string) ([]Artifact, error) {
	var artifacts []Artifact
	page := 1

	for {
		url := fmt.Sprintf("%s/projects/%s/repositories/%s/artifacts?page=%d&page_size=%d", apiEndpoint, projectName, repositoryName, page, MaxPageSize)
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

	return artifacts, nil
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

	// Check API availability
	err := checkAPI(apiEndpoint, auth)
	if err != nil {
		fmt.Printf("API not available: %v\n", err)
		return
	}

	// Fetch statistics
	stats, err := fetchStatistics(apiEndpoint, auth)
	if err != nil {
		fmt.Printf("Error fetching statistics: %v\n", err)
		return
	}

	fmt.Printf("Statistics: %+v\n", stats)

	// Fetch all projects
	projects, err := fetchAllProjects(apiEndpoint, auth)
	if err != nil {
		fmt.Printf("Error fetching projects: %v\n", err)
		return
	}

	for _, project := range projects {
		fmt.Printf("Project: %+v\n", project)

		// Fetch repositories for the project
		repositories, err := fetchRepositories(apiEndpoint, project.Name, auth)
		if err != nil {
			fmt.Printf("Error fetching repositories for project %s: %v\n", project.Name, err)
			continue
		}

		for _, repository := range repositories {
			fmt.Printf("  Repository: %+v\n", repository)

			// Fetch artifacts for the repository
			artifacts, err := fetchArtifacts(apiEndpoint, project.Name, repository.Name, auth)
			if err != nil {
				fmt.Printf("Error fetching artifacts for repository %s: %v\n", repository.Name, err)
				continue
			}

			for _, artifact := range artifacts {
				fmt.Printf("    Artifact: %+v\n", artifact)

				// Process sub-artifacts if any (multi-architecture images)
				if len(artifact.References) > 0 {
					for _, reference := range artifact.References {
						fmt.Printf("      Sub-Artifact: %+v\n", reference)
					}
				}
			}
		}
	}
}
