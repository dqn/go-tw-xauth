package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dqn/go-tw-xauth/xauth"
)

func main() {
	flag.Parse()
	if flag.NArg() != 4 {
		flag.Usage()
		os.Exit(1)
	}

	resp, err := xauth.Do(flag.Arg(0), flag.Arg(1), flag.Arg(2), flag.Arg(3))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("oauth token:", resp.OauthToken)
	fmt.Println("oauth token secret:", resp.OauthTokenSecret)
	fmt.Println("screen name:", resp.ScreenName)
	fmt.Println("user id:", resp.UserID)
	fmt.Println("xauth expires:", resp.XAuthExpires)

	os.Exit(0)
}
