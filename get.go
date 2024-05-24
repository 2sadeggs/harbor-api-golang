package main

// Function to get repository name by repository ID
func getRepoNameByID(repoID int, repositories []Repository) string {
	for _, repo := range repositories {
		if repo.ID == repoID {
			return repo.Name
		}
	}
	return ""
}

// Function to get project name by project ID
func getProjectNameByID(projectID int, projects []Project) string {
	for _, project := range projects {
		if project.ProjectID == projectID {
			return project.Name
		}
	}
	return ""
}
