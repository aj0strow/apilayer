// mailboxlayer api docs: https://mailboxlayer.com/documentation
package mailbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Client is a API Client for Mailbox Layer.
type Client struct {
	AccessKey string
	http.Client
}

// CheckResponse is the response object for Check(email) request.
type CheckResponse struct {
	Email       string  `json:"email"`
	DidYouMean  string  `json:"did_you_mean"`
	User        string  `json:"user"`
	Domain      string  `json:"domain"`
	FormatValid Bool    `json:"format_valid"`
	MXFound     Bool    `json:"mx_found"`
	SMTPCheck   Bool    `json:"smtp_check"`
	CatchAll    Bool    `json:"catch_all"`
	Role        Bool    `json:"role"`
	Disposable  Bool    `json:"disposable"`
	Free        Bool    `json:"free"`
	Score       float32 `json:"score"`
}

// Check sends a request to the mailbox layer API to get a CheckResponse.
func (c *Client) Check(email string) (*CheckResponse, error) {
	q := url.Values{}
	q.Set("access_key", c.AccessKey)
	q.Set("email", email)
	r, err := http.NewRequest("GET", "https://apilayer.net/api/check?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	w, err := c.Do(r)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(w.Body)
	if w.StatusCode == 200 {
		checkResponse := &CheckResponse{}
		err = decoder.Decode(checkResponse)
		if err != nil {
			return nil, err
		}
		return checkResponse, nil
	} else {
		errorResponse := &ErrorResponse{}
		err = decoder.Decode(errorResponse)
		if err != nil {
			return nil, err
		}
		return nil, errorResponse.Error
	}
}

// ErrorResponse is the standard API layer error object.
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   *Error `json:"error"`
}

// Error is the actual error details.
type Error struct {
	Code int    `json:"code"`
	Type string `json:"type"`
	Info string `json:"info"`
}

func (err *Error) Error() string {
	return fmt.Sprintf("apilayer/mailbox: %d %s %s", err.Code, err.Type, err.Info)
}

var _ error = (*Error)(nil)

// Bool is used to unmarshal null into a bool value.
type Bool bool

func (b *Bool) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("true")) {
		*b = true
	}
	return nil
}

var _ json.Unmarshaler = (*Bool)(nil)
