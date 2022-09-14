package domain

type Service struct {
	Name        string        `yaml:"name"`
	Command     string        `yaml:"command"`
	Environment []Environment `yaml:"env"`
	WorkingDir  string        `yaml:"working_dir"`
}
