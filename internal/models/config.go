package models

type SQLConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}
type SearchConfig struct {
	ApiKey string `yaml:"apiKey"`
	Query  string `yaml:"query"`
}
