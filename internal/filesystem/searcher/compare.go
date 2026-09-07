package searcher

import "strings"

func compareNames(left string, right string) int {
	return strings.Compare(strings.ToLower(left), strings.ToLower(right))
}
