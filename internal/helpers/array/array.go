// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package array

// DiffArrays returns the added and removed items between two string arrays.
// This can be useful for comparing plan vs state arrays.
func Diff(newArr, oldArr []string) (added []string, removed []string) {
	addedElements := make([]string, 0)
	removedElements := make([]string, 0)
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
			addedElements = append(addedElements, item)
		}
	}
	for _, item := range oldArr {
		if _, found := newMap[item]; !found {
			removedElements = append(removedElements, item)
		}
	}

	return addedElements, removedElements
}

// ArrayContains returns true if the given array contains the given item.
func Contains[T comparable](arr []T, item T) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

// Except returns a slice of elements that are in 'a' but not in 'b'.
func Except[T comparable](a, b []T) []T {
	bSet := make(map[T]struct{}, len(b))
	for _, value := range b {
		bSet[value] = struct{}{}
	}

	var diff []T
	for _, value := range a {
		if _, found := bSet[value]; !found {
			diff = append(diff, value)
		}
	}

	return diff
}

// Find returns the first element in the array that satisfies the predicate.
func Find[T comparable](arr []T, predicate func(T) bool) T {
	for _, v := range arr {
		if predicate(v) {
			return v
		}
	}
	var result T
	return result
}
