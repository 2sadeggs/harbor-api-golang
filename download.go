package main

import (
	"fmt"
	"os/exec"
)

// 用于下载并保存所有制品
func downloadAndSaveArtifacts(baseURL, auth string) error {
	//_, _, _, uris, _, _, err := fetchAllURIs(baseURL, auth)
	uris, err := fetchNonUnknownArchURIs(baseURL, auth)
	if err != nil {
		return err
	}

	for _, uri := range uris {
		fmt.Printf("Downloading artifact: %s\n", uri)
		cmd := exec.Command("docker", "pull", uri)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to download artifact %s: %v\nOutput: %s", uri, err, string(output))
		}
		fmt.Printf("Successfully downloaded artifact: %s\n", uri)
	}

	return nil
}
