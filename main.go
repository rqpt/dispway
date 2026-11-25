package main

import (
	"log"

	"github.com/rqpt/dispway/internal/display"
	"github.com/rqpt/picker"
)

func main() {
	outputMap, err := display.GetOutputMap()
	if err != nil {
		log.Fatal(err)
	}

	outputNames := display.GetOutputNames(outputMap)

	selectedOutput, err := picker.Run(outputNames)
	if err != nil {
		log.Fatal(err)
	}

	modes := outputMap[selectedOutput]

	modeStrings := make([]string, len(modes))

	for i, mode := range modes {
		modeStrings[i] = mode.String()
	}

	selectedMode, err := picker.Run(modeStrings)
	if err != nil {
		log.Fatal(err)
	}

	err = display.ConfigureOutput(selectedOutput, selectedMode)
	if err != nil {
		log.Fatal(err)
	}
}
