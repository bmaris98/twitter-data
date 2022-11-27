package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bbalet/stopwords"
	"mvdan.cc/xurls/v2"
)

var urlFormat = "https://api.twitter.com/1.1/search/tweets.json?q=%s&lang=en&count=10&tweet_mode=extended&result_type=recent"
var bearerHeader = "Bearer " + os.Getenv("TWITTER_BEARER")
var tmpRoot = "./tmp"

type SearchResult struct {
	Statuses []Tweet `json:"statuses"`
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

type Tweet struct {
	Id       float64 `json:"id"`
	FullText string  `json:"full_text"`
	Entities struct {
		Hashtags []Hashtag `json:"hashtags"`
		Mentions []Mention `json:"user_mentions"`
	} `json:"entities"`
}

type Hashtag struct {
	Text string `json:"text"`
}

type Mention struct {
	ScreenName string `json:"screen_name"`
}

type PostData struct {
	Words      []string `json:"words"`
	Mentions   []string `json:"mentions"`
	ManualTags []string `json:"manual_tags"`
}

func (data PostData) ToCsv() string {
	return fmt.Sprintf("%s~|~%s~|~%s", strings.Join(data.Words, ","), strings.Join(data.Mentions, ","), strings.Join(data.ManualTags, ","))
}

type void struct{}

var nullptr void

func main() {
	client := &http.Client{}
	active_queries := make(map[string]void)
	active_queries["%23elections"] = nullptr

	for {
		for query := range active_queries {
			RunQuery(query, client)
		}
		time.Sleep(10 * time.Second)
	}
}

func RunQuery(query string, client *http.Client) {
	// Create a new request using http
	url := fmt.Sprintf(urlFormat, query)
	req, err := http.NewRequest("GET", url, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearerHeader)

	// Send req using http Client
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	var searchResult SearchResult
	err = json.Unmarshal(body, &searchResult)
	if err != nil {
		log.Print("Error deserializing.\n - ", err)
	}

	postData := MapTweetsToPostData(searchResult)

	lines := Map(postData, func(single PostData) string {
		return single.ToCsv()
	})
	csv := strings.Join(lines, "\n")
	log.Println(csv)
	// marshalled, err := json.Marshal(postData)

	dirPath := fmt.Sprintf("%s/%s", tmpRoot, query)
	err = os.MkdirAll(dirPath, 0777)
	Check(err)
	fileName := fmt.Sprintf("%s/%d.csv", dirPath, time.Now().UnixNano())
	err = os.WriteFile(fileName, []byte(csv), 0777)
	Check(err)
}

func MapTweetsToPostData(searchResult SearchResult) []PostData {
	var postData []PostData
	for _, v := range searchResult.Statuses {
		originalTags := ExtractTagsFromTweet(&v)
		lowerCaseTags := Map(originalTags, func(s string) string {
			return strings.ToLower(s)
		})

		mentions := ExtractMentionsFromTweet(&v)

		sanitizedText := SanitizeTwitterText(v.FullText, originalTags, mentions)
		var words []string

		for _, v := range strings.Split(stopwords.CleanString(sanitizedText, "en", true), " ") {
			if len(v) > 1 {
				words = append(words, v)
			}
		}

		postData = append(postData, PostData{Words: words, ManualTags: lowerCaseTags, Mentions: mentions})
	}
	return postData
}

func ExtractTagsFromTweet(tweet *Tweet) []string {
	return Map(tweet.Entities.Hashtags, func(h Hashtag) string {
		return h.Text
	})
}

func ExtractMentionsFromTweet(tweet *Tweet) []string {
	return Map(tweet.Entities.Mentions, func(m Mention) string {
		return m.ScreenName
	})
}

func SanitizeTwitterText(text string, tags []string, mentions []string) string {
	urls := xurls.Strict().FindAllString(text, -1)
	formattedTags := Map(tags, func(tag string) string {
		return "#" + tag
	})
	formattedMentions := Map(tags, func(mention string) string {
		return "@" + mention
	})
	substringsToRemove := append(urls, formattedTags...)
	substringsToRemove = append(substringsToRemove, formattedMentions...)
	return RemoveSubstringsFromText(text, substringsToRemove)
}

// and or the
func RemoveSubstringsFromText(text string, substrings []string) string {
	for _, v := range substrings {
		text = strings.ReplaceAll(text, v, "")
	}
	return text
}

func Map[TIn, TOut any](input []TIn, mapper func(TIn) TOut) []TOut {
	var output []TOut
	for _, v := range input {
		output = append(output, mapper(v))
	}
	return output
}
