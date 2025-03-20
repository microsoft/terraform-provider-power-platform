#!/bin/bash

# Check if gh CLI is installed
if ! command -v gh &> /dev/null; then
    echo "Error: GitHub CLI (gh) is not installed. Please install it first."
    exit 1
fi

# Get the current branch name
current_branch=$(git branch --show-current)
echo "Current branch: $current_branch"

# Try to extract issue number from branch name
# Common patterns: issue-123, 123-description, feature/123-description
issue_number=""

# Try pattern like "issue-123"
if [[ "$current_branch" =~ issues?[-/]([0-9]+) ]]; then
    issue_number="${BASH_REMATCH[1]}"
# Try pattern with numbers at start or after slash/dash
elif [[ "$current_branch" =~ (^|/|-)([0-9]+)[-/] ]]; then
    issue_number="${BASH_REMATCH[2]}"
fi

# If no issue found in branch name, try to get the PR that this branch is for
if [[ -z "$issue_number" ]]; then
    echo "Trying to find associated PR for this branch..."
    pr_data=$(gh pr list --head "$current_branch" --json number --limit 1 2>/dev/null)
    
    if [[ -n "$pr_data" && "$pr_data" != "[]" ]]; then
        pr_number=$(echo "$pr_data" | grep -o '"number":[0-9]*' | grep -o '[0-9]*')
        if [[ -n "$pr_number" ]]; then
            echo "Found PR #$pr_number, checking for linked issues..."
            issue_data=$(gh pr view "$pr_number" --json closingIssues --jq '.closingIssues[0].number' 2>/dev/null)
            if [[ -n "$issue_data" && "$issue_data" != "null" ]]; then
                issue_number=$issue_data
            fi
        fi
    fi
fi

# If still no issue found, ask user
if [[ -z "$issue_number" ]]; then
    echo "Could not identify an issue number from branch name: $current_branch"
    echo "Please enter issue number:"
    read -r issue_number
fi

# Validate that we have an issue number
if [[ -z "$issue_number" ]]; then
    echo "No issue number provided. Exiting."
    exit 1
fi

echo "Using issue number: $issue_number"

# Create the directory if it doesn't exist
mkdir -p "/workspaces/terraform-provider-power-platform/.github/prompts"

# Use gh CLI to get the issue content and save it to file
output_file="/workspaces/terraform-provider-power-platform/.github/prompts/.userstory.prompt.md"
if gh issue view "$issue_number" > "$output_file" 2>/dev/null; then
    echo "Successfully saved issue #$issue_number to $output_file"
    echo "Path: $output_file"
else
    echo "Error: Could not fetch issue #$issue_number. Make sure it exists and you have permissions."
    exit 1
fi
