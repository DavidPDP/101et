package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"os"
)

func main() {
	loadEnvironmentVariables()
	configureAPI()
}

type VibesVariables struct {
	SpotifyClientId     string
	SpotifyClientSecret string 
}

var vibesVariables VibesVariables

func loadEnvironmentVariables() {
	var v1 = os.Getenv("SPOTIFY_CLIENT_ID")
	var v2 = os.Getenv("SPOTIFY_CLIENT_SECRET")
	if v1 == "" || v2 == "" {
		log.Fatal("Env variables not loaded")
	}
	vibesVariables.SpotifyClientId = v1
	vibesVariables.SpotifyClientSecret = v2
}

// Starts the TCP listener over localhost:8080 and serves the handlers 
// to handle requests on incoming connections
func configureAPI() {
	http.HandleFunc("/login", spotifyAuthHandler)
	http.HandleFunc("/callback", spotifyCallbackHandler)
	http.ListenAndServe(":8080", nil)
}

// Returns the Spotify authentication URL with the code_challenge (state)
// to comply with RFC 7636 (https://oauth.net/2/pkce/) (TODO: RFC)
func spotifyAuthHandler(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin")) // cors
	
	var state = "a_ksd7d}a}sd{'?nnn235r"
	var scope = "user-read-private user-read-email"

	base, _ := url.Parse("https://accounts.spotify.com/authorize")
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", vibesVariables.SpotifyClientId)
	params.Add("scope", scope)
	params.Add("redirect_uri", "http://localhost:5500/website/")
	params.Add("state", state)
	base.RawQuery = params.Encode()

	w.Write([]byte(base.String()))

}

type SpotifyAccessData struct { 
	Code  string 
	State string 
}

type SpotifyAccessResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

var client = &http.Client{}

// Receives the authorization code (result of user authentication) to 
// request the access token. Once having the access token, proceeds to 
// request the user's profile information (which is protected with OAuth)
func spotifyCallbackHandler(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin")) // cors
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // cors

	var data SpotifyAccessData
	_ = json.NewDecoder(r.Body).Decode(&data)

	if data.State == "a_ksd7d}a}sd{'?nnn235r" {
		
		params := url.Values{
			"code": { data.Code }, 
			"redirect_uri": { "http://localhost:5500/website/" }, 
			"grant_type": { "authorization_code" },
		}

		request, _ := http.NewRequest(
			"POST", 
			"https://accounts.spotify.com/api/token", 
			strings.NewReader(params.Encode()), 
		)
		request.SetBasicAuth(vibesVariables.SpotifyClientId, vibesVariables.SpotifyClientSecret)
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		response, _ := client.Do(request)
		
		var dataResponse SpotifyAccessResponse

		defer response.Body.Close()
		json.NewDecoder(response.Body).Decode(&dataResponse)
		
		var spotifyProfile = spotifyGetProfile(dataResponse.AccessToken)
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(spotifyProfile)
		
	}

}

type SpotifyProfile struct {
	DisplayName string `json:"display_name"`
}

func spotifyGetProfile(accessToken string) SpotifyProfile {
	
	request, _ := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	request.Header.Add("Authorization", "Bearer " + accessToken)

	response, _ := client.Do(request)

	var data SpotifyProfile
	defer response.Body.Close()
	_ = json.NewDecoder(response.Body).Decode(&data)

	return data

}
