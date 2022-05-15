package service

type ProductConfig struct {
	Addr string `yaml:"addr"`
	Port int    `yaml:"port"`
}

type ConfigDB struct {
	Username string `yaml:"username"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Dbname   string `yaml:"dbname"`
	Password string `yaml:"password"`
}

type Config struct {
	Product ProductConfig `yaml:"product"`
	DB      ConfigDB      `yaml:"db"`
}
