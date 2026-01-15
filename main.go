package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	orgName := os.Getenv("ORGANISATION_NAME")
	repoName := os.Getenv("REPOSITORY_NAME")
	bearerToken := os.Getenv("BEARER_TOKEN")

	if orgName == "" || repoName == "" || bearerToken == "" {
		log.Fatal("Missing required environment variables: ORGANISATION_NAME, REPOSITORY_NAME, BEARER_TOKEN")
	}

	fmt.Printf("Fetching artifacts for %s/%s...\n", orgName, repoName)

	artifacts, err := FetchAllArtifacts(orgName, repoName, bearerToken)
	if err != nil {
		log.Fatalf("Failed to fetch artifacts: %v", err)
	}

	fmt.Printf("\nTotal artifacts found: %d\n", artifacts.TotalCount)
	fmt.Printf("Retrieved %d artifacts from API\n", len(artifacts.Artifacts))
	for i, artifact := range artifacts.Artifacts {
		fmt.Printf("\nArtifact %d:\n", i+1)
		fmt.Printf("  Name: %s\n", artifact.Name)
		fmt.Printf("  ID: %d\n", artifact.ID)
		fmt.Printf("  Size: %d bytes\n", artifact.SizeInBytes)
		fmt.Printf("  Created: %s\n", artifact.CreatedAt)
		fmt.Printf("  Expires: %s\n", artifact.ExpiresAt)
		fmt.Printf("  Expired: %v\n", artifact.Expired)
	}
}
