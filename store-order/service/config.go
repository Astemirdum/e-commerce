package service

type OrderConfig struct {
	Addr string `yaml:"addr"`
	Port int    `yaml:"port"`
}

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
	Order   OrderConfig   `yaml:"order"`
	DB      ConfigDB      `yaml:"db"`
	Product ProductConfig `yaml:"product"`
}
