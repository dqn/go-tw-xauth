package xauth

import (
	"os"
	"testing"
)

func TestXAuth(t *testing.T) {
	resp, err := XAuth(os.Getenv("CK"), os.Getenv("CS"), os.Getenv("SN"), os.Getenv("PW"))
	if err != nil {
		t.Fatal(err)
	}
	println(resp)
}
