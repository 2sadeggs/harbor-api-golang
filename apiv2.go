package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// Fetch all projects
func fetchAllProjects(baseURL, auth string) ([]Project, error) {
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
func fetchAllRepositories(baseURL, auth string) ([]Repository, error) {
	var allRepositories []Repository
	projects, err := fetchAllProjects(baseURL, auth)
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		var repositories []Repository
		page := 1

		for {
			url := fmt.Sprintf("%s/projects/%s/repositories?page=%d&page_size=%d", baseURL, project.Name, page, MaxPageSize)
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
func fetchAllArtifacts(baseURL, auth string) ([]Artifact, error) {
	var allArtifacts []Artifact
	repositories, err := fetchAllRepositories(baseURL, auth)
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
				baseURL, repoNameParts[0], doubleEncodedRepoName, page, MaxPageSize)
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

// Fetch all single-architecture and multi-architecture artifact digests
func fetchAllArtifactsWithTypes(apiEndpoint, auth string) (map[string][]string, error) {
	artifacts, err := fetchAllArtifacts(apiEndpoint, auth)
	if err != nil {
		return nil, err
	}

	artifactDigests := map[string][]string{
		"single_architecture": {},
		"multi_architecture":  {},
	}

	for _, artifact := range artifacts {
		if len(artifact.References) == 0 {
			artifactDigests["single_architecture"] = append(artifactDigests["single_architecture"], artifact.Digest)
		} else {
			artifactDigests["multi_architecture"] = append(artifactDigests["multi_architecture"], artifact.Digest)
		}
	}

	return artifactDigests, nil
}

// Fetch all single-architecture and multi-architecture artifact digests, including child digests for multi-architecture artifacts
func fetchAllArtifactsWithChildDigests(baseURL, auth string) (map[string]interface{}, error) {
	artifacts, err := fetchAllArtifacts(baseURL, auth)
	if err != nil {
		return nil, err
	}

	artifactDigests := map[string]interface{}{
		"single_architecture": []string{},
		"multi_architecture":  []map[string]interface{}{},
	}

	for _, artifact := range artifacts {
		if len(artifact.References) == 0 {
			artifactDigests["single_architecture"] = append(artifactDigests["single_architecture"].([]string), artifact.Digest)
		} else {
			childDigests := []string{}
			for _, reference := range artifact.References {
				childDigests = append(childDigests, reference.ChildDigest)
			}
			multiArchEntry := map[string]interface{}{
				"digest":       artifact.Digest,
				"childDigests": childDigests,
			}
			artifactDigests["multi_architecture"] = append(artifactDigests["multi_architecture"].([]map[string]interface{}), multiArchEntry)
		}
	}

	return artifactDigests, nil
}
