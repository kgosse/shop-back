package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/kgosse/training/golang/src/projects/shop/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
