package newsblur

type ApiFeed struct {
	ID          int    `json:"id"`
	Ps          int    `json:"ps"`           // positive/focus count
	Nt          int    `json:"nt"`           // neutral/unread count
	Ng          int    `json:"ng"`           // negative/hidden count
	FeedAddress string `json:"feed_address"` // link to the feed (usually .xml/.atom)
	FeedLink    string `json:"feed_link"`    // link to the website
	FeedTitle   string `json:"feed_title"`
}

type ApiStory struct {
	StoryAuthors     string `json:"story_authors"`
	StoryPermalink   string `json:"story_permalink"`
	StoryTimestamp   int64  `json:"story_timestamp,string"`
	StoryHash        string `json:"story_hash"`
	ID               string `json:"id"`
	StoryDate        string `json:"story_date"`
	ShortParsedDate  string `json:"short_parsed_date"`
	GUIDHash         string `json:"guid_hash"`
	StoryFeedID      int    `json:"story_feed_id"`
	LongParsedDate   string `json:"long_parsed_date"`
	ReadStatus       int    `json:"read_status"`
	HasModifications bool   `json:"has_modifications"`
	StoryTitle       string `json:"story_title"`
	StoryContent     string `json:"story_content"`
	Starred          bool   `json:"starred"`
}

type LoginOutput struct {
	Authenticated bool        `json:"authenticated"`
	Code          int         `json:"code"`
	Errors        interface{} `json:"errors"`
	Result        string      `json:"result"`
}

type ReaderFeedsOutputRaw struct {
	Folders []interface{} `json:"folders"`
	Feeds   []ApiFeed     `json:"feeds"`
}

func (raw ReaderFeedsOutputRaw) toOutput() (*ReaderFeedsOutput, error) {
	output := ReaderFeedsOutput{
		Folders: make([]Folder, 0),
		Feeds:   raw.Feeds,
	}

	emptyFolder := Folder{
		Title:   "",
		FeedIDs: []int{},
	}

	for _, element := range raw.Folders {
		switch element.(type) {
		case float64, float32:
			// Feed without folder
			emptyFolder.FeedIDs = append(emptyFolder.FeedIDs, int(element.(float64)))
		case map[string]interface{}:
			// Feed with folder
			folders := element.(map[string]interface{})
			for folder, feeds := range folders {
				feedIDs := []int{}
				for _, feedId := range feeds.([]interface{}) {
					feedIDs = append(feedIDs, int(feedId.(float64)))
				}

				// Add folder if it's not empty
				if len(feedIDs) > 0 {
					output.Folders = append(output.Folders, Folder{
						Title:   folder,
						FeedIDs: feedIDs,
					})
				}
			}
		}
	}

	if len(emptyFolder.FeedIDs) > 0 {
		output.Folders = append(output.Folders, emptyFolder)
	}

	return &output, nil
}

type Folder struct {
	Title   string
	FeedIDs []int
}

type ReaderFeedsOutput struct {
	Folders []Folder
	Feeds   []ApiFeed `json:"feeds"`
}

type StoriesOutput struct {
	Stories []ApiStory `json:"stories"`
}
