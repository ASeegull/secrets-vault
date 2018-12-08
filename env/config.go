package env

type Config struct {
	DB_HOST     string `envconfig:"DB_HOST" default:"localhost:5432"`
	DB_USER     string `envconfig:"DB_USER" default:"postgres"`
	DB_PASSWORD string `envconfig:"DB_PASSWORD" required:"true"`
	DB_NAME     string `envconfig:"DB_NAME" required:"true"`
	DB_SSLMODE  string `envconfig:"DB_SSLMODE" default:"disable"`

	HOST        string `envconfig:"HOST" default:"localhost"`
	PORT        string `envconfig:"PORT" default:"5000"`
	SECRET_SALT string `envconfig:"SECRET_SALT" default:"secret"`
}
