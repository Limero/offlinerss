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
