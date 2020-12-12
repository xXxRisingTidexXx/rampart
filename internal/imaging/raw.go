package imaging

import (
	"crypto/sha1"
	"image"
)

type Raw struct {
	Hash   [sha1.Size]byte
	Label  string
	Effect Effect
	Source image.Image
}
