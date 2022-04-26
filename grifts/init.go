package grifts

import (
	"github.com/magikid/goshorter/actions"

	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
