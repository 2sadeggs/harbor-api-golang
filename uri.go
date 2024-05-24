package main

import (
	"fmt"
	"net/url"
)

// Fetch all artifact URIs for all repositories
func fetchAllArtifactURIs(baseURL, auth string) ([]string, error) {
	// Parse baseURL to extract harborHost
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid baseURL: %v", err)
	}

	harborHost := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	artifacts, err := fetchAllArtifacts(baseURL, auth)
	if err != nil {
		return nil, err
	}

	repositories, err := fetchAllRepositories(baseURL, auth)
	if err != nil {
		return nil, err
	}

	projects, err := fetchAllProjects(baseURL, auth)
	if err != nil {
		return nil, err
	}

	var uris []string

	for _, artifact := range artifacts {
		projectName := getProjectNameByID(artifact.ProjectID, projects)
		if projectName == "" {
			return nil, fmt.Errorf("project name not found for project ID: %d", artifact.ProjectID)
		}

		repoName := getRepoNameByID(artifact.RepositoryID, repositories)
		if repoName == "" {
			return nil, fmt.Errorf("repository name not found for repository ID: %d", artifact.RepositoryID)
		}

		var uri string
		if len(artifact.References) == 0 {
			// Single-architecture artifact
			uri = fmt.Sprintf("%s/%s/%s@%s",
				baseURL, projectName, repoName, artifact.Digest)
			uris = append(uris, uri)
		} else {
			// Multi-architecture artifact
			for _, reference := range artifact.References {
				if reference.Platform.Architecture == "unknown" || reference.Platform.Os == "unknown" {
					// 如果当前制品引用 reference 的平台架构或平台OS为 unknown 则跳出该 reference 循环 进入下一个 reference
					continue
				}
				uri = fmt.Sprintf("%s/%s/%s@%s::%s",
					baseURL, projectName, repoName, artifact.Digest, reference.ChildDigest)
				uris = append(uris, uri)
			}
		}

	}

	return uris, nil
}
