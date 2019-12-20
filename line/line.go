package line

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gobuffalo/envy"
)

type LineProfileResp struct {
	ID     string `json:"userId"`
	Name   string `json:"displayName"`
	Avatar string `json:"pictureUrl"`
	Status string `json:"statusMessage"`
}

var LINE_PROFILE_URL = ""

func init() {
	LINE_PROFILE_URL = envy.Get("LINE_PROFILE_URL", "")
}

func GetProfile(accessToken string) (*LineProfileResp, error) {
	var profile *LineProfileResp
	client := &http.Client{}
	req, err := http.NewRequest("GET", LINE_PROFILE_URL, nil)
	if err != nil {
		return profile, err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return profile, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return profile, err
	}

	err = json.Unmarshal(body, &profile)
	if err != nil {
		return profile, err
	}

	return profile, err
}
