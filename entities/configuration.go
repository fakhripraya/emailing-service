package entities

// Configuration Entity
type Configuration struct {
	EmailCredential EmailCredential
}

// EmailCredential is an entity that stores an smtp mailer credentials
type EmailCredential struct {
	Username string
	Password string
}
