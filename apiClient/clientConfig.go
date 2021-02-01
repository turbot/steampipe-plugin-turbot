package apiClient

type ClientConfig struct {
	Credentials     ClientCredentials
	CredentialsPath string
	Profile         string
}

type ClientCredentials struct {
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
	Workspace string
}
