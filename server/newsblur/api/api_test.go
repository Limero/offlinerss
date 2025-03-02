package api

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed response/readerfeeds.json
var respReaderFeeds []byte

func TestLogin(t *testing.T) {
}

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

	assert.Equal(t, []Feed{
		{ID: 8312062, Ps: 0, Nt: 0, Ng: 0, FeedAddress: "https://webzine.puffy.cafe/atom.xml", FeedLink: "https://webzine.puffy.cafe/", FeedTitle: "OpenBSD Webzine"},
		{ID: 7600810, Ps: 0, Nt: 0, Ng: 0, FeedAddress: "https://undeadly.org/cgi?action=rss", FeedLink: "https://undeadly.org/", FeedTitle: "OpenBSD Journal"},
	}, output.Feeds)
}

func TestReaderUnreadStoryHashes(t *testing.T) {
}

func TestReaderStarredStoryHashes(t *testing.T) {
}

func TestReaderRiverStories_StoryHash(t *testing.T) {
}

func TestMarkStoryHashesAsRead(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/reader/mark_story_hashes_as_read", r.RequestURI)
		require.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		require.NoError(t, r.ParseForm())
		require.Equal(t, url.Values{"story_hash": []string{"a", "b"}}, r.Form)

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("")) // TODO
		require.NoError(t, err)
	}))
	defer ts.Close()

	api := &Newsblur{
		Hostname: ts.URL,
		client:   &http.Client{},
	}
	require.NoError(t, api.MarkStoryHashesAsRead([]string{"a", "b"}))
}

func TestMarkStoryHashAsUnread(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/reader/mark_story_hash_as_unread", r.RequestURI)
		require.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		require.NoError(t, r.ParseForm())
		require.Equal(t, url.Values{"story_hash": []string{"a", "b"}}, r.Form)

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(""))
		require.NoError(t, err)
	}))
	defer ts.Close()

	api := &Newsblur{
		Hostname: ts.URL,
		client:   &http.Client{},
	}
	require.NoError(t, api.MarkStoryHashAsUnread([]string{"a", "b"}))
}

func TestMarkStoryHashAsStarred(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/reader/mark_story_hash_as_starred", r.RequestURI)
		require.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		require.NoError(t, r.ParseForm())
		require.Equal(t, url.Values{"story_hash": []string{"a"}}, r.Form)

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(""))
		require.NoError(t, err)
	}))
	defer ts.Close()

	api := &Newsblur{
		Hostname: ts.URL,
		client:   &http.Client{},
	}
	require.NoError(t, api.MarkStoryHashAsStarred("a"))
}

func MarkStoryHashAsUnstarred(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/reader/mark_story_hash_as_unstarred", r.RequestURI)
		require.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		require.NoError(t, r.ParseForm())
		require.Equal(t, url.Values{"story_hash": []string{"a", "b"}}, r.Form)

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(""))
		require.NoError(t, err)
	}))
	defer ts.Close()

	api := &Newsblur{
		Hostname: ts.URL,
		client:   &http.Client{},
	}
	require.NoError(t, api.MarkStoryHashAsUnstarred([]string{"a", "b"}))
}
