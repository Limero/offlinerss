package helpers

func Cond[T any](condition bool, primary, fallback T) T {
	/*
		If condition is true, return primary, else return fallback
	*/
	if condition {
		return primary
	}
	return fallback
}
