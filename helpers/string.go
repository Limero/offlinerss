package helpers

func CondString(condition bool, primary string, fallback string) string {
	/*
		If condition is true, return primary string, else return fallback string
	*/
	if condition {
		return primary
	}
	return fallback
}

func RemoveDuplicates(slice []string) []string {
	/*
		Remove duplicates in string slice
	*/
	encountered := make(map[string]bool)
	result := []string{}

	for _, v := range slice {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}

	return result
}
