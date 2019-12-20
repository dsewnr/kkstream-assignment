package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	csrf "github.com/gobuffalo/mw-csrf"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/google/uuid"
	"github.com/unrolled/secure"

	"assignment/models"

	"github.com/gobuffalo/buffalo-pop/pop/popmw"
	i18n "github.com/gobuffalo/mw-i18n"
	"github.com/gobuffalo/packr/v2"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_assignment_session",
		})

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)
		app.Use(generateSessionID) // Generate Session ID
		app.Use(checkIsLogin)      // Check is User Login

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))

		// Setup and use translations:
		app.Use(translations())

		app.GET("/", HomeHandler)

		// admin := app.Group("/admin")
		// admin.GET("/login", AdminLogin)
		// admin.GET("/logout", AdminLogout)

		// Login/Logout
		auth := app.Group("/auth")
		auth.GET("/login", AuthLogin)
		auth.GET("/logout", AuthLogout)
		auth.POST("/doLogin", AuthDoLogin)
		auth.GET("/callback/line", AuthCallbackLine)
		auth.Middleware.Skip(csrf.New, AuthDoLogin, AuthCallbackLine)
		auth.Middleware.Skip(checkIsLogin, AuthLogin, AuthLogout, AuthDoLogin, AuthCallbackLine)

		// Upload APIs
		api := app.Group("/api")
		api.Use(checkIsAdmin)
		api.POST("/upload", ApiUpload)
		api.POST("/uploadMulti", ApiUploadMulti)
		api.POST("/uploadGcs", ApiUploadGcs)
		api.Middleware.Skip(csrf.New, ApiUpload, ApiUploadMulti, ApiUploadGcs)
		api.Middleware.Skip(checkIsLogin, ApiUpload, ApiUploadMulti, ApiUploadGcs)
		api.Middleware.Skip(checkIsAdmin, ApiUpload, ApiUploadMulti)

		if ENV == "test" {
			api.Middleware.Skip(checkIsAdmin, ApiUploadGcs)
		}

		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(packr.New("app:locales", "../locales"), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}

func generateSessionID(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		sessionIDKey := envy.Get("SESSION_ID", "development")
		if c.Session().Get(sessionIDKey) == nil {
			c.Session().Set(sessionIDKey, uuid.New().String())
		}
		err := next(c)
		// do some work after calling the next handler
		return err
	}
}

// Check is User Login
func checkIsLogin(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		// do some work before calling the next handler

		if c.Session().Get("user") == nil {
			return c.Redirect(307, "authLoginPath()")
		}
		err := next(c)
		// do some work after calling the next handler
		return err
	}
}

// Check is Admin Role
func checkIsAdmin(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		// do some work before calling the next handler

		if c.Session().Get("is_admin") == nil {
			return c.Redirect(403, "rootPath()")
		}
		err := next(c)
		// do some work after calling the next handler
		return err
	}
}
