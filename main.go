package main

import "hypervigo/process"

func main() {
	processManager := process.NewProcessManager()

	processManager.RunCommand("sleep", "5")
	processManager.Loop()
}
