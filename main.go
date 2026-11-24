package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/rqpt/picker"
)

func main() {
	cmd := exec.Command("swaymsg", "-t", "get_outputs", "-r")

	cmdOutput, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	var outputs []Output

	err = json.Unmarshal(cmdOutput, &outputs)
	if err != nil {
		log.Fatal(err)
	}

	outputModes := make(map[string][]Mode)
	for _, o := range outputs {
		outputModes[o.Name] = o.Modes
	}

	outputNames := make([]string, 0, len(outputModes))
	for name := range outputModes {
		outputNames = append(outputNames, name)
	}

	selectedOutput, err := picker.Run(outputNames)
	if err != nil {
		log.Fatal(err)
	}

	modeStrings := make([]string, 0, len(outputModes[selectedOutput]))
	seen := make(map[string]bool)

	for _, m := range outputModes[selectedOutput] {
		mode := fmt.Sprintf(
			"%dx%d @ %.3f Hz",
			m.Width,
			m.Height,
			float64(m.Refresh)/1000,
		)

		if !seen[mode] {
			modeStrings = append(modeStrings, mode)
			seen[mode] = true
		}
	}

	selectedMode, err := picker.Run(modeStrings)
	if err != nil {
		log.Fatal(err)
	}

	cmd = exec.Command(
		"swaymsg",
		"output",
		selectedOutput,
		"mode",
		strings.ReplaceAll(selectedMode, " ", ""),
	)

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

type Output struct {
	Name  string `json:"name"`
	Modes []Mode `json:"modes"`
}

type Mode struct {
	Width              int    `json:"width"`
	Height             int    `json:"height"`
	Refresh            int    `json:"refresh"`
	PictureAspectRatio string `json:"picture_aspect_ratio"`
}
