package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/roblaughlin/twitch-chatbot/env"
)

// AuthFile relative path for authentication info for the bot
const AuthFile string = "auth.json"

// TwitchAuthResponse response from Twitch's authorization server
type TwitchAuthResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int      `json:"expires_in"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

// CSRFToken returns a string representation of a CSRFToken.
// Generates the token with a given byte length.
func CSRFToken(length uint64) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)

	return base64.URLEncoding.EncodeToString(bytes), err
}

func main() {
	env, err := env.Validate("auth.env", []string{
		"AUTH_ADDR",
		"AUTH_PORT",
		"AUTHORIZE_API",
		"CLIENT_ID",
		"CLIENT_SECRET",
		"REDIRECT_URI",
		"RESPONSE_TYPE",
		"SCOPE",
		"GRANT_TYPE",
		"TOKEN_API",
	})

	if err != nil {
		log.Fatal(err)
	}

	redirectParams := url.Values{}
	redirectParams.Set("client_id", env["CLIENT_ID"])
	redirectParams.Set("redirect_uri", env["REDIRECT_URI"])
	redirectParams.Set("response_type", env["RESPONSE_TYPE"])
	redirectParams.Set("scope", env["SCOPE"])
	redirectParams.Set("force_verify", env["FORCE_VERIFY"])

	http.HandleFunc("/verify", func(writer http.ResponseWriter, reader *http.Request) {
		verify(writer, reader, env["AUTHORIZE_API"], redirectParams)
	})

	authorizationParams := url.Values{}
	authorizationParams.Set("client_id", env["CLIENT_ID"])
	authorizationParams.Set("client_secret", env["CLIENT_SECRET"])
	authorizationParams.Set("redirect_uri", env["REDIRECT_URI"])
	authorizationParams.Set("grant_type", env["GRANT_TYPE"])

	http.HandleFunc("/oauth/redirect", func(writer http.ResponseWriter, reader *http.Request) {
		oauth(writer, reader, env["TOKEN_API"], authorizationParams)
	})

	http.HandleFunc("/oauth/authorized", authorized)

	address := fmt.Sprintf("%s:%s", env["AUTH_ADDR"], env["AUTH_PORT"])
	fmt.Printf("Starting HTTP Server at %s...\n", address)
	fmt.Printf("Head to %s/verify to authenticate!\n", address)
	http.ListenAndServe(address, nil)
}

func verify(writer http.ResponseWriter, reader *http.Request, endpoint string, params url.Values) {
	/*
		state, err := CSRFToken(256)
		if err != nil {
			log.Fatal("Could not generate CSRF token.")
		}

		params.Set("state", state)
	*/

	url := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	fmt.Printf("Redirecting to: %s...\n", url)
	http.Redirect(writer, reader, url, http.StatusPermanentRedirect)
}

func oauth(writer http.ResponseWriter, reader *http.Request, endpoint string, params url.Values) {
	query := reader.URL.Query()
	code := query.Get("code")
	if code == "" {
		log.Fatal("No code received in URL params.")
	}

	params.Set("code", code)

	resp, err := http.PostForm(endpoint, params)

	if err != nil {
		fmt.Println("There was an error with the final authorization phase.")
	} else {
		defer resp.Body.Close()

		fmt.Println("Decoding response into JSON...")
		var JSONResponse TwitchAuthResponse
		err := json.NewDecoder(resp.Body).Decode(&JSONResponse)

		if err != nil {
			log.Fatal("There was a problem decoding the JSON response.")
		}

		fmt.Printf("Saving: %s...\n", AuthFile)
		file, _ := json.MarshalIndent(JSONResponse, "", "    ")
		ioutil.WriteFile(AuthFile, file, 0644)
		fmt.Printf("You may now close the server.\n\n")
	}
}

func authorized(writer http.ResponseWriter, reader *http.Request) {
	fmt.Println("success")
}
