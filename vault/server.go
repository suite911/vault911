package vclient

import (
	"encoding/json"
	"strconv"
	"unicode/utf8"

	"github.com/suite911/vault911/vault"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

func Recv(ctx *fasthttp.RequestCtx, key [32]byte) ([]byte, error) {
	if !ctx.IsPost() {
		return nil, errors.Wrap(errors.New("not POST"), "(*fasthttp.RequestCtx).IsPost")
	}
	args := ctx.PostArgs()
	tsBytes := args.Peek("ts")
	if !utf8.Valid(tsBytes) {
		return nil, errors.Wrap(errors.New("ts not valid Unicode"), "utf8.Valid")
	}
	tsString := string(tsBytes)
	ts, e := strconv.ParseUint(tsString, 10, 64)
	if e != nil {
		return nil, errors.Wrap(e, "strconv.ParseUint")
	}
	var v Vault
	v.TimeStamp = ts
	v.Payload = args.Peek("ct")
	pt, ok := v.Decode(key)
	if !ok {
		return nil, errors.Wrap(errors.New("unable to decode message"), "(*vault.Vault).Decode")
	}
	return pt, nil
}

func Reply(ctx *fasthttp.RequestCtx, message []byte, key [32]byte) (http500 string, err error) {
	var v Vault
	v.Init(message, key)
	b, e := json.Marshal(v)
	if e != nil {
		return "Internal Server Error: Unable to marshal JSON", errors.Wrap(e, "json.Marshal")
	}
	_, e := ctx.Write(b)
	if e != nil {
		return "Internal Server Error: Unable to write reply", errors.Wrap(e, "(*fasthttp.RequestCtx).Write")
	}
	return "", nil
}
