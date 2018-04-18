package main

import (
	"encoding/json"
	"fmt"
	"github.com/dylanowen/btt-scripts/utils"
	"github.com/kyokomi/emoji"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var statusMap = map[string]string{
	"SUCCESS":   "",
	"UNSTABLE":  emojize("question") + "Unstable",
	"FAILURE":   emojize("bangbang") + "Failed",
	"NOT_BUILT": "\u23F8Paused",
	"ABORTED":   emojize("heavy_multiplication_x") + "Aborted",
}

var BuildingStatus = emojize("arrows_counterclockwise") + "Running"

const ResultKey = "result"
const BuildingKey = "building"

func main() {

	if len(os.Args) <= 1 {
		fatal("Missing jenkins url")
	}

	var rawJenkinsUrl = os.Args[1]
	var emojiNames = os.Args[2:]
	var emojiPrefix = ""
	for _, emojiName := range emojiNames {
		emojiPrefix += emojize(emojiName)
	}

	var jenkinsUrl string
	if parsedUrl, err := url.Parse(rawJenkinsUrl); err != nil {
		fatal("Invalid jenkins url")
	} else {
		parsedUrl.Query().Add("tree", ResultKey+","+BuildingKey+",executor")

		jenkinsUrl = parsedUrl.String()
	}

	var response *http.Response
	var body map[string]interface{}
	var building bool
	var status string

	var err = utils.ChainErr(
		func() (err error) {
			response, err = http.Get(jenkinsUrl)

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
			rawBuilding, hasBuilding := body[BuildingKey]
			if !hasBuilding || rawBuilding == nil {
				building = false
			} else {
				switch rawBuilding.(type) {
				case bool:
					building = rawBuilding.(bool)
				default:
					return errors.New("Unexpected type of " + BuildingKey)
				}
			}

			return nil
		},
		func() error {
			rawStatus, hasStatus := body[ResultKey]
			if !hasStatus || rawStatus == nil {
				status = ""
			} else {
				switch rawStatus.(type) {
				case string:
					status = rawStatus.(string)
				default:
					return errors.New("Unexpected type of " + ResultKey)
				}
			}

			return nil
		},
		func() error {
			var statusMessage string

			if building {
				statusMessage = BuildingStatus
			} else {
				var foundStatus bool
				statusMessage, foundStatus = statusMap[status]
				if !foundStatus {
					statusMessage = emojize("question") + "Unknown"
				}
			}

			fmt.Print(emojiPrefix, statusMessage)

			return nil
		},
	)

	if err != nil {
		fatal("Error", err)
	}
}

/// the emoji library adds a padding string at the end, remove it
func emojize(emojiString string) string {
	return strings.Trim(emoji.Sprint(":"+emojiString+":"), emoji.ReplacePadding)
}

func fatal(a ...interface{}) {
	fmt.Print("Error")
	fmt.Fprintln(os.Stderr, a...)

	os.Exit(1)
}
