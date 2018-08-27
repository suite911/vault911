package main

import (
	"flag"
	"strings"
	"unicode/utf8"

	"github.com/suite911/vault911/vclient"
	"github.com/suite911/vault911/vkey"
	"github.com/valyala/fasthttp"
)

func main() {
	pKeyFile := flag.String("f", "", "Key file")
	pPW := flag.String("pw", "password", "Password to use")
	pSalt := flag.String("salt", "salt", "Per-user salt to use")
	pURL := flag.String("s", "http://localhost:8080", "Server address")
	flag.Parse()
	key, err := vkey.New(*pPW, *pSalt, *pKeyFile)
	if err != nil {
		panic(err)
	}
	var args fasthttp.Args
	_, plaintext, err := vclient.Post(*pURL, strings.Join(flag.Args(), " "), &args, key)
	if err != nil {
		panic(err)
	}
	if !utf8.Valid(plaintext) {
		panic("Invalid UTF-8")
	}
	if _, err := io.Write(plaintext); err != nil {
		panic(err)
	}
}
