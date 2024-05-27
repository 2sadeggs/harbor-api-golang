package main

import (
	"encoding/json"
	"fmt"
)

// Fetch all projects
func fetchHarborProjects(baseURL, auth string) ([]Project, error) {
	var projects []Project
	page := 1

	for {
		url := fmt.Sprintf("%s/projects?page=%d&page_size=%d", baseURL, page, MaxPageSize)
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
func fetchProjectRepositories(baseURL, projectName, auth string) ([]Repository, error) {
	var repositories []Repository
	page := 1

	for {
		url := fmt.Sprintf("%s/projects/%s/repositories?page=%d&page_size=%d", baseURL, projectName, page, MaxPageSize)
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
func fetchProjectRepositoryArtifacts(baseURL, projectName, repositoryName, auth string) ([]Artifact, error) {
	var artifacts []Artifact
	page := 1

	for {
		url := fmt.Sprintf("%s/projects/%s/repositories/%s/artifacts?page=%d&page_size=%d", baseURL, projectName, repositoryName, page, MaxPageSize)
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
