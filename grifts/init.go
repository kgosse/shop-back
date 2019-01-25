package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/kgosse/shop-back/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
