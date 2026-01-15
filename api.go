package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchArtifacts(orgName, repoName, bearerToken string, page int) (*ArtifactsResponse, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/artifacts?page=%d&per_page=30", orgName, repoName, page)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var artifactsResponse ArtifactsResponse
	if err := json.Unmarshal(body, &artifactsResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &artifactsResponse, nil
}

func FetchAllArtifacts(orgName, repoName, bearerToken string) (*ArtifactsResponse, error) {
	// Fetch first page to get total count
	firstPage, err := FetchArtifacts(orgName, repoName, bearerToken, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch first page: %w", err)
	}

	// If there's only one page or no artifacts, return immediately
	if firstPage.TotalCount <= 30 {
		return firstPage, nil
	}

	// Calculate total number of pages
	totalPages := (firstPage.TotalCount + 29) / 30 // Ceiling division

	// Initialize result with first page
	allArtifacts := &ArtifactsResponse{
		TotalCount: firstPage.TotalCount,
		Artifacts:  make([]Artifact, 0, firstPage.TotalCount),
	}
	allArtifacts.Artifacts = append(allArtifacts.Artifacts, firstPage.Artifacts...)

	// Fetch remaining pages
	for page := 2; page <= totalPages; page++ {
		fmt.Printf("Fetching page %d of %d...\n", page, totalPages)
		pageResponse, err := FetchArtifacts(orgName, repoName, bearerToken, page)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch page %d: %w", page, err)
		}
		allArtifacts.Artifacts = append(allArtifacts.Artifacts, pageResponse.Artifacts...)
	}

	return allArtifacts, nil
}

func DeleteArtifact(orgName, repoName, bearerToken string, artifactID int64) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/artifacts/%d", orgName, repoName, artifactID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
