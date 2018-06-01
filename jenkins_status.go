package main

import (
	"fmt"
	"github.com/dylanowen/btt-scripts/btt"
	"github.com/dylanowen/btt-scripts/jenkins"
	"github.com/dylanowen/btt-scripts/utils"
	"github.com/kyokomi/emoji"
	"image/color"
	"os"
	"strings"
)

type statusOutput struct {
	Text  string
	Color color.RGBA
}

var statusMap = map[string]statusOutput{
	"SUCCESS": {
		Text:  "",
		Color: color.RGBA{0x6a, 0xa5, 0x34, 0xff},
	},
	"UNSTABLE": {
		Text:  emojize("question") + "Unstable",
		Color: color.RGBA{0xf5, 0xab, 0x42, 0xff},
	},
	"FAILURE": {
		Text:  emojize("bangbang") + "Failed",
		Color: color.RGBA{0xcf, 0x42, 0x49, 0xff},
	},
	"NOT_BUILT": {
		Text:  "\u23F8Paused",
		Color: color.RGBA{0x94, 0x93, 0x93, 0xff},
	},
	"ABORTED": {
		Text:  emojize("heavy_multiplication_x") + "Aborted",
		Color: color.RGBA{0x94, 0x93, 0x93, 0xff},
	},
}

var buildingStatus = statusOutput{
	Text:  emojize("arrows_counterclockwise") + "Running",
	Color: color.RGBA{0x40, 0x85, 0xde, 0xff},
}

var unknownStatus = statusOutput{
	Text:  emojize("question") + "Unknown",
	Color: color.RGBA{0x94, 0x93, 0x93, 0xff},
}

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

	var jenkinsClient, err = jenkins.NewClient(rawJenkinsUrl)
	if err != nil {
		fatal(err)
	}

	var buildStatus *jenkins.BuildStatus

	err = utils.ChainErr(
		func() (err error) {
			buildStatus, err = jenkinsClient.GetBuildStatus()

			return
		},
		func() error {
			var statusMessage statusOutput

			if buildStatus.Building {
				statusMessage = buildingStatus
			} else {
				var foundStatus bool
				statusMessage, foundStatus = statusMap[buildStatus.Status]
				if !foundStatus {
					statusMessage = unknownStatus
				}
			}

			// construct our script result
			var scriptResult = &btt.ScriptResult{
				Text:            emojiPrefix + statusMessage.Text,
				BackgroundColor: &statusMessage.Color,
			}

			if resultString, err := scriptResult.ToJson().String(); err == nil {
				fmt.Print(resultString)
			} else {
				return err
			}

			return nil
		},
	)

	if err != nil {
		fatal("Error", err)
	}
}

/// the emoji library adds a padding string at the end, remove it
func emojize(emojiString string) string {
	if _, ok := emoji.CodeMap()[":"+emojiString+":"]; ok {
		return strings.Trim(emoji.Sprint(":"+emojiString+":"), emoji.ReplacePadding)
	} else {
		return emojiString
	}
}

func fatal(a ...interface{}) {
	fmt.Print("Error")
	fmt.Fprintln(os.Stderr, a...)

	os.Exit(1)
}
