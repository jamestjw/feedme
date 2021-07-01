// Use Instagram API to fetch posts

package instagram

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

type InstagramFeedResponse struct {
	Data struct {
		User struct {
			Name          string `json:"full_name"`
			Username      string `json:"username"`
			IsPrivate     bool   `json:"is_private"`
			Description   string `json:"biography"`
			ProfilePicUrl string `json:"profile_pic_url_hd"`
			Media         struct {
				Count int                 `json:"count"`
				Edges []InstagramPostEdge `json:"edges"`
			} `json:"edge_owner_to_timeline_media"`
		} `json:"user"`
	} `json:"graphql"`
}

type InstagramPostEdge struct {
	Node InstagramPostEdgeNode `json:"node"`
}

type InstagramPostEdgeNode struct {
	Shortcode  string `json:"shortcode"`
	DisplayUrl string `json:"display_url"`
	VideoUrl   string `json:"video_url"`
	IsVideo    bool   `json:"is_video"`
	CreatedAt  int    `json:"taken_at_timestamp"`
	Caption    struct {
		Edges []InstagramCaptionEdge `json:"edges"`
	} `json:"edge_media_to_caption"`
	Sidecar struct {
		Edges []InstagramPostEdge `json:"edges"`
	} `json:"edge_sidecar_to_children"`
}

type InstagramCaptionEdge struct {
	Node struct {
		Text string `json:"text"`
	} `json:"node"`
}

func FetchFeed(accountID *string) (InstagramFeedResponse, error) {
	igResponse := InstagramFeedResponse{}

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.instagram.com/%s/channel/?__a=1", *accountID), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", "feedme-cli-tool/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return igResponse, errors.New(fmt.Sprintf("Failed to fetch instagram feeds for account ID %s", *accountID))
	}

	if err != nil {
		return igResponse, errors.New(fmt.Sprintf("Failed to fetch instagram feeds for account ID %s", *accountID))
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return igResponse, errors.New(fmt.Sprintf("Failed to parse instagram feeds for account ID %s", *accountID))
	}

	jsonErr := json.Unmarshal(body, &igResponse)
	if jsonErr != nil {
		return igResponse, errors.New(fmt.Sprintf("Failed to parse instagram feeds for account ID %s", *accountID))
	}

	if reflect.ValueOf(igResponse).IsZero() {
		return igResponse, errors.New(fmt.Sprintf("Instagram account %s doesn't exist", *accountID))
	}

	return igResponse, nil
}
