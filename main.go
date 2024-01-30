package main

import "hypervigo/process"

func main() {
	processManager := process.NewProcessManager()

	processManager.RunCommand([]string{"sleep", "2"}, process.Always)
	processManager.RunCommand([]string{"sleep", "3"}, process.Never)
	processManager.Loop()
}
