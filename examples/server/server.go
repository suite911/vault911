package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"unicode/utf8"

	"github.com/suite911/vault911/vserver"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/sha3"
)

var key [32]byte

func main() {
	pPort := flag.Int("p", 8080, "Port on which to listen")
	pPW := flag.String("pw", "password", "Password to use")
	pSalt := flag.String("salt", "salt", "Per-user salt to use")
	flag.Parse()
	key = sha3.Sum256(argon2.IDKey(*pPW, *pSalt, 1, 64*1024, 4, 32))
	if err := fasthttp.ListenAndServe(":"+strconv.Itoa(*pPort), handler); err != nil {
		panic(err)
	}
}

func handler(ctx *fasthttp.RequestCtx) {
	b, err := vserver.Recv(ctx, key)
	if err != nil || !utf8.Valid(b) {
		ctx.Error("Bad Request", 400)
		return
	}
	message := string(b)
	reply := fmt.Sprintf("You said %q.", message)
	http500, err := vserver.Reply(ctx, []byte(reply), key)
	if err != nil {
		log.Println(err)
		ctx.Error(http500, 500)
		return
	}
	ctx.SetStatusCode(200)
}
