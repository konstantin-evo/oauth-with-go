package model

// TokenIntrospect struct represents the response from the token introspection endpoint.
type TokenIntrospect struct {
	Exp    int         `json:"exp"`
	Nbf    int         `json:"nbf"`
	Iat    int         `json:"iat"`
	Jti    string      `json:"jti"`
	Aud    interface{} `json:"aud"`
	Typ    string      `json:"typ"`
	Acr    string      `json:"acr"`
	Active bool        `json:"active"`
}

type TokenClaim struct {
	Exp            int         `json:"exp"`
	Iat            int         `json:"iat"`
	AuthTime       int         `json:"auth_time"`
	Jti            string      `json:"jti"`
	Iss            string      `json:"iss"`
	Aud            interface{} `json:"aud"`
	Sub            string      `json:"sub"`
	Typ            string      `json:"typ"`
	Azp            string      `json:"azp"`
	SessionState   string      `json:"session_state"`
	Acr            string      `json:"acr"`
	AllowedOrigins []string    `json:"allowed-origins"`
	RealmAccess    struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	ResourceAccess struct {
		Account struct {
			Roles []string `json:"roles"`
		} `json:"account"`
	} `json:"resource_access"`
	Scope             string `json:"scope"`
	Sid               string `json:"sid"`
	EmailVerified     bool   `json:"email_verified"`
	PreferredUsername string `json:"preferred_username"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
}
