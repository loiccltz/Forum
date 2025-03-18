package backend

import (
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "golang.org/x/oauth2/github"
    "os"
)


//google
var googleOAuthConfig = &oauth2.Config{
    ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
    ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
    RedirectURL:  "https://localhost/auth/google/callback",
    Scopes:       []string{"email", "profile"},
    Endpoint:     google.Endpoint,
}

//gituhb
var GithubOAuthConfig = &oauth2.Config{
    ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
    ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
    RedirectURL:  "https://localhost:8080/auth/github/callback",
    Scopes:       []string{"user:email"},
    Endpoint:     github.Endpoint,
}


//mettre le reste git facebook