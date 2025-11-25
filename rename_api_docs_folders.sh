#!/bin/bash

# Script to rename subfolders with "resource_" or "datasource_" prefix
# based on the content of the .md file inside

cd /workspaces/terraform-provider-power-platform/provider_api_documentation || exit 1

for dir in */; do
    # Remove trailing slash
    dir="${dir%/}"
    
    # Find the .md file in the directory
    mdfile="$dir/${dir}.md"
    
    if [ ! -f "$mdfile" ]; then
        echo "Warning: $mdfile not found, skipping $dir"
        continue
    fi
    
    # Check if it's a datasource or resource by looking at the content
    # Datasources typically have "Data Source:" or "data source" in the title/header
    # Resources typically have "Resource:" or don't have "data source"
    
    if grep -qi "data source" "$mdfile" | head -20; then
        prefix="datasource"
    else
        prefix="resource"
    fi
    
    # Create new directory name
    new_name="${prefix}_${dir}"
    
    # Rename if different
    if [ "$dir" != "$new_name" ]; then
        mv "$dir" "$new_name"
        echo "Renamed $dir to $new_name"
    else
        echo "Skipped $dir (already has correct prefix)"
    fi
done

echo "Renaming complete!"
