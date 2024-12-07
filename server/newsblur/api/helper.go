package newsblur

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetWithBody(client *http.Client, url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, nil
}

func PostWithBody(client *http.Client, url string, data url.Values) ([]byte, error) {
	resp, err := client.PostForm(url, data)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, nil
}
