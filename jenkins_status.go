package main

import (
	"encoding/json"
	"fmt"
	"github.com/dylanowen/btt-scripts/utils"
	"github.com/kyokomi/emoji"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var statusMap = map[string]string{
	"NOT_EXECUTED":         emojize("arrows_counterclockwise") + "Waiting",
	"ABORTED":              emojize("heavy_multiplication_x") + "Aborted",
	"SUCCESS":              "",
	"IN_PROGRESS":          emojize("arrows_counterclockwise") + "Running",
	"PAUSED_PENDING_INPUT": emojize("double_vertical_bar") + "Paused",
	"FAILED":               emojize("red_circle") + "Failed",
	"UNSTABLE":             emojize("question") + "Unstable",
}

func main() {

	if len(os.Args) <= 1 {
		log.Fatalln("Missing jenkins url")
	}

	var jenkinsUrl = os.Args[1]
	var emojiNames = os.Args[2:]
	var emojiPrefix = ""
	for _, emojiName := range emojiNames {
		emojiPrefix += emojize(emojiName)
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
			rawStatus, hasStatus := body["status"]
			if !hasStatus {
				status = ""
			} else {
				switch rawStatus.(type) {
				case string:
					status = rawStatus.(string)
				default:
					return errors.New("Unexpected type of status")
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
		println("Error: ", err)
	}
}

/// the emoji library adds a padding string at the end, remove it
func emojize(emojiString string) string {
	return strings.Trim(emoji.Sprint(":"+emojiString+":"), emoji.ReplacePadding)
}
