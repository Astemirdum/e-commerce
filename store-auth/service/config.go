package service

type AuthConfig struct {
	Addr      string `yaml:"addr"`
	Port      int    `yaml:"port"`
	SecretKey string `yaml:"key"`
}

type ConfigDB struct {
	Username string `yaml:"username"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Dbname   string `yaml:"dbname"`
	Password string `yaml:"password"`
}

type Config struct {
	Auth AuthConfig `yaml:"auth"`
	DB   ConfigDB   `yaml:"db"`
}
