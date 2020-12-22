package imaging

import (
	"image"
)

type Raw struct {
	Hash   string
	Group  string
	Label  string
	Effect Effect
	Source image.Image
}
