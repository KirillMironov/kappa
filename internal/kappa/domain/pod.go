package domain

type Pod struct {
	Name        string        `yaml:"name"`
	Command     []string      `yaml:"command"`
	Environment []Environment `yaml:"env"`
	WorkingDir  string        `yaml:"workingDir"`
}
