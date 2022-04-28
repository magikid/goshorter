package actions

import (
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/magikid/goshorter/models"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"

	"github.com/gobuffalo/pop/v6"
	"github.com/pkg/errors"
)

func init() {
	gothic.Store = App().SessionStore

	goth.UseProviders(
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://localhost:8080/auth/github/callback"),
	)
}

func AuthCallback(c buffalo.Context) error {
	gothUser, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(401, err)
	}

	if gothUser.UserID != "115515" {
		return c.Error(401, errors.New("You are not authorized to use this service"))
	}

	q := models.DB.Where("provider = ? AND provider_id = ?", gothUser.Provider, gothUser.UserID)
	found, err := q.Exists("users")
	if err != nil {
		return errors.WithStack(err)
	}

	user := &models.User{}
	if found {
		err = q.First(user)
		if err != nil {
			return errors.WithStack(err)
		}
	} else {
		user.Name = gothUser.Name
		user.Email = gothUser.Email
		user.Provider = gothUser.Provider
		user.ProviderID = gothUser.UserID

		validationIssues, err := models.DB.ValidateAndCreate(user)
		if err != nil {
			return errors.WithStack(err)
		}

		if validationIssues.HasAny() {
			c.Flash().Add("danger", validationIssues.String())
			return c.Redirect(302, "/")
		}
	}

	c.Session().Set("current_user_id", user.ID)
	err = c.Session().Save()
	if err != nil {
		return errors.WithStack(err)
	}

	c.Flash().Add("success", "You have been logged in")
	return c.Redirect(302, "shortenedLinkPath()")
}

func AuthDestroy(c buffalo.Context) error {
	c.Session().Clear()
	c.Flash().Add("success", "You have been logged out")
	return c.Redirect(302, "/")
}

func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			tx := c.Value("tx").(*pop.Connection)
			if err := tx.Find(u, uid); err != nil {
				return errors.WithStack(err)
			}
			c.Set("current_user", u)
		}
		return next(c)
	}
}

func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid == nil {
			c.Flash().Add("danger", "You must be authorized to see that page")
			return c.Redirect(302, "/")
		}
		return next(c)
	}
}
