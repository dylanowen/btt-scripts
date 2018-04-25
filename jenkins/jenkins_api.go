package jenkins

import (
	"encoding/json"
	"errors"
	"github.com/dylanowen/btt-scripts/utils"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	url *url.URL
}

type BuildStatus struct {
	Building bool
	Status   string
}

const resultKey = "result"
const buildingKey = "building"

func NewClient(rawJenkinsUrl string) (*Client, error) {
	if parsedUrl, err := url.Parse(rawJenkinsUrl); err != nil {
		return nil, errors.New("invalid jenkins url")
	} else {
		parsedUrl.Query().Add("tree", resultKey+","+buildingKey+",executor")

		return &Client{
			url: parsedUrl,
		}, nil
	}
}

func (c *Client) GetBuildStatus() (buildStatus *BuildStatus, err error) {
	var response *http.Response
	var body map[string]interface{}
	var building bool
	var status string

	err = utils.ChainErr(
		func() (err error) {
			response, err = http.Get(c.url.String())

			return
		},
		func() error {
			bytes, err := ioutil.ReadAll(response.Body)
			defer response.Body.Close()
			if err != nil {
				return err
			}
			return json.Unmarshal(bytes, &body)
		},
		func() error {
			rawBuilding, hasBuilding := body[buildingKey]
			if !hasBuilding || rawBuilding == nil {
				building = false
			} else {
				switch rawBuilding.(type) {
				case bool:
					building = rawBuilding.(bool)
				default:
					return errors.New("unexpected type of " + buildingKey)
				}
			}

			return nil
		},
		func() error {
			rawStatus, hasStatus := body[resultKey]
			if !hasStatus || rawStatus == nil {
				status = ""
			} else {
				switch rawStatus.(type) {
				case string:
					status = rawStatus.(string)
				default:
					return errors.New("unexpected type of " + resultKey)
				}
			}

			return nil
		},
		func() error {
			buildStatus = &BuildStatus{
				Building: building,
				Status:   status,
			}

			return nil
		},
	)

	return
}
