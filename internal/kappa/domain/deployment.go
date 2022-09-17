package domain

type Deployment struct {
	Services []Service `yaml:"services"`
}
