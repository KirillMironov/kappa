package domain

type Pod struct {
	Name        string        `yaml:"name"`
	Command     string        `yaml:"command"`
	Args        []string      `yaml:"args"`
	Environment []Environment `yaml:"env"`
	WorkingDir  string        `yaml:"workingDir"`
}
