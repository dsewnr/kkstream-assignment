package actions

import "github.com/gobuffalo/buffalo"

// AdminLogin default implementation.
func AdminLogin(c buffalo.Context) error {
	c.Session().Set("admin", true)
	return c.Render(200, r.HTML("admin/login.html"))
}

// AdminLogout default implementation.
func AdminLogout(c buffalo.Context) error {
	c.Session().Delete("admin")
	return c.Render(200, r.HTML("admin/logout.html"))
}
