package config

type Cfg struct {
	Server Server `yaml:"Server"`
	Redis  Redis  `yaml:"Redis"`
	Mysql  Mysql  `yaml:"Mysql"`
	JWT    JWT    `yaml:"JWT"`
}
type Server struct {
	Port string `yaml:"Port"`
}

type Redis struct {
	Address  string `yaml:"Address"`
	Password string `yaml:"Password"`
	DB       int    `yaml:"Port"`
}

type Mysql struct {
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	Address  string `yaml:"Address"`
	Port     string `yaml:"Port"`
	DB       string `yaml:"DB"`
}

type JWT struct {
	Key         string `yaml:"Key"`
	ExpiresTime int64  `yaml:"ExpiresTime"`
	Issuer      string `yaml:"Issuer"`
}
