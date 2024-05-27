package main

import "fmt"

// Print projects in a formatted way
func printProjects(projects []Project) {
	fmt.Println("Projects:")
	for _, project := range projects {
		fmt.Printf("  - Name: %s\n", project.Name)
		fmt.Printf("    ID: %d\n", project.ProjectID)
		fmt.Printf("    Owner: %s\n", project.OwnerName)
		fmt.Printf("    Repo Count: %d\n", project.RepoCount)
		fmt.Printf("    Creation Time: %s\n", project.CreationTime)
		fmt.Println()
	}
}

// Print repositories in a formatted way
func printRepositories(repositories []Repository) {
	fmt.Println("Repositories:")
	for _, repository := range repositories {
		fmt.Printf("  - Name: %s\n", repository.Name)
		fmt.Printf("    ID: %d\n", repository.ID)
		fmt.Printf("    Project ID: %d\n", repository.ProjectID)
		fmt.Printf("    Artifact Count: %d\n", repository.ArtifactCount)
		fmt.Printf("    Pull Count: %d\n", repository.PullCount)
		fmt.Printf("    Creation Time: %s\n", repository.CreationTime)
		fmt.Println()
	}
}

// Print artifacts in a formatted way
func printArtifacts(artifacts []Artifact) {
	fmt.Println("Artifacts:")
	for _, artifact := range artifacts {
		fmt.Printf("  - Digest: %s\n", artifact.Digest)
		fmt.Printf("    ID: %d\n", artifact.ID)
		fmt.Printf("    Project ID: %d\n", artifact.ProjectID)
		fmt.Printf("    Repository ID: %d\n", artifact.RepositoryID)
		fmt.Printf("    Media Type: %s\n", artifact.MediaType)
		fmt.Printf("    Size: %d\n", artifact.Size)
		fmt.Printf("    Push Time: %s\n", artifact.PushTime)
		fmt.Println()
	}
}

// 打印所有制品的 URI 列表
func printAllURIs(singleArchURIs, multiArchURIs, multiArchWithChildURIs, allURIs, nonUnknownArchURIs, unknownArchURIs []string) {
	fmt.Println("Single-architecture URIs:")
	for _, uri := range singleArchURIs {
		fmt.Println(uri)
	}

	fmt.Println("\nMulti-architecture URIs:")
	for _, uri := range multiArchURIs {
		fmt.Println(uri)
	}

	fmt.Println("\nMulti-architecture with child URIs:")
	for _, uri := range multiArchWithChildURIs {
		fmt.Println(uri)
	}

	fmt.Println("\nAll URIs:")
	for _, uri := range allURIs {
		fmt.Println(uri)
	}

	fmt.Println("\nNon-unknown architecture and OS URIs:")
	for _, uri := range nonUnknownArchURIs {
		fmt.Println(uri)
	}

	fmt.Println("\nUnknown architecture or OS URIs:")
	for _, uri := range unknownArchURIs {
		fmt.Println(uri)
	}
}

// 打印 URI 列表的函数
func printURIs(label string, uris []string) {
	fmt.Printf("%s:\n", label)
	for _, uri := range uris {
		fmt.Println(uri)
	}
	fmt.Println()
}

// 打印 URI
func printArtifactsWithTypes(artifactMap map[string][]string) {
	fmt.Println("Artifacts with types:")
	for key, uris := range artifactMap {
		fmt.Printf("%s:\n", key)
		for _, uri := range uris {
			fmt.Printf("- %s\n", uri)
		}
	}
}
