package configuration

type Configuration struct {
	// indexed by the process config name
	Processes map[string]ProcessConfiguration `yaml:"processes"`
	LogsPath  string                          `yaml:"logs"`
}

type ProcessConfiguration struct {
	Command       string   `yaml:"command"`
	Cwd           string   `yaml:"cwd"`
	Env           []string `yaml:"env"`
	StopSignal    string   `yaml:"stop_signal"`
	RestartPolicy string   `yaml:"restart"`
}
