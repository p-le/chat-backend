package main

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleAuthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:2711/auth/callback/google",
		ClientID:     "866874494879-35970trihs9fsuu3c3lsj0tqarm6ntdo.apps.googleusercontent.com",
		ClientSecret: "wTJ_oQSUN-znU2JBiQ6d2iae",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
	oauthStateString = "lequangphu"
)

func loginHandler(writer http.ResponseWriter, req *http.Request) {
	seqs := strings.Split(req.URL.Path, "/")
	action := seqs[2]
	provider := seqs[3]
	switch action {
	case "login":
		switch provider {
		case "google":
			url := googleAuthConfig.AuthCodeURL(oauthStateString)
			http.Redirect(writer, req, url, http.StatusTemporaryRedirect)
		case "github":
			// TODO
		case "facebook":
			// TODO
		}
	case "callback":
		switch provider {
		case "google":
			state := req.FormValue("state")
			if state != oauthStateString {
				http.Error(
					writer,
					fmt.Sprintf("Invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state),
					http.StatusExpectationFailed,
				)
				return
			}
			code := req.FormValue("code")
			token, err := googleAuthConfig.Exchange(oauth2.NoContext, code)
			if err != nil {
				http.Error(
					writer,
					fmt.Sprintf("Code exchange failed with '%s'\n", err),
					http.StatusUnauthorized,
				)
				return
			}
			// oauthClient := googleAuthConfig.Client(oauth2.NoContext, token)
			response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
			if err != nil {
				http.Error(
					writer,
					fmt.Sprintf("Error when trying to complete auth for %s: %s", provider, err),
					http.StatusInternalServerError,
				)
				return
			}
			defer response.Body.Close()
			// contents, err := ioutil.ReadAll(response.Body)
			http.SetCookie(writer, &http.Cookie{
				Name:  "auth",
				Value: "Test",
				Path:  "/",
			})
			// TODO parse contents get name
			writer.Header().Set("Location", "/chat")
			writer.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
	default:
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "Auth action %s not supported", action)
	}
}
