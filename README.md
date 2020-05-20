# go-tw-xauth

Twitter xAuth in Go.

## Installation

```bash
$ go get github.com/dqn/go-tw-xauth
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/dqn/go-tw-xauth/xauth"
)

func main() {
	resp, err := xauth.Do("CONSUMER_KEY", "CONSUMER_SECRET", "SCREEN_NAME", "PASSWORD")
	if err != nil {
		// Handle error
	}

	fmt.Println("oauth token:", resp.OauthToken)
	fmt.Println("oauth token secret:", resp.OauthTokenSecret)
	fmt.Println("screen name:", resp.ScreenName)
	fmt.Println("user id:", resp.UserID)
	fmt.Println("xauth expires:", resp.XAuthExpires)
}
```

## CLI

```bash
$ go get github.com/dqn/go-tw-xauth/cmd/xauth
```

```bash
$ xauth <consumer-key> <consumer-secret> <screen-name> <password>
oauth token: XXXXXXXXXXXXXXXXXX-XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
oauth token secret: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
screen name: XXXXX
user id: XXXXXXXXXXXXXXXXXX
xauth expires: 0
```
