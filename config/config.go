package config

var (
	// API Settings
	APP_PORT = 5000
	MODE     = `DEV`
	DEBUG    = true
	// CORS Settings
	ALLOW_ORIGINS = `*`
	ALLOW_HEADERS = `*`
	// Auth Settings
	JWT_SECRET  = `Enter Your Secret`
	JWT_EXPIRES = int64(84600) // One Day
	SALT        = `SuperSALTYnotSweet`
	// DB Settings
	DB_USERNAME = `ryan`
	DB_PASSWORD = `123`
	DB_HOST     = `s.a`
	DB_DATABASE = `test`
	DB_PORT     = `3306`
	// SMTP Settings
	SMTP_ENABLED = false
	SMTP_HOST    = ``
	SMTP_USER    = ``
	SMTP_PASS    = ``
)
