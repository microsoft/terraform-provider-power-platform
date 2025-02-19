package helpers

// Contains checks if a slice contains a specific element.
func Contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}
