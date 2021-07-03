// Format output based on results from Instagram's API
package instagram

import (
	"encoding/xml"
	"fmt"
	"time"
)

type RSSContent struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Title         string    `xml:"title"`
		Link          string    `xml:"link"`
		Description   string    `xml:"description"`
		PubDateStr    string    `xml:"pubDate"`
		LastBuildDate string    `xml:"lastBuildDate"`
		Language      string    `xml:"language"`
		ProfileName   string    `xml:"profileName"`
		ProfilePicUrl string    `json:"profile_pic_url_hd"`
		Items         []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string         `xml:"title"`
	Link        string         `xml:"link"`
	Description string         `xml:"description,omitempty"`
	PubDate     int            `xml:"-"`
	PubDateStr  string         `xml:"pubDate"`
	MediaItems  []RSSMediaItem `xml:"mediaItem"`
	DisplayUrl  string         `xml:"displayUrl"`
}

type RSSMediaItem struct {
	IsVideo    bool   `xml:"is_video"`
	DisplayUrl string `xml:"displayUrl"`
	VideoUrl   string `xml:"videoUrl,omitempty"`
}

func GenerateOutput(feed InstagramFeedResponse) (RSSContent, error) {
	currTimeString := time.Now().Format(time.RFC822)
	content := newInitialisedRSSContent(currTimeString)
	content.Channel.Title = fmt.Sprintf("Instagram feed of %s", feed.Data.User.Username)
	content.Channel.Description = feed.Data.User.Description
	content.Channel.Link = fmt.Sprintf("https://instagram.com/%s", feed.Data.User.Username)
	content.Channel.ProfileName = feed.Data.User.Name
	content.Channel.ProfilePicUrl = feed.Data.User.ProfilePicUrl

	latestTime := 0

	for _, edge := range feed.Data.User.Media.Edges {
		item := buildRSSItem(edge, feed.Data.User.Username)
		if item.PubDate > latestTime {
			latestTime = item.PubDate
		}
		content.Channel.Items = append(content.Channel.Items, item)
	}

	// Derive last build date based on last post date
	if latestTime != 0 {
		content.Channel.LastBuildDate = unixTimeToString(latestTime)
	} else {
		// If there are no posts, default to current time
		content.Channel.LastBuildDate = currTimeString
	}

	return content, nil
}

func newInitialisedRSSContent(timeString string) RSSContent {
	content := RSSContent{}
	content.Version = "2.0"
	content.Channel.Language = "en-us"
	content.Channel.PubDateStr = timeString
	return content
}

func buildRSSItem(post InstagramPostEdge, accountID string) RSSItem {
	var caption string
	var edges []InstagramPostEdge

	node := post.Node

	if len(node.Caption.Edges) > 0 {
		caption = node.Caption.Edges[0].Node.Text
	} else {
		caption = ""
	}

	item := RSSItem{
		Description: caption,
		Link:        fmt.Sprintf("https://instagram.com/p/%s", node.Shortcode),
		Title:       fmt.Sprintf("%s's post", accountID),
		PubDateStr:  unixTimeToString(node.CreatedAt),
		PubDate:     node.CreatedAt,
		DisplayUrl:  node.DisplayUrl,
	}

	if len(node.Sidecar.Edges) == 0 {
		edges = append(edges, post)
	} else {
		edges = node.Sidecar.Edges
	}

	for _, slide_edge := range edges {
		item.MediaItems = append(item.MediaItems, RSSMediaItem{
			IsVideo:    slide_edge.Node.IsVideo,
			VideoUrl:   slide_edge.Node.VideoUrl,
			DisplayUrl: slide_edge.Node.DisplayUrl,
		})
	}

	return item
}

func unixTimeToString(t int) string {
	return time.Unix(int64(t), 0).Format(time.RFC822)
}
