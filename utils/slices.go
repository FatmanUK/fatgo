package utils

// loop and copy nonempty - O(n)
// join and split - O(???)
func RemoveEmptyStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func CompareSlices[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !(a[i] == b[i]) {
			return false
		}
	}
	return true
}
