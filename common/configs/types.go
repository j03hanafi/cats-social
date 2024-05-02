package configs

const (
	dbName     = "DB_NAME"
	dbPort     = "DB_PORT"
	dbHost     = "DB_HOST"
	dbUsername = "DB_USERNAME"
	dbPassword = "DB_PASSWORD"
	dbParams   = "DB_PARAMS"
	jwtSecret  = "JWT_SECRET"
	bcryptSalt = "BCRYPT_SALT"
)

type RuntimeConfig struct {
	App appCfg `mapstructure:"App"`
	API apiCfg `mapstructure:"API"`
	DB  dbCfg  `mapstructure:"DB"`
}

type appCfg struct {
	Name           string `mapstructure:"Name"`
	Host           string `mapstructure:"Host"`
	Port           int    `mapstructure:"Port"`
	Version        string `mapstructure:"Version"`
	Env            string `mapstructure:"Env"`
	ContextTimeout int    `mapsctructure:"ContextTimeout"`
}

type apiCfg struct {
	BaseURL    string `mapstructure:"BaseURL"`
	Timeout    int    `mapstructure:"Timeout"`
	DebugMode  bool   `mapstructure:"DebugMode"`
	BCryptSalt int    `mapstructure:"BCRYPT_SALT"`
	JWT        jwt    `mapstructure:"JWT"`
}

type jwt struct {
	Expire    int    `mapstructure:"Expire"`
	JWTSecret string `mapstructure:"JWT_SECRET"`
}

type dbCfg struct {
	Name     string   `mapstructure:"DB_NAME"`
	Port     int      `mapstructure:"DB_PORT"`
	Host     string   `mapstructure:"DB_HOST"`
	Username string   `mapstructure:"DB_USERNAME"`
	Password string   `mapstructure:"DB_PASSWORD"`
	Params   []string `mapstructure:"DB_PARAMS"`
}
