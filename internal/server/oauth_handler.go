package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/pglet/pglet/internal/auth"
	"github.com/pglet/pglet/internal/config"
	"github.com/pglet/pglet/internal/utils"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const (
	redirectUrlParameter  = "redirect_url"
	groupsUrlParameter    = "groups"
	principalIdCookieName = "pid"
)

func githubAuthHandler(c *gin.Context) {
	oauthHandler(c, auth.GitHubAuth)
}

func azureAuthHandler(c *gin.Context) {
	oauthHandler(c, auth.AzureAuth)
}

func oauthHandler(c *gin.Context, authProvider string) {
	code := c.Query("code")
	stateID := c.Query("state")

	if code == "" {
		// initial flow
		redirectURL := c.Query(redirectUrlParameter)
		groupsEnabled := c.Query(groupsUrlParameter) == "1"

		stateID, err := saveOAuthState(c.Writer, &auth.State{
			RedirectURL:   redirectURL,
			AuthProvider:  authProvider,
			GroupsEnabled: groupsEnabled,
		})

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// redirect to authorize page
		oauthConfig := auth.GetOauthConfig(authProvider, groupsEnabled)
		c.Redirect(302, oauthConfig.AuthCodeURL(stateID))
	} else {

		// load state from cookie
		if stateID == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid state"))
			return
		}

		state, err := getOAuthState(c.Request, stateID)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		oauthConfig := auth.GetOauthConfig(authProvider, state.GroupsEnabled)

		// request token
		token, err := oauthConfig.Exchange(oauth2.NoContext, code)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// create new principal and update its details from API
		principal := auth.NewPrincipal(authProvider, state.GroupsEnabled)
		principal.SetToken(token)
		err = principal.UpdateDetails()

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		log.Debugln(utils.ToJSON(principal))

		deleteOAuthState(c.Writer, stateID)
		savePrincipalID(c.Writer, principal.UID)
		c.Redirect(302, state.RedirectURL)
	}
}

func saveOAuthState(w http.ResponseWriter, state *auth.State) (string, error) {
	id, _ := utils.GenerateRandomString(32)
	state.Id = id

	sc := getSecureCookie()

	// serialize to a secure cookie
	encoded, err := sc.Encode(id, state)
	if err != nil {
		return "", err
	}

	cookie := &http.Cookie{
		Name:     id,
		Value:    encoded,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	return id, nil
}

func getOAuthState(r *http.Request, stateID string) (*auth.State, error) {
	sc := getSecureCookie()
	cookie, err := r.Cookie(stateID)

	if err != nil {
		return nil, err
	}

	state := &auth.State{}
	err = sc.Decode(stateID, cookie.Value, &state)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func deleteOAuthState(w http.ResponseWriter, stateID string) {
	cookie := &http.Cookie{
		Name:     stateID,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   -1,
	}
	http.SetCookie(w, cookie)
}

func savePrincipalID(w http.ResponseWriter, principalID string) error {
	sc := getSecureCookie()

	// serialize to a secure cookie
	encoded, err := sc.Encode(principalIdCookieName, principalID)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     principalIdCookieName,
		Value:    encoded,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	return nil
}

func getPrincipalID(r *http.Request) (string, error) {
	sc := getSecureCookie()
	cookie, err := r.Cookie(principalIdCookieName)

	if err == http.ErrNoCookie {
		return "", nil
	} else if err != nil {
		return "", err
	}

	principalID := ""
	err = sc.Decode(principalIdCookieName, cookie.Value, &principalID)
	if err != nil {
		return "", err
	}
	return principalID, nil
}

func getSecureCookie() *securecookie.SecureCookie {
	return securecookie.New(utils.GetCipherKey(config.CookieSecret()), utils.GetCipherKey(config.MasterSecretKey()))
}
