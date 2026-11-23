package main

import (
	"os"
	"os/exec"
	// "github.com/rqpt/picker"
)

func main() {
	swayMsg := exec.Command("swaymsg", "-t", "get_outputs")
	swayMsg.Stdin = os.Stdin
	swayMsg.Stdout = os.Stdout

	swayMsg.Run()
}
