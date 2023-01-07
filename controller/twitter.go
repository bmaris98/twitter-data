package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bbalet/stopwords"
	"mvdan.cc/xurls/v2"
)

var urlFormat = "https://api.twitter.com/1.1/search/tweets.json?q=%s&lang=en&count=10&tweet_mode=extended&result_type=mixed"
var bearerHeader = "Bearer " + os.Getenv("TWITTER_BEARER")
var tmpRoot = "./tmp"

type SearchResult struct {
	Statuses []Tweet `json:"statuses"`
}

func (sr SearchResult) GetMaximumId() float64 {
	var max float64 = 0
	for _, v := range sr.Statuses {
		if v.Id > max {
			max = v.Id
		}
	}
	return max
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

func (data PostData) ToWordList() string {
	allWords := append(append(data.Words, data.Mentions...), data.ManualTags...)
	return strings.Join(allWords, " ")
}

func RunQuery(query string, lastReadId float64, client *http.Client, kafkaCtx KafkaContext, hdfsConn HdfsConnection) float64 {
	url := fmt.Sprintf(urlFormat, url.QueryEscape(query))
	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", bearerHeader)

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

	if searchResult.GetMaximumId() == lastReadId {
		log.Println("No new tweets for prompt '%s'", query)
		return lastReadId
	} else {
		log.Println("Got new data for prompt '%s'", query)
	}

	postData := MapTweetsToPostData(searchResult)

	csvLines := Map(postData, func(single PostData) string {
		return single.ToCsv()
	})
	csv := strings.Join(csvLines, "\n")
	kafkaCtx.PushMsg(query, csv)

	wordsLines := Map(postData, func(single PostData) string {
		return single.ToWordList()
	})
	allWords := strings.Join(wordsLines, "\n")

	hdfsFileName := fmt.Sprintf("%d.txt", time.Now().UnixNano())
	hdfsConn.CreateFile(query, hdfsFileName, allWords)

	return searchResult.GetMaximumId()
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

func RemoveSubstringsFromText(text string, substrings []string) string {
	for _, v := range substrings {
		text = strings.ReplaceAll(text, v, "")
	}
	return text
}
