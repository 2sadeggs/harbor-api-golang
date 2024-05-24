package main

import (
	"fmt"
	"net/url"
)

func fetchAllURIs(baseURL, auth string) (
	singleArchURIs []string,
	multiArchURIs []string,
	multiArchWithChildURIs []string,
	allURIs []string,
	nonUnknownArchURIs []string,
	unknownArchURIs []string,
	err error,
) {
	// 解析 baseURL 以提取 harborHost
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("invalid baseURL: %v", err)
	}
	harborHost := fmt.Sprintf("%s", u.Host)

	repositories, err := fetchAllRepositories(baseURL, auth)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	artifacts, err := fetchAllArtifacts(baseURL, auth)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	// 遍历所有制品
	for _, artifact := range artifacts {
		// 根据制品的 RepositoryID 获取 RepositoryName
		repoName := getRepoNameByID(artifact.RepositoryID, repositories)
		if repoName == "" {
			return nil, nil, nil, nil, nil, nil, fmt.Errorf("repository name not found for repository ID: %d", artifact.RepositoryID)
		}

		if len(artifact.References) == 0 {
			// 单架构制品
			uri := fmt.Sprintf("%s/%s@%s", harborHost, repoName, artifact.Digest)
			singleArchURIs = append(singleArchURIs, uri)
			allURIs = append(allURIs, uri)
			nonUnknownArchURIs = append(nonUnknownArchURIs, uri)
		} else {
			// 多架构制品
			uri := fmt.Sprintf("%s/%s@%s", harborHost, repoName, artifact.Digest)
			multiArchURIs = append(multiArchURIs, uri)

			for _, reference := range artifact.References {
				childURI := fmt.Sprintf("%s/%s@%s::%s", harborHost, repoName, artifact.Digest, reference.ChildDigest)
				multiArchWithChildURIs = append(multiArchWithChildURIs, childURI)

				childDigestURI := fmt.Sprintf("%s/%s@%s", harborHost, repoName, reference.ChildDigest)

				if reference.Platform.Architecture != "unknown" && reference.Platform.Os != "unknown" {
					nonUnknownArchURIs = append(nonUnknownArchURIs, childDigestURI)
				} else {
					unknownArchURIs = append(unknownArchURIs, childDigestURI)
				}

				allURIs = append(allURIs, childDigestURI)
			}
		}
	}

	// 返回各个类型的 URI 列表和错误信息
	return singleArchURIs, multiArchURIs, multiArchWithChildURIs, allURIs, nonUnknownArchURIs, unknownArchURIs, nil
}

func fetchSingleArchURIs(baseURL, auth string) ([]string, error) {
	singleArchURIs, _, _, _, _, _, err := fetchAllURIs(baseURL, auth)
	if err != nil {
		return nil, err
	}
	return singleArchURIs, nil
}

func fetchMultiArchURIs(baseURL, auth string) ([]string, error) {
	_, multiArchURIs, _, _, _, _, err := fetchAllURIs(baseURL, auth)
	if err != nil {
		return nil, err
	}
	return multiArchURIs, nil
}

func fetchMultiArchWithChildURIs(baseURL, auth string) ([]string, error) {
	_, _, multiArchWithChildURIs, _, _, _, err := fetchAllURIs(baseURL, auth)
	if err != nil {
		return nil, err
	}
	return multiArchWithChildURIs, nil
}

func fetchAllURIsList(baseURL, auth string) ([]string, error) {
	_, _, _, allURIs, _, _, err := fetchAllURIs(baseURL, auth)
	if err != nil {
		return nil, err
	}
	return allURIs, nil
}

func fetchNonUnknownArchURIs(baseURL, auth string) ([]string, error) {
	_, _, _, _, nonUnknownArchURIs, _, err := fetchAllURIs(baseURL, auth)
	if err != nil {
		return nil, err
	}
	return nonUnknownArchURIs, nil
}

func fetchUnknownArchURIs(baseURL, auth string) ([]string, error) {
	_, _, _, _, _, unknownArchURIs, err := fetchAllURIs(baseURL, auth)
	if err != nil {
		return nil, err
	}
	return unknownArchURIs, nil
}
