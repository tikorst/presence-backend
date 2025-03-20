package mobile

type LoginResult struct {
    Name       string `json:"name"`
    Username   string `json:"username"`
    TokenString string `json:"token_string"`
}

type LoginResponse struct {
    LoginResult *LoginResult `json:"loginResult"`
    Error       bool         `json:"error"`
    Message     string       `json:"message"`
}
