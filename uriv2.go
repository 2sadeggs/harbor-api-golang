package main

import (
	"fmt"
	"net/url"
)

func fetchAllArtifactsWithTypes(baseURL, auth string) (map[string][]string, error) {
	// 解析 baseURL 以提取 harborHost
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid baseURL: %v", err)
	}
	harborHost := fmt.Sprintf("%s", u.Host)

	repositories, err := fetchAllRepositories(baseURL, auth)
	if err != nil {
		return nil, err
	}

	artifacts, err := fetchAllArtifacts(baseURL, auth)
	if err != nil {
		return nil, err
	}

	// 初始化存储 URI 列表的 map
	uriMap := map[string][]string{
		"single_architecture":   {},
		"multi_architecture":    {},
		"multi_arch_with_child": {},
		"all_uris":              {},
		"non_unknown_arch_uris": {},
		"unknown_arch_uris":     {},
	}

	// 遍历所有制品
	for _, artifact := range artifacts {
		// 根据制品的 RepositoryID 获取 RepositoryName
		repoName := getRepoNameByID(artifact.RepositoryID, repositories)
		if repoName == "" {
			return nil, fmt.Errorf("repository name not found for repository ID: %d", artifact.RepositoryID)
		}

		if len(artifact.References) == 0 {
			// 单架构制品
			uri := fmt.Sprintf("%s/%s@%s", harborHost, repoName, artifact.Digest)
			uriMap["single_architecture"] = append(uriMap["single_architecture"], uri)
			uriMap["all_uris"] = append(uriMap["all_uris"], uri)
			uriMap["non_unknown_arch_uris"] = append(uriMap["non_unknown_arch_uris"], uri)
		} else {
			// 多架构制品
			uri := fmt.Sprintf("%s/%s@%s", harborHost, repoName, artifact.Digest)
			uriMap["multi_architecture"] = append(uriMap["multi_architecture"], uri)

			for _, reference := range artifact.References {
				childURI := fmt.Sprintf("%s/%s@%s::%s", harborHost, repoName, artifact.Digest, reference.ChildDigest)
				uriMap["multi_arch_with_child"] = append(uriMap["multi_arch_with_child"], childURI)

				childDigestURI := fmt.Sprintf("%s/%s@%s", harborHost, repoName, reference.ChildDigest)

				if reference.Platform.Architecture != "unknown" && reference.Platform.Os != "unknown" {
					uriMap["non_unknown_arch_uris"] = append(uriMap["non_unknown_arch_uris"], childDigestURI)
				} else {
					uriMap["unknown_arch_uris"] = append(uriMap["unknown_arch_uris"], childDigestURI)
				}

				uriMap["all_uris"] = append(uriMap["all_uris"], childDigestURI)
			}
		}
	}

	// 返回存储 URI 列表的 map 和错误信息
	return uriMap, nil
}
