# GitHub Artifact Cleanup

A Go application that automatically cleans up old GitHub Actions artifacts from a repository based on a configurable retention policy.

## Overview

This tool helps manage GitHub Actions artifacts by:
- Fetching all artifacts from a specified repository
- Identifying artifacts older than a configured number of days
- Automatically deleting outdated artifacts to free up storage space

## Features

- **Paginated API calls**: Handles repositories with hundreds or thousands of artifacts
- **Configurable retention period**: Set how many days to keep artifacts
- **Detailed logging**: Shows progress and summary of deleted vs kept artifacts
- **Environment-based configuration**: Secure credential management via `.env` file

## Prerequisites

- Go 1.16 or higher
- GitHub Personal Access Token with the following permissions:
  - `repo` (Full control of private repositories)
  - `actions` (Read and write actions artifacts)

## Setup

1. Clone this repository

2. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

3. Edit `.env` and fill in your values:
   ```
   ORGANISATION_NAME=your-org-name
   REPOSITORY_NAME=your-repo-name
   BEARER_TOKEN=ghp_your_github_token
   NUMBER_OF_DAYS=30
   ```

4. Install dependencies:
   ```bash
   go mod download
   ```

## Configuration

| Variable | Description | Example |
|----------|-------------|---------|
| `ORGANISATION_NAME` | GitHub organization or username | `mycompany` |
| `REPOSITORY_NAME` | Repository name | `my-project` |
| `BEARER_TOKEN` | GitHub Personal Access Token | `ghp_xxxx...` |
| `NUMBER_OF_DAYS` | Retention period in days | `30` |

## Usage

### Build the application:
```bash
go build
```

### Run the application:
```bash
./github-artifact-cleanup
```

Or run directly without building:
```bash
go run .
```

## Example Output

```
Fetching artifacts for mycompany/my-project...
Fetching page 2 of 5...
Fetching page 3 of 5...
Fetching page 4 of 5...
Fetching page 5 of 5...

Total artifacts found: 150
Retrieved 150 artifacts from API
Deleting artifacts older than 30 days (before 2025-12-16)

Deleting artifact: integration-test-results (ID: 5098109034, Created: 2025-12-10)
Deleting artifact: build-artifacts (ID: 5087423112, Created: 2025-12-08)
Deleting artifact: coverage-report (ID: 5076234891, Created: 2025-12-05)

Cleanup complete!
Artifacts deleted: 45
Artifacts kept: 105
```

## How It Works

1. **Load Configuration**: Reads environment variables from `.env` file
2. **Fetch Artifacts**: Makes paginated API calls to retrieve all artifacts (30 per page)
3. **Calculate Cutoff Date**: Determines which artifacts are older than the retention period
4. **Delete Old Artifacts**: Iterates through artifacts and deletes those older than the cutoff
5. **Report Results**: Displays summary of deleted and kept artifacts

## API Endpoints Used

- `GET /repos/{org}/{repo}/actions/artifacts?page={page}&per_page=30` - List artifacts
- `DELETE /repos/{org}/{repo}/actions/artifacts/{artifact_id}` - Delete artifact

## Security Notes

- Never commit your `.env` file (it's already in `.gitignore`)
- Store your GitHub token securely
- Use tokens with minimal required permissions
- Rotate tokens regularly

## License

See LICENSE file for details.
