package main

import (
	"fmt"
)

// Fetch all artifact URIs for all repositories
func fetchAllArtifactURIsV2(scheme, harborHost, auth string) (
	NonUnknownUris, singleArchUris, multiArchUris,
	multiArchUrisWithChildDigest, multiArchUrisWithChildDigestAndUnknownArch []string, err error) {

	baseURL := fmt.Sprintf("%s://%s/api/v2.0", scheme, harborHost)

	artifacts, err := fetchAllArtifacts(baseURL, auth)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	repositories, err := fetchAllRepositories(baseURL, auth)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	projects, err := fetchAllProjects(baseURL, auth)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	for _, artifact := range artifacts {
		projectName := getProjectNameByID(artifact.ProjectID, projects)
		if projectName == "" {
			return nil, nil, nil, nil, nil, fmt.Errorf("project name not found for project ID: %d", artifact.ProjectID)
		}

		repoName := getRepoNameByID(artifact.RepositoryID, repositories)
		if repoName == "" {
			return nil, nil, nil, nil, nil, fmt.Errorf("repository name not found for repository ID: %d", artifact.RepositoryID)
		}

		var NonUnknown string
		if len(artifact.References) == 0 {
			// Single-architecture artifact
			uri := fmt.Sprintf("%s/%s/%s@%s",
				harborHost, projectName, repoName, artifact.Digest)
			//uris = append(uris, uri)
			NonUnknown = fmt.Sprintf("%s/%s/%s@%s",
				harborHost, projectName, repoName, artifact.Digest)
			singleArchUris = append(singleArchUris, uri)
			NonUnknownUris = append(NonUnknownUris, NonUnknown)
		} else {
			// Multi-architecture artifact
			uri := fmt.Sprintf("%s/%s/%s@%s",
				harborHost, projectName, repoName, artifact.Digest)
			//uris = append(uris, uri)
			multiArchUris = append(multiArchUris, uri)
			for _, reference := range artifact.References {
				uriMultiWithChild := fmt.Sprintf("%s/%s/%s@%s::%s",
					harborHost, projectName, repoName, artifact.Digest, reference.ChildDigest)
				//uris = append(uris, uri)
				multiArchUrisWithChildDigest = append(multiArchUrisWithChildDigest, uriMultiWithChild)
				if reference.Platform.Architecture == "unknown" || reference.Platform.Os == "unknown" {
					// 如果当前制品引用 reference 的平台架构或平台OS为 unknown 则跳出该 reference 循环 进入下一个 reference
					uriMultiWithChildAndUnknownArch := fmt.Sprintf("%s/%s/%s@%s::%s",
						harborHost, projectName, repoName, artifact.Digest, reference.ChildDigest)
					multiArchUrisWithChildDigestAndUnknownArch = append(multiArchUrisWithChildDigestAndUnknownArch, uriMultiWithChildAndUnknownArch)
					continue
				}
				NonUnknown = fmt.Sprintf("%s/%s/%s@%s",
					harborHost, projectName, repoName, reference.ChildDigest)
				NonUnknownUris = append(NonUnknownUris, NonUnknown)
			}
		}
	}

	return NonUnknownUris, singleArchUris, multiArchUris, multiArchUrisWithChildDigest, multiArchUrisWithChildDigestAndUnknownArch, nil
}
