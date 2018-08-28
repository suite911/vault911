package vault

import (
	"encoding/binary"
	"time"

	"golang.org/x/crypto/salsa20"
	"golang.org/x/crypto/sha3"
)

type Vault struct {
	TimeStamp uint64 `json:"ts"` // Time stamp
	Payload   []byte `json:"ct"` // Encrypted message ciphertext
}

func New(plaintext []byte, key Key) *Vault {
	return new(Vault).Init(plaintext, key)
}

func (v *Vault) Init(plaintext []byte, key Key) *Vault {
	its := time.Now().UTC().UnixNano()
	if its < 0 {
		panic("bad system time")
	}
	ts := uint64(its)
	buf := make([]byte, 8, 8+len(key))
	binary.LittleEndian.PutUint64(buf, ts)
	buf = append(buf, key...)
	dig := sha3.Sum256(buf)
	auth := sha3.Sum256(plaintext)
	plaintext = append(plaintext, auth...)
	ct := make([]byte, len(plaintext))
	salsa20.XORKeyStream(ct, plaintext, buf[:8], &dig)
	v.TimeStamp = ts
	v.Payload = ct
	return v
}

func (v *Vault) Decrypt(key Key) (plaintext []byte, ok bool) {
	ts, ct := v.TimeStamp, v.Payload
	if len(ct) < 32 {
		return nil, false
	}
	buf := make([]byte, 8, 8+len(key))
	binary.LittleEndian.PutUint64(buf, ts)
	buf = append(buf, key...)
	dig := sha3.Sum256(buf)
	pt := make([]byte, len(ct))
	salsa20.XORKeyStream(pt, ct, buf[:8], &dig)
	end := len(pt) - 32
	auth := pt[end:]
	pt = pt[:end]
	actual := sha3.Sum256(pt)
	if actual[:] != auth {
		return nil, false
	}
	return pt, true
}
