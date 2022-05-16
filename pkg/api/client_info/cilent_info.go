package client_info

type Client struct {
	clientId     string `json:"client_Id,omitempty"`
	clientSecret string `json:"client_Secret,omitempty"`
}
