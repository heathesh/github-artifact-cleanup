package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	orgName := os.Getenv("ORGANISATION_NAME")
	repoName := os.Getenv("REPOSITORY_NAME")
	bearerToken := os.Getenv("BEARER_TOKEN")
	numberOfDaysStr := os.Getenv("NUMBER_OF_DAYS")

	if orgName == "" || repoName == "" || bearerToken == "" || numberOfDaysStr == "" {
		log.Fatal("Missing required environment variables: ORGANISATION_NAME, REPOSITORY_NAME, BEARER_TOKEN, NUMBER_OF_DAYS")
	}

	numberOfDays, err := strconv.Atoi(numberOfDaysStr)
	if err != nil {
		log.Fatalf("Invalid NUMBER_OF_DAYS value: %v", err)
	}

	fmt.Printf("Fetching artifacts for %s/%s...\n", orgName, repoName)

	artifacts, err := FetchAllArtifacts(orgName, repoName, bearerToken)
	if err != nil {
		log.Fatalf("Failed to fetch artifacts: %v", err)
	}

	fmt.Printf("\nTotal artifacts found: %d\n", artifacts.TotalCount)
	fmt.Printf("Retrieved %d artifacts from API\n", len(artifacts.Artifacts))

	// Calculate cutoff date
	cutoffDate := time.Now().AddDate(0, 0, -numberOfDays)
	fmt.Printf("Deleting artifacts older than %d days (before %s)\n\n", numberOfDays, cutoffDate.Format("2006-01-02"))

	// Filter and delete old artifacts
	var deletedCount, skippedCount int
	for _, artifact := range artifacts.Artifacts {
		if artifact.CreatedAt.Before(cutoffDate) {
			fmt.Printf("Deleting artifact: %s (ID: %d, Created: %s)\n", artifact.Name, artifact.ID, artifact.CreatedAt.Format("2006-01-02"))
			if err := DeleteArtifact(orgName, repoName, bearerToken, artifact.ID); err != nil {
				log.Printf("Failed to delete artifact %d: %v\n", artifact.ID, err)
			} else {
				deletedCount++
			}
		} else {
			skippedCount++
		}
	}

	fmt.Printf("\nCleanup complete!\n")
	fmt.Printf("Artifacts deleted: %d\n", deletedCount)
	fmt.Printf("Artifacts kept: %d\n", skippedCount)
}
