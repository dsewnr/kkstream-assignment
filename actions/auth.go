package actions

import (
	"assignment/line"
	"assignment/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// Line Token Response
type LineTokenResp struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

var LINE_AUTH_URL = ""
var LINE_TOKEN_URL = ""
var LINE_PROFILE_URL = ""
var LINE_CHANNEL_ID = ""
var LINE_CHANNEL_SECRET = ""
var LINE_CALLBACK = ""

func init() {
	LINE_AUTH_URL = envy.Get("LINE_AUTH_URL", "")
	LINE_TOKEN_URL = envy.Get("LINE_TOKEN_URL", "")
	LINE_PROFILE_URL = envy.Get("LINE_PROFILE_URL", "")
	LINE_CHANNEL_ID = envy.Get("LINE_CHANNEL_ID", "")
	LINE_CHANNEL_SECRET = envy.Get("LINE_CHANNEL_SECRET", "")
	LINE_CALLBACK = envy.Get("LINE_CALLBACK", "")
}

// AuthLogin default implementation.
func AuthLogin(c buffalo.Context) error {

	u, err := url.Parse(LINE_AUTH_URL)
	if err != nil {
		return c.Error(500, err)
	}
	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", LINE_CHANNEL_ID)
	q.Set("redirect_uri", LINE_CALLBACK)
	q.Set("state", c.Session().Get(envy.Get("SESSION_ID", "development")).(string))
	q.Set("scope", "openid profile")

	u.RawQuery = q.Encode()

	c.Set("LineLoginURL", u.String())
	return c.Render(200, r.HTML("auth/login.html"))
}

// AuthCallbackLine default implementation.
func AuthCallbackLine(c buffalo.Context) error {
	if c.Param("error_code") != "" {
		return c.Redirect(307, "authLoginPath()")
	} else {
		v := url.Values{}
		v.Add("grant_type", "authorization_code")
		v.Add("code", c.Param("code"))
		v.Add("redirect_uri", LINE_CALLBACK)
		v.Add("client_id", LINE_CHANNEL_ID)
		v.Add("client_secret", LINE_CHANNEL_SECRET)
		resp, err := http.Post(LINE_TOKEN_URL,
			"application/x-www-form-urlencoded",
			strings.NewReader(v.Encode()))
		if err != nil {
			c.Logger().Debug(err)
			return c.Redirect(307, "authLoginPath()")
		}

		// Not successful request.
		if resp.StatusCode != 200 {
			c.Logger().Debug(resp.StatusCode)
			return c.Redirect(307, "authLoginPath()")
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.Logger().Debug(err)
			return c.Redirect(307, "authLoginPath()")
		}

		var lineTokenResp LineTokenResp
		err = json.Unmarshal(body, &lineTokenResp)
		if err != nil {
			c.Logger().Debug(err)
			return c.Redirect(307, "authLoginPath()")
		}

		// Getting Line profile via Line API
		profile, err := line.GetProfile(lineTokenResp.AccessToken)
		if err != nil {
			c.Logger().Debug(err)
			return c.Redirect(307, "authLoginPath()")
		}

		// Set sessions for user logged in
		c.Session().Set("user", profile.ID)
		c.Session().Set("is_admin", false)

		return c.Redirect(307, "rootPath()")
	}
}

// AuthDoLogin default implementation.
func AuthDoLogin(c buffalo.Context) error {
	username := c.Param("username")
	password := c.Param("password")

	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)

	// find a user by name and password
	err := tx.Where("name = ? and password = ?", username, password).First(u)
	if err != nil {
		c.Logger().Info(err)
	} else {
		// Set sessions for user logged in
		c.Session().Set("user", u.ID)
		c.Session().Set("is_admin", u.IsAdmin)
	}
	return c.Redirect(307, "rootPath()")
}

// AuthLogout default implementation.
func AuthLogout(c buffalo.Context) error {
	// Remove sessions for user logged out
	c.Session().Delete("is_admin")
	c.Session().Delete("user")
	return c.Redirect(307, "authLoginPath()")
}
