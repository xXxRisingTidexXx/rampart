package imaging

import (
	"crypto/sha1"
)

type Raw struct {
	Hash   [sha1.Size]byte
	Label  string
	Effect *Effect
	Bytes  []byte
}
