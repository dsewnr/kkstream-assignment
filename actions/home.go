package actions

import "github.com/gobuffalo/buffalo"

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	isAdmin := c.Session().Get("is_admin") == true
	routePath := "apiUploadMultiPath"
	if isAdmin {
		routePath = "apiUploadGcsPath"
	}
	uploadURI := ""
	ri, err := app.Routes().Lookup(routePath)
	if err != nil {
		c.Logger().Debug("No such route")
	} else {
		uploadURI = ri.Path
	}
	c.Set("isAdmin", isAdmin)
	c.Set("uploadURI", uploadURI)
	c.Set("userID", c.Session().Get("user"))
	return c.Render(200, r.HTML("_index.html"))
}
