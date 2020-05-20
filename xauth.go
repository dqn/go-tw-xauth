package xauth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type XAuthResponse struct {
}

const endpoint = "https://api.twitter.com/oauth/access_token"

func getKeys(m map[string]string) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	return keys
}

func makePairs(format string, m map[string]string) []string {
	keys := getKeys(m)
	sort.Strings(keys)

	pairs := make([]string, len(keys))
	for i, k := range keys {
		pairs[i] = fmt.Sprintf(format, k, m[k])
	}

	return pairs
}

func mergeMaps(a, b map[string]string) map[string]string {
	m := make(map[string]string, len(a)+len(b))

	for k, v := range a {
		m[k] = v
	}
	for k, v := range b {
		m[k] = v
	}

	return m
}

func makeOAuthSignature(method, endpoint string, params map[string]string) (string, error) {
	ps := strings.Join(makePairs(`%s=%s`, params), "&")
	base := method + "&" + url.QueryEscape(endpoint) + "&" + url.QueryEscape(ps)
	key := url.QueryEscape(params["oauth_consumer_key"]) + "&"

	h := hmac.New(sha1.New, []byte(key))
	_, err := h.Write([]byte(base))
	if err != nil {
		return "", err
	}

	s := url.QueryEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))

	return s, nil
}

func makeAuthorization(params map[string]string) string {
	pairs := makePairs(`%s="%s"`, params)
	return "OAuth " + strings.Join(pairs, ", ")
}

func XAuth(consumerKey, consumerSecret, screenName, password string) (*XAuthResponse, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	oauth := map[string]string{
		"oauth_consumer_key":     consumerKey,
		"oauth_nonce":            timestamp,
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        timestamp,
		"oauth_version":          "1.0a",
	}

	data := map[string]string{
		"x_auth_mode":     "client_auth",
		"x_auth_password": password,
		"x_auth_username": screenName,
	}

	merged := mergeMaps(oauth, data)

	sign, err := makeOAuthSignature("POST", endpoint, merged)
	if err != nil {
		return nil, err
	}

	println(sign)

	return nil, nil
}
