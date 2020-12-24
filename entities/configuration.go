package entities

// Configuration Entity
type Configuration struct {
	AppConfig       AppConfiguration
	EmailCredential EmailCredential
}

// AppConfiguration is an entity that stores the app main configuration
type AppConfiguration struct {
	Host string
	Port string
}

// EmailCredential is an entity that stores an smtp mailer credentials
type EmailCredential struct {
	Username string
	Password string
}
