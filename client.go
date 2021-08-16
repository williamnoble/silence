package silence

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const baseURL string = "https://hacker-news.firebaseio.com/v0"
const linkURL string = "https://news.ycombinator.com/item?id="

type Client struct {
}

type NewsItem struct {
	Title string
	URL   string
	Host  string
	ID    int
	Link  string `json:"-"`
}

// RetrieveTopStories queries the HackerNews API and retrieves a list of
// HN NewsItems by their ID.
func (c Client) RetrieveTopStories() ([]int, error) {
	var idList []int
	err := request(baseURL, "/topstories.json", &idList)
	if err != nil {
		return []int{}, err
	}

	return idList, nil
}

// RetrieveNewsItemById retrieves a single HackerNews Item by its ID.
func (c Client) RetrieveNewsItemById(id int) (NewsItem, error) {
	newsItem := NewsItem{}
	pathWithId := fmt.Sprintf("/item/%d.json", id)

	err := request(baseURL, pathWithId, &newsItem)
	if err != nil {
		return NewsItem{}, err
	}

	parsedURL, err := url.Parse(newsItem.URL)
	if err != nil {
		return NewsItem{}, err
	}

	newsItem.Host = strings.TrimPrefix(parsedURL.Hostname(), "www.")
	newsItem.ID = id
	newsItemString := strconv.Itoa(id)
	newsItem.Link = fmt.Sprintf("%s%s", linkURL, newsItemString)

	return newsItem, nil
}

func request(baseURL string, path string, target interface{}) error {
	fullPath := fmt.Sprintf("%s%s", baseURL, path)

	data, err := http.Get(fullPath)
	if err != nil {
		return fmt.Errorf("failed to decode JSON for url: %s", fullPath)
	}

	dec := json.NewDecoder(data.Body)
	err = dec.Decode(target)
	if err != nil {
		return fmt.Errorf("failed to decode JSON for url: %s", fullPath)
	}
	return nil
}
