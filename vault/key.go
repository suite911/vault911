package vault

import (
	"io/ioutil"
	"unicode/utf8"

	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/sha3"
)

type Key [32]byte

func NewKey(password, salt []byte, keyFile string) (*Key, error) {
	k := new(Key)
	e := k.Init(password, salt, keyFile)
	return k, e
}

func (k *Key) Init(password, salt []byte, keyFile string) error {
	if !utf8.Valid(password) {
		return errors.Wrap(errors.New("password is not valid Unicode"), "utf8.Valid")
	}
	a2id := argon2.IDKey(password, salt, 1, 64*1024, 4, uint32(len(*k)))
	if len(a2id) != len(*k) {
		return errors.Wrap(errors.New("wrong length of IDKey"), "argon2.IDKey")
	}
	copy((*k)[:], a2id)
	if len(keyFile) > 0 {
		b, e := ioutil.ReadFile(keyFile)
		if e != nil {
			return errors.Wrap(e, "ioutil.ReadFile")
		}
		dig := sha3.Sum256(b)
		j := len(*k)
		if j > len(dig) {
			j = len(dig)
		}
		for i := 0; i < j; i++ {
			(*k)[i] ^= dig[i]
		}
	}
	return nil
}
