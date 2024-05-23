package main

import (
	"fmt"
	"net/url"
	"strings"
)

// Helper function to extract and double encode repository name
func extractAndDoubleEncodeRepoName(fullRepoName string) (string, error) {
	repoNameParts := strings.SplitN(fullRepoName, "/", 2)
	if len(repoNameParts) != 2 {
		return "", fmt.Errorf("invalid repository name format: %s", fullRepoName)
	}
	// First URL encoding
	firstEncodedRepoName := url.PathEscape(repoNameParts[1])
	// Second URL encoding
	doubleEncodedRepoName := url.PathEscape(firstEncodedRepoName)
	return doubleEncodedRepoName, nil
}
