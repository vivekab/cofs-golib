package golibarray

func SafeIndex[T any](arr []T, idx int, defaultValue T) T {
	if idx < 0 || idx >= len(arr) {
		return defaultValue
	}
	return arr[idx]
}

func Contains[T comparable](arr []T, searchValue T) bool {
	for _, v := range arr {
		if v == searchValue {
			return true
		}
	}
	return false
}
