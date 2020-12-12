package imaging

import (
	"crypto/sha1"
)

type Asset struct {
	Hash   [sha1.Size]byte
	Label  string
	Effect string
	Bytes  []byte
}
