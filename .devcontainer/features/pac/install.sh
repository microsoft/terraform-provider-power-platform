#!/bin/sh
set -e

echo "Activating feature 'Power Platform CLI'"

# Install .NET Core SDK
sudo apt-get update && \
  sudo apt-get install -y dotnet-sdk-8.0

# Install Power Platform CLI to vscode user's path
dotnet tool install Microsoft.PowerApps.CLI.Tool --interactive false --tool-path /home/vscode/.dotnet/tools
