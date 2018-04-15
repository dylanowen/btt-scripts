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
	"NOT_EXECUTED":         emojize("arrows_counterclockwise") + "Waiting",
	"ABORTED":              emojize("heavy_multiplication_x") + "Aborted",
	"SUCCESS":              "",
	"IN_PROGRESS":          emojize("arrows_counterclockwise") + "Running",
	"PAUSED_PENDING_INPUT": emojize("vertical_traffic_light") + "Paused",
	"FAILED":               emojize("bangbang") + "Failed",
	"UNSTABLE":             emojize("question") + "Unstable",
}

const RESULT_KEY = "result"

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
		parsedUrl.Query().Add("tree", RESULT_KEY)

		jenkinsUrl = parsedUrl.String()
	}

	var response *http.Response
	var body map[string]interface{}
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
			rawStatus, hasStatus := body[RESULT_KEY]
			if !hasStatus {
				status = ""
			} else {
				switch rawStatus.(type) {
				case string:
					status = rawStatus.(string)
				default:
					return errors.New("Unexpected type of " + RESULT_KEY)
				}
			}

			return nil
		},
		func() error {
			statusMessage, foundStatus := statusMap[status]
			if !foundStatus {
				statusMessage = emojize("question") + "Unknown"
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
