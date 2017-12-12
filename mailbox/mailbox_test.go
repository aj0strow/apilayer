package mailbox

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	client := &Client{
		AccessKey: "123-not-real",
	}
	client.Timeout = time.Second * 10
}

func TestCheckResponse(t *testing.T) {
	tests := []struct {
		Body     string
		Response interface{}
	}{
		{
			Body: `{
    "catch_all": null,
    "did_you_mean": "",
    "disposable": false,
    "domain": "gmail.com",
    "email": "alexander.ostrow@gmail.com",
    "format_valid": true,
    "free": true,
    "mx_found": true,
    "role": false,
    "score": 0.8,
    "smtp_check": true,
    "user": "alexander.ostrow"
}`,
			Response: &CheckResponse{
				Domain:      "gmail.com",
				Email:       "alexander.ostrow@gmail.com",
				FormatValid: true,
				Free:        true,
				MXFound:     true,
				Role:        false,
				Score:       0.8,
				SMTPCheck:   true,
				User:        "alexander.ostrow",
			},
		},
		{
			Body: `{
    "catch_all": null,
    "did_you_mean": "alexander.ostrow@gmail.com",
    "disposable": false,
    "domain": "gail.com",
    "email": "alexander.ostrow@gail.com",
    "format_valid": true,
    "free": false,
    "mx_found": true,
    "role": false,
    "score": 0.32,
    "smtp_check": false,
    "user": "alexander.ostrow"
}`,
			Response: &CheckResponse{
				DidYouMean:  "alexander.ostrow@gmail.com",
				Domain:      "gail.com",
				Email:       "alexander.ostrow@gail.com",
				FormatValid: true,
				MXFound:     true,
				Score:       0.32,
				SMTPCheck:   false,
				User:        "alexander.ostrow",
			},
		},
		{
			Body: `{
    "catch_all": null,
    "did_you_mean": "",
    "disposable": false,
    "domain": "gmail",
    "email": "alexander.ostrow@gmail",
    "format_valid": false,
    "free": false,
    "mx_found": null,
    "role": false,
    "score": 0.64,
    "smtp_check": false,
    "user": "alexander.ostrow"
}`,
			Response: &CheckResponse{
				Domain: "gmail",
				Email:  "alexander.ostrow@gmail",
				Score:  0.64,
				User:   "alexander.ostrow",
			},
		},
	}
	for _, tt := range tests {
		checkResponse := &CheckResponse{}
		err := json.Unmarshal([]byte(tt.Body), checkResponse)
		if err != nil {
			t.Errorf(err.Error())
		}
		if !reflect.DeepEqual(checkResponse, tt.Response) {
			t.Errorf("invalid check response", tt.Response)
		}
	}
}

func TestErrorResponse(t *testing.T) {
	tests := []struct {
		Body     string
		Response interface{}
	}{
		{
			Body: `{
    "error": {
        "code": 210,
        "info": "Please specify an email address. [Example: support@apilayer.com]",
        "type": "no_email_address_supplied"
    },
    "success": false
}`,
			Response: &ErrorResponse{
				Error: &Error{
					Code: 210,
					Info: "Please specify an email address. [Example: support@apilayer.com]",
					Type: "no_email_address_supplied",
				},
			},
		},
	}
	for _, tt := range tests {
		response := &ErrorResponse{}
		err := json.Unmarshal([]byte(tt.Body), response)
		if err != nil {
			t.Errorf(err.Error())
		}
		if !reflect.DeepEqual(response, tt.Response) {
			t.Errorf("invalid check response", tt.Response)
		}
	}
}
