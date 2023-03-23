package krew

import (
	"fmt"
	"os"
)

const (
	krewIndexRepoName  = "artificial-krew-repository"
	krewIndexRepoOwner = "artificialinc"
)

// GetKrewIndexRepoName returns the krew-index repo name
func GetKrewIndexRepoName() string {
	override := os.Getenv("UPSTREAM_KREW_INDEX_REPO_NAME")
	if override != "" {
		fmt.Println("overriding krew index repo name")
		return override
	}

	return krewIndexRepoName
}

// GetKrewIndexRepoOwner returns the krew-index repo owner
func GetKrewIndexRepoOwner() string {
	override := os.Getenv("UPSTREAM_KREW_INDEX_REPO_OWNER")
	if override != "" {
		fmt.Println("overriding krew index repo owner")
		return override
	}

	return krewIndexRepoOwner
}
