package config

type MinIOConn struct {
	Ep     string `yaml:"ep"`
	Secure bool   `yaml:"secure"`
	Bucket string `yaml:"bucket"`
}
