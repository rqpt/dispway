// Package display implements abstractions for dealing with getting and setting
// the monitor display configuration of swaywm.
package display

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

// GetOutputMap calls `swaymsg -t get_outputs -r`, which outputs a json body
// describing the current display configuration, and then returns an outputMap.
func GetOutputMap() (outputMap, error) {
	cmd := exec.Command("swaymsg", "-t", "get_outputs", "-r")

	cmdOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var outputs []output

	err = json.Unmarshal(cmdOutput, &outputs)
	if err != nil {
		return nil, err
	}

	outputMap := make(outputMap, len(outputs))

	for _, output := range outputs {
		output.clearDuplicateModes()
		outputMap[output.Name] = output.Modes
	}

	return outputMap, nil
}

// GetOutputNames returns the keys of an outputMap.
func GetOutputNames(outputMap outputMap) []string {
	names := make([]string, 0, len(outputMap))

	for name := range outputMap {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

// ConfigureOutput sets the current display configuration by calling swaymsg
// using the provided output and mode.
func ConfigureOutput(output, mode string) error {
	cmd := exec.Command(
		"swaymsg",
		"output",
		output,
		"mode",
		strings.ReplaceAll(mode, " ", ""),
	)

	return cmd.Run()
}

// A map keyed by the display outputs available, and the various modes
// available for each as values.
type outputMap map[string][]mode

type output struct {
	Name  string `json:"name"`
	Modes []mode `json:"modes"`
}

func (o *output) clearDuplicateModes() {
	seen := make(map[mode]struct{})
	unique := make([]mode, 0, len(o.Modes))

	for _, mode := range o.Modes {
		if _, exists := seen[mode]; exists {
			continue
		}

		seen[mode] = struct{}{}
		unique = append(unique, mode)
	}

	o.Modes = unique
}

type mode struct {
	Width   int `json:"width"`
	Height  int `json:"height"`
	Refresh int `json:"refresh"`
}

func (m mode) String() string {
	return fmt.Sprintf(
		"%dx%d @ %.3f Hz",
		m.Width,
		m.Height,
		float64(m.Refresh)/1000,
	)
}
