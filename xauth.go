package xauth

import (
	"fmt"
	"sort"
	"strings"
)

type XAuthResponse struct {
}

func getKeys(m map[string]string) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	return keys
}

func makePairs(m map[string]string) []string {
	keys := getKeys(m)
	sort.Strings(keys)

	pairs := make([]string, len(keys))
	for i, k := range keys {
		pairs[i] = fmt.Sprintf(`%s="%s"`, k, m[k])
	}

	return pairs
}

func makeAuthorization(params map[string]string) string {
	pairs := makePairs(params)
	return "OAuth " + strings.Join(pairs, ", ")
}

func XAuth(consumerKey, consumerSecret, screenName, password string) (*XAuthResponse, error) {

	return nil, nil
}
