package interactive

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"szg/configuration"
	"szg/process"
)

func NewCommandHandler(processManager *process.ProcessManager) *InteractiveCommandHandler {
	return &InteractiveCommandHandler{
		ProcessManager: processManager,
	}
}

type InteractiveCommandHandler struct {
	ProcessManager *process.ProcessManager
}

func parseCommand(command string) (string, string) {
	cmd := strings.Split(command, " ")

	return cmd[0], strings.Join(cmd[1:], " ")
}

func (ich *InteractiveCommandHandler) HandleCommand(command string) {
	command, args := parseCommand(command)

	switch command {
	case "start":
		reader := configuration.NewReader()

		activeConf, err := reader.Read(args)
		if err != nil {
			log.Printf("Cannot read configuration file '%s' because of %s\n", args, err)

			return
		}

		fmt.Printf("Starting processes from configuration file '%s'\n", args)

		ich.ProcessManager.SetLogsPath(activeConf.LogsPath)
		go ich.ProcessManager.RunFromConfiguration(activeConf)
	case "stop":
		fmt.Println("Stopping process")
		err := ich.ProcessManager.Stop(args)

		if err != nil {
			fmt.Printf("Cannot stop process '%s' because of %s\n", args, err)
		}
	case "list":
		fmt.Println("Listing running processes")
		for pid, p := range ich.ProcessManager.Processes {
			fmt.Printf("Process '%s' (%d): %s (status: %s)\n", p.Name, pid, p.Command, p.Status)
		}
	case "exit":
		fmt.Println("Exiting...")
		ich.ProcessManager.StopAll()
		os.Exit(0)
	case "help":
		fmt.Println("Available commands:")
		fmt.Println("start <configuration file>")
		fmt.Println("list")
		fmt.Println("exit")
		fmt.Println("help")
	default:
		fmt.Println("Unknown command")
	}
}

func (ich *InteractiveCommandHandler) Loop() {
	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("szg> ")
		command, err := r.ReadString('\n')
		if err != nil {
			continue
		}

		ich.HandleCommand(strings.TrimSpace(command))
	}
}
