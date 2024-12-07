package newsblur

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// Login as an existing user.
// POST /api/login
// https://www.newsblur.com/api#/api/login
func (nb *Newsblur) Login(username, password string) (output *LoginOutput, err error) {
	formData := url.Values{
		"username": {username},
		"password": {password},
	}

	body, err := PostWithBody(nb.client, nb.Hostname+"/api/login", formData)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &output); err != nil {
		return nil, err
	}

	if !output.Authenticated {
		return nil, fmt.Errorf("Failed to login to NewsBlur. %v", output.Errors)
	}

	return output, nil
}
