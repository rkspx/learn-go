package youtube

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type Response struct {
	Items []Item `json:"items"`
}

type Item struct {
	Stats Stats `json:"statistics"`
}

type Stats struct {
	ViewCount       string `json:"viewCount"`
	SubscriberCount string `json:"subscriberCount"`
	VideoCount      string `json:"videoCount"`
}

func GetStatistics() ([]byte, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/youtube/v3/channels", nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("key", os.Getenv("YT_API_KEY"))
	q.Add("id", os.Getenv("YT_CHANNEL_ID"))
	q.Add("part", "statistics")
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return json.Marshal(response.Items[0].Stats)
}
