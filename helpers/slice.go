package helpers

func RemoveDuplicates[T comparable](slice []T) []T {
	/*
		Remove duplicates in unsorted slice
	*/
	encountered := make(map[T]bool)
	result := make([]T, 0, len(slice))

	for _, v := range slice {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}

	return result
}
