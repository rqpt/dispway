package main

import (
	"log"

	"github.com/rqpt/dispway/internal/display"
	"github.com/rqpt/picker"
)

func main() {
	outputs, err := display.GetOutputs()
	if err != nil {
		log.Fatal(err)
	}

	outputNames := display.GetOutputNames(outputs)

	selectedOutput, err := picker.Run(outputNames)
	if err != nil {
		log.Fatal(err)
	}

	var modes []display.Mode

	for _, output := range outputs {
		if output.Name == selectedOutput {
			modes = output.Modes
			break
		}
	}

	modes = display.GetUniqueModes(modes)

	modeStrings := make([]string, len(modes))

	for i, mode := range modes {
		modeStrings[i] = mode.String()
	}

	selectedModeString, err := picker.Run(modeStrings)
	if err != nil {
		log.Fatal(err)
	}

	var selectedMode display.Mode

	for _, mode := range modes {
		if mode.String() == selectedModeString {
			selectedMode = mode
			break
		}
	}

	err = display.ConfigureOutput(selectedOutput, selectedMode)
	if err != nil {
		log.Fatal(err)
	}
}
