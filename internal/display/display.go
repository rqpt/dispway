package display

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

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

func (m Mode) String() string {
	return fmt.Sprintf(
		"%dx%d @ %.3f Hz",
		m.Width,
		m.Height,
		float64(m.Refresh)/1000,
	)
}

func (m Mode) SwayString() string {
	return strings.ReplaceAll(m.String(), " ", "")
}

func GetOutputs() ([]Output, error) {
	cmd := exec.Command("swaymsg", "-t", "get_outputs", "-r")

	cmdOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var outputs []Output

	err = json.Unmarshal(cmdOutput, &outputs)
	if err != nil {
		return nil, err
	}

	return outputs, nil
}

func GetOutputNames(outputs []Output) []string {
	names := make([]string, 0, len(outputs))

	for _, output := range outputs {
		names = append(names, output.Name)
	}

	sort.Strings(names)

	return names
}

func GetUniqueModes(modes []Mode) []Mode {
	seen := make(map[Mode]struct{})
	unique := make([]Mode, 0, len(modes))

	for _, mode := range modes {
		if _, exists := seen[mode]; exists {
			continue
		}

		seen[mode] = struct{}{}
		unique = append(unique, mode)
	}

	return unique
}

func ConfigureOutput(output string, mode Mode) error {
	cmd := exec.Command(
		"swaymsg",
		"output",
		output,
		"mode",
		mode.SwayString(),
	)

	return cmd.Run()
}
