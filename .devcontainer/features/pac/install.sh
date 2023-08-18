#!/bin/sh
set -e

echo "Activating feature 'Power Platform CLI'"

/usr/local/dotnet/current/dotnet tool install --global Microsoft.PowerApps.CLI.Tool --interactive false