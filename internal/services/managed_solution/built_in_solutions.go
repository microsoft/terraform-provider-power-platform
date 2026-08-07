// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package managedsolution

import (
	_ "embed"
	"encoding/json"
	"strings"
)

//go:embed built-in-solutions.json
var builtInSolutionsJSON []byte

type builtInSolutionRegistryPayload struct {
	SolutionUniqueNames []string `json:"solution_unique_names"`
}

var builtInSolutionNames = loadBuiltInSolutionNames()

func loadBuiltInSolutionNames() map[string]struct{} {
	var payload builtInSolutionRegistryPayload
	if err := json.Unmarshal(builtInSolutionsJSON, &payload); err != nil {
		panic(err)
	}

	names := make(map[string]struct{}, len(payload.SolutionUniqueNames))
	for _, name := range payload.SolutionUniqueNames {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" {
			continue
		}
		names[strings.ToLower(trimmed)] = struct{}{}
	}

	return names
}

func isBuiltInSolutionDependency(name string) bool {
	_, exists := builtInSolutionNames[strings.ToLower(strings.TrimSpace(name))]
	return exists
}
