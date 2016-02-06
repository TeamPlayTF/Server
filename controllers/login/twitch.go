package login

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/TF2Stadium/Helen/config"
	"github.com/TF2Stadium/Helen/controllers/controllerhelpers"
	"github.com/TF2Stadium/Helen/models"
	"golang.org/x/net/xsrftoken"
)

type reply struct {
	AccessToken string   `json:"access_token"`
	Scope       []string `json:"scope"`
}

func TwitchLogin(w http.ResponseWriter, r *http.Request) {
	session, err := controllerhelpers.GetSessionHTTP(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	steamID, ok := session.Values["steam_id"]
	if !ok {
		http.Error(w, "You are not logged in.", http.StatusUnauthorized)
		return
	}

	player, _ := models.GetPlayerBySteamID(steamID.(string))

	loginURL := url.URL{
		Scheme: "https",
		Host:   "api.twitch.tv",
		Path:   "kraken/oauth2/authorize",
	}

	twitchRedirectURL := "http://" + config.Constants.ListenAddress + "/" + "twitchAuth"

	values := loginURL.Query()
	values.Set("response_type", "code")
	values.Set("client_id", config.Constants.TwitchClientID)
	values.Set("redirect_uri", twitchRedirectURL)
	values.Set("scope", "channel_check_subscription user_subscriptions channel_subscriptions")
	values.Set("state", xsrftoken.Generate(config.Constants.CookieStoreSecret, player.SteamID, "GET"))
	loginURL.RawQuery = values.Encode()

	http.Redirect(w, r, loginURL.String(), http.StatusTemporaryRedirect)
}

func TwitchAuth(w http.ResponseWriter, r *http.Request) {
	session, err := controllerhelpers.GetSessionHTTP(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	steamID, ok := session.Values["steam_id"]
	if !ok {
		http.Error(w, "You are not logged in.", http.StatusUnauthorized)
	}

	player, _ := models.GetPlayerBySteamID(steamID.(string))

	values := r.URL.Query()
	code := values.Get("code")
	if code == "" {
		http.Error(w, "No code given", http.StatusBadRequest)
		return
	}

	state := values.Get("state")
	if state == "" || !xsrftoken.Valid(state, config.Constants.CookieStoreSecret, player.SteamID, "GET") {
		http.Error(w, "Missing or Invalid XSRF token", http.StatusBadRequest)
		return
	}

	twitchRedirectURL := "http://" + config.Constants.ListenAddress + "/" + "twitchAuth"

	// successful login, try getting access token now
	tokenURL := url.URL{
		Scheme: "https",
		Host:   "api.twitch.tv",
		Path:   "kraken/oauth2/token",
	}
	values = tokenURL.Query()
	values.Set("client_id", config.Constants.TwitchClientID)
	values.Set("client_secret", config.Constants.TwitchClientSecret)
	values.Set("grant_type", "authorization_code")
	values.Set("redirect_uri", twitchRedirectURL)
	values.Set("code", code)
	values.Set("state", state)

	req, err := http.NewRequest("POST", tokenURL.String(), strings.NewReader(values.Encode()))
	if err != nil {
		logrus.Error(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	reply := reply{}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&reply)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	player.TwitchAccessToken = reply.AccessToken
	player.Save()

	http.Redirect(w, r, config.Constants.LoginRedirectPath, http.StatusTemporaryRedirect)
}