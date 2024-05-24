package main

import (
	"fmt"
)

// Fetch all artifact URIs for all repositories
func fetchAllArtifactURIsV3(scheme, harborHost, auth string) (NonUnknownArchURIs []string, err error) {

	baseURL := fmt.Sprintf("%s://%s/api/v2.0", scheme, harborHost)

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

	for _, artifact := range artifacts {
		projectName := getProjectNameByID(artifact.ProjectID, projects)
		if projectName == "" {
			return nil, fmt.Errorf("project name not found for project ID: %d", artifact.ProjectID)
		}

		repoName := getRepoNameByID(artifact.RepositoryID, repositories)
		if repoName == "" {
			return nil, fmt.Errorf("repository name not found for repository ID: %d", artifact.RepositoryID)
		}

		var NonUnknownArchURI string
		if len(artifact.References) == 0 {
			// Single-architecture artifact
			NonUnknownArchURI = fmt.Sprintf("%s/%s/%s@%s", harborHost, projectName, repoName, artifact.Digest)
			NonUnknownArchURIs = append(NonUnknownArchURIs, NonUnknownArchURI)
		} else {
			// Multi-architecture artifact
			for _, reference := range artifact.References {
				if reference.Platform.Architecture == "unknown" || reference.Platform.Os == "unknown" {
					// 如果当前制品引用 reference 的平台架构或平台OS为 unknown 则跳出该 reference 循环 进入下一个 reference
					continue
				}
				NonUnknownArchURI = fmt.Sprintf("%s/%s/%s@%s", harborHost, projectName, repoName, reference.ChildDigest)
				NonUnknownArchURIs = append(NonUnknownArchURIs, NonUnknownArchURI)
			}
		}
	}

	return NonUnknownArchURIs, nil
}
