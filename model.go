package main

// Data structures (omitted for brevity, use the ones from the initial example)

type HealthStatus struct {
	Status     string       `json:"status"`
	Components []HealthComp `json:"components"`
}

type HealthComp struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// HarborStatistics 结构用于解析 Harbor 统计信息的 JSON 数据
type HarborStatistics struct {
	PrivateProjectCount int `json:"private_project_count"`
	PrivateRepoCount    int `json:"private_repo_count"`
	PublicProjectCount  int `json:"public_project_count"`
	PublicRepoCount     int `json:"public_repo_count"`
	TotalProjectCount   int `json:"total_project_count"`
	TotalRepoCount      int `json:"total_repo_count"`
}

type ProjectSummary struct {
	ProjectAdminCount int `json:"project_admin_count"`
	Quota             struct {
		Hard struct {
			Storage int `json:"storage"`
		} `json:"hard"`
		Used struct {
			Storage int `json:"storage"`
		} `json:"used"`
	} `json:"quota"`
	RepoCount int `json:"repo_count"`
}

// Project 结构用于解析项目的 JSON 数据
type Project struct {
	CreationTime       string `json:"creation_time"`
	CurrentUserRoleID  int    `json:"current_user_role_id"`
	CurrentUserRoleIDs []int  `json:"current_user_role_ids"`
	CVEAllowlist       struct {
		CreationTime string `json:"creation_time"`
		ID           int    `json:"id"`
		Items        []struct {
			ID        int    `json:"id"`
			ProjectID int    `json:"project_id"`
			Severity  string `json:"severity"`
			VulName   string `json:"vul_name"`
		} `json:"items"`
		ProjectID  int    `json:"project_id"`
		UpdateTime string `json:"update_time"`
	} `json:"cve_allowlist"`
	Metadata struct {
		Public string `json:"public"`
	} `json:"metadata"`
	Name       string `json:"name"`
	OwnerID    int    `json:"owner_id"`
	OwnerName  string `json:"owner_name"`
	ProjectID  int    `json:"project_id"`
	RepoCount  int    `json:"repo_count"`
	UpdateTime string `json:"update_time"`
}

type Repository struct {
	UpdateTime    string `json:"update_time"`
	Description   string `json:"description"`
	PullCount     int    `json:"pull_count"`
	CreationTime  string `json:"creation_time"`
	ArtifactCount int    `json:"artifact_count"`
	ProjectID     int    `json:"project_id"`
	ID            int    `json:"id"`
	Name          string `json:"name"`
}

type Artifact struct {
	Size              int         `json:"size"`
	PushTime          string      `json:"push_time"`
	ScanOverview      interface{} `json:"scan_overview"`
	Tags              []Tag       `json:"tags"`
	PullTime          string      `json:"pull_time"`
	Labels            []Label     `json:"labels"`
	References        []Reference `json:"references"`
	ManifestMediaType string      `json:"manifest_media_type"`
	ExtraAttrs        interface{} `json:"extra_attrs"`
	ID                int         `json:"id"`
	Digest            string      `json:"digest"`
	Icon              string      `json:"icon"`
	RepositoryID      int         `json:"repository_id"`
	AdditionLinks     interface{} `json:"addition_links"`
	MediaType         string      `json:"media_type"`
	ProjectID         int         `json:"project_id"`
	Type              string      `json:"type"`
	Annotations       interface{} `json:"annotations"`
}

type Tag struct {
	RepositoryID int    `json:"repository_id"`
	Name         string `json:"name"`
	PushTime     string `json:"push_time"`
	PullTime     string `json:"pull_time"`
	Signed       bool   `json:"signed"`
	ID           int    `json:"id"`
	Immutable    bool   `json:"immutable"`
	ArtifactID   int    `json:"artifact_id"`
}

type Label struct {
	UpdateTime   string `json:"update_time"`
	Description  string `json:"description"`
	Color        string `json:"color"`
	CreationTime string `json:"creation_time"`
	Deleted      bool   `json:"deleted"`
	Scope        string `json:"scope"`
	ProjectID    int    `json:"project_id"`
	ID           int    `json:"id"`
	Name         string `json:"name"`
}

type Reference struct {
	Platform    Platform    `json:"platform"`
	ChildDigest string      `json:"child_digest"`
	Urls        []string    `json:"urls"`
	ParentID    int         `json:"parent_id"`
	ChildID     int         `json:"child_id"`
	Annotations interface{} `json:"annotations"`
}

type Platform struct {
	Os           string   `json:"os"`
	Variant      string   `json:"variant"`
	Architecture string   `json:"architecture"`
	OsFeatures   []string `json:"'os.features'"`
	OsVersion    string   `json:"'os.version'"`
}
