package xauth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const endpoint = "https://api.twitter.com/oauth/access_token"

type XAuthResponse struct {
	OauthToken       string
	OauthTokenSecret string
	UserID           string
	ScreenName       string
	XAuthExpires     string
}

type XAuthError struct {
	Err struct {
		Text string `xml:",chardata"`
		Code string `xml:"code,attr"`
	} `xml:"error"`
}

func (e *XAuthError) Error() string {
	return fmt.Sprintf("code %s: %s", e.Err.Code, e.Err.Text)
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

func buildQueryString(m map[string]string) string {
	pairs := make([]string, len(m))
	i := 0
	for k, v := range m {
		pairs[i] = k + "=" + url.QueryEscape(v)
		i++
	}
	s := strings.Join(pairs, "&")

	return s
}

func makeOAuthSignature(method, endpoint string, params map[string]string, consumerSecret string) (string, error) {
	ps := strings.Join(makePairs("%s=%s", params), "&")
	base := method + "&" + url.QueryEscape(endpoint) + "&" + url.QueryEscape(ps)
	key := url.QueryEscape(consumerSecret) + "&"

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

func Do(consumerKey, consumerSecret, screenName, password string) (*XAuthResponse, error) {
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

	sign, err := makeOAuthSignature("POST", endpoint, mergeMaps(oauth, data), consumerSecret)
	if err != nil {
		return nil, err
	}

	oauth["oauth_signature"] = sign
	token := makeAuthorization(oauth)

	reader := strings.NewReader(buildQueryString(data))
	req, err := http.NewRequest("POST", endpoint, reader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var xe XAuthError
	if xml.Unmarshal(b, &xe) == nil {
		return nil, &xe
	}

	q, err := url.ParseQuery(string(b))
	if len(q) != 5 {
		err = fmt.Errorf("unknown error: %s", b)
		return nil, err
	}

	xr := XAuthResponse{
		OauthToken:       q["oauth_token"][0],
		OauthTokenSecret: q["oauth_token_secret"][0],
		UserID:           q["user_id"][0],
		ScreenName:       q["screen_name"][0],
		XAuthExpires:     q["x_auth_expires"][0],
	}

	return &xr, nil
}
