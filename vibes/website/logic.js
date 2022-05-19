const spotifyRequest = new Request('http://localhost:8080/login', {method: 'GET'});

/** 
 * Starts Spotify OAuth 2.0 Flow, calling the Vibes backend to 
 * retrieve the Spotify authentication URL to redirect you. 
 */
function spotifyAuth() {
    fetch(spotifyRequest)
        .then(resp => resp.text())
        .then(spotifyAuthUrl => { 
            window.location.assign(spotifyAuthUrl);
         }).catch(err => { console.log(err) });
}

const spotifyBtn = document.getElementById('spotify-auth'); 
spotifyBtn.addEventListener('click', () => { spotifyAuth() });

/**
 * Checks the query parameters in the current URL. If found it, calls 
 * the Vibes backend to submit the authentication code (to continue 
 * the OAuth 2.0 flow). Once the flow ends, receives the user's profile data.
 */
function vibesAuth() {
    if(window.location.search) {
        
        const queryParams = new URLSearchParams(window.location.search);
        window.history.replaceState({}, document.title, "/website/"); // clean query params
      
        const vibesRequest = new Request('http://localhost:8080/callback', 
            { method: 'POST', headers: {'Content-Type': 'application/json'}, 
              body: JSON.stringify({ Code: queryParams.get('code'), State:  queryParams.get('state') })
            }
        );

        fetch(vibesRequest)
            .then(resp => resp.text())
            .then(userProfile => {
                document.getElementById('spotify-auth').style.setProperty('display', 'none');
                const userInfo = document.getElementById('spotify-user');
                userInfo.style.setProperty('display', 'block');
                userInfo.textContent = "@" + JSON.parse(userProfile).display_name;
            }).catch(err => { console.log(err) });

    }
}

vibesAuth();