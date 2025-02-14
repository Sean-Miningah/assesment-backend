package utils

type GoogleApiResponse struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	Picture       string `json:"picture"`
}

type SMSResponse struct {
	APIKey   string
	Username string
}

type SMSRequest struct {
	Username     string   `json:"username"`
	Message      string   `json:"message"`
	SenderID     string   `json:"senderId,omitempty"`
	MaskedNumber string   `json:"maskedNumber,omitempty"`
	Telco        string   `json:"telco,omitempty"`
	PhoneNumbers []string `json:"phoneNumbers"`
}
