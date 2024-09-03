package utils

import "strings"

// similar to https://github.com/zspecza/common-tags
func OneLineTrim(in string) string {
	ll := strings.Split(in, "\n")
	for i, l := range ll {
		ll[i] = strings.Trim(l, "\t ")
	}
	return strings.Join(ll, " ")
}
