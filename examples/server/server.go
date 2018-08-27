package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"unicode/utf8"

	"github.com/suite911/vault911/vault"
	"github.com/valyala/fasthttp"
)

var key vault.Key

func main() {
	pPort := flag.Int("p", 8080, "Port on which to listen")
	pKeyFile := flag.String("f", "", "Key file")
	pPW := flag.String("pw", "password", "Password to use")
	pSalt := flag.String("salt", "salt", "Per-user salt to use")
	flag.Parse()
	if err := key.Init(*pPW, *pSalt, *pKeyFile); err != nil {
		panic(err)
	}
	if err := fasthttp.ListenAndServe(":"+strconv.Itoa(*pPort), handler); err != nil {
		panic(err)
	}
}

func handler(ctx *fasthttp.RequestCtx) {
	b, err := vault.Recv(ctx, key)
	if err != nil || !utf8.Valid(b) {
		ctx.Error("Bad Request", 400)
		return
	}
	message := string(b)
	reply := fmt.Sprintf("You said %q.", message)
	http500, err := vault.Reply(ctx, []byte(reply), key)
	if err != nil {
		log.Println(err)
		ctx.Error(http500, 500)
		return
	}
	ctx.SetStatusCode(200)
}
