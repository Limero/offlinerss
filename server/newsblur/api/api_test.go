package newsblur

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed response/readerfeeds.json
var respReaderFeeds []byte

func TestReaderFeeds(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/reader/feeds?v=2", r.RequestURI)

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(respReaderFeeds))
		require.NoError(t, err)
	}))
	defer ts.Close()

	api := &Newsblur{
		Hostname: ts.URL,
		client:   &http.Client{},
	}
	output, err := api.ReaderFeeds()
	require.NoError(t, err)

	assert.Len(t, output.Folders, 2)
	assert.Equal(t, []Folder{
		{
			Title: "BSD",
			FeedIDs: []int{
				7600810,
				8312062,
			},
		},
		{
			Title: "",
			FeedIDs: []int{
				1000000,
				1111111,
			},
		},
	}, output.Folders)

	assert.Equal(t, []ApiFeed{
		{ID: 8312062, Ps: 0, Nt: 0, Ng: 0, FeedAddress: "https://webzine.puffy.cafe/atom.xml", FeedLink: "https://webzine.puffy.cafe/", FeedTitle: "OpenBSD Webzine"},
		{ID: 7600810, Ps: 0, Nt: 0, Ng: 0, FeedAddress: "https://undeadly.org/cgi?action=rss", FeedLink: "https://undeadly.org/", FeedTitle: "OpenBSD Journal"},
	}, output.Feeds)
}
