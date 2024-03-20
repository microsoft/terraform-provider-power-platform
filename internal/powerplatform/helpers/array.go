// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform_helpers

func DiffArrays(newArr, oldArr []string) ([]string, []string) {
	added := make([]string, 0)
	removed := make([]string, 0)

	oldMap := make(map[string]bool)
	for _, item := range oldArr {
		oldMap[item] = true
	}

	newMap := make(map[string]bool)
	for _, item := range newArr {
		newMap[item] = true
	}

	for _, item := range newArr {
		if _, found := oldMap[item]; !found {
			added = append(added, item)
		}
	}

	for _, item := range oldArr {
		if _, found := newMap[item]; !found {
			removed = append(removed, item)
		}
	}

	return added, removed
}
