package helpers

func Cond[T any](condition bool, primary, fallback T) T {
	/*
		If condition is true, return primary, else return fallback
		https://github.com/golang/go/issues/66062
	*/
	if condition {
		return primary
	}
	return fallback
}
