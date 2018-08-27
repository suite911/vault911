package main

import (
	"flag"
	"strings"
	"unicode/utf8"

	"github.com/suite911/vault911/vclient"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/sha3"
)

func main() {
	pPW := flag.String("p", "password", "Password to use")
	pSalt := flag.String("s", "salt", "Per-user salt to use")
	pURL := flag.String("u", "http://localhost:8080", "Server address")
	flag.Parse()
	key := sha3.Sum256(argon2.IDKey(*pPW, *pSalt, 1, 64*1024, 4, 32))
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
