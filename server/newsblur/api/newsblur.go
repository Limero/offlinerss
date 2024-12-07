package newsblur

import (
	"net/http"
)

type Newsblur struct {
	Hostname string
	client   *http.Client
}

func New(client *http.Client) *Newsblur {
	return &Newsblur{
		Hostname: "https://www.newsblur.com",
		client:   client,
	}
}
