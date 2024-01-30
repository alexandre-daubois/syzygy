package main

import "hypervigo/process"

func main() {
	processManager := process.NewProcessManager()

	processManager.RunCommand("sleep", "2")
	processManager.RunCommand("sleep", "3")
	processManager.Loop()
}
