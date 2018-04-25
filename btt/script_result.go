package btt

import (
	"encoding/json"
	"fmt"
	"image/color"
)

type ScriptResult struct {
	Text            string
	IconData        string
	IconPath        string
	BackgroundColor *color.RGBA
}

func (sr *ScriptResult) ToJson() *ScriptResultJson {
	var backgroundColor = ""
	if sr.BackgroundColor != nil {
		var bc = sr.BackgroundColor
		backgroundColor = fmt.Sprintf("%d,%d,%d,%d", bc.R, bc.G, bc.B, bc.A)
	}

	return &ScriptResultJson{
		Text:            sr.Text,
		IconData:        sr.IconData,
		IconPath:        sr.IconPath,
		BackgroundColor: backgroundColor,
	}
}

// {"text":"newTitle", "icon_data":"base64_icon_data", "icon_path":"file_path_to_new_icon", "background_color": "255,85,100,255"}
type ScriptResultJson struct {
	Text            string `json:"text"`
	IconData        string `json:"icon_data"`
	IconPath        string `json:"icon_path"`
	BackgroundColor string `json:"background_color"`
}

func (srj *ScriptResultJson) String() (string, error) {
	if result, err := json.Marshal(srj); err == nil {
		return string(result), nil
	} else {
		return "", err
	}
}
