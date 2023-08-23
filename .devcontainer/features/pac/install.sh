#!/bin/sh
set -e

echo "Activating feature 'Power Platform CLI'"

# Install .NET Core SDK
sudo apt-get update && \
  sudo apt-get install -y dotnet-sdk-7.0

# Install Power Platform CLI
dotnet tool install --global Microsoft.PowerApps.CLI.Tool --interactive false 
# BUGBUG: dotnet tool install installs to /root/.dotnet/tools when running as root which is not accessible by the remote user in the dev container.
