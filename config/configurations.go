package config

type Database struct {
	DbName string
	User   string
	Pass   string
}

type Jwt struct {
	SecretKey string
}

type Config struct {
	Database Database
	Jwt      Jwt
}

func LoadConfig() Config {
	return Config{Database{"db_postgres", "user", "pass"}, Jwt{"wwwwwwwwwwww"}}
}
