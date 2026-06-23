package wellknown

type Info struct {
	Version string   `json:"version"`
	Auth    AuthInfo `json:"auth"`
}

type AuthInfo struct {
	Providers []string `json:"providers"`
}
