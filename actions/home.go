package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/x/responder"
)

func HomeHandler(c buffalo.Context) error {
	return responder.Wants("html", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.HTML("home/index.plush.html"))
	}).Respond(c)
}
