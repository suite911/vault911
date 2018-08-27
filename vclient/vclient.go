package vclient

import (
	"encoding/json"
	"strconv"

	"github.com/suite911/vault911/vault"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

func Post(url string, message []byte, args *fasthttp.Args, key [32]byte) (statusCode int, plaintext []bytes, err error) {
	var v vault.Vault
	v.Init(message, key)
	args.Set("ts", strconv.Itoa(v.TimeStamp))
	args.SetBytesV("ct", v.Payload)
	sc, body, e := fasthttp.Post(nil, url, args)
	if e == nil && (sc < 200 || sc > 299) {
		e = errors.New("HTTP "+strconv.Itoa(sc))
	}
	if e != nil {
		return sc, nil, errors.Wrap(e, "fasthttp.Post")
	}
	var reply Vault
	if e := json.Unmarshal(body, &reply); e != nil {
		return sc, nil, errors.Wrap(e, "json.Unmarshal")
	}
	pt, ok := reply.Decode(key)
	if !ok {
		return sc, nil, errors.Wrap(errors.New("unable to decode reply"), "(*vault.Vault).Decode")
	}
	return sc, pt, nil
}
