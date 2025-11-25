#!/bin/bash

# Script to organize API documentation files into subfolders
# Each .md file will be moved to a subfolder named after the resource/datasource

cd /workspaces/terraform-provider-power-platform/provider_api_documentation || exit 1

for file in *.md; do
    # Skip if no .md files exist
    [ -e "$file" ] || continue
    
    # Extract the resource/datasource name (filename without .md extension)
    name="${file%.md}"
    
    # Create subfolder if it doesn't exist
    mkdir -p "$name"
    
    # Move the file into the subfolder
    mv "$file" "$name/"
    
    echo "Moved $file to $name/"
done

echo "Organization complete!"
