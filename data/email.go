package data

import (
	"net"
	"regexp"
	"strings"

	"github.com/hashicorp/go-hclog"
)

// Email defines a struct for emailing flow
type Email struct {
	logger hclog.Logger
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// NewEmail creates a new Email struct
func NewEmail(logger hclog.Logger) *Email {
	email := &Email{logger}

	return email
}

// IsEmailValid checks if the email provided passes the required structure
// and length test. It also checks the domain has a valid MX record.
func (email *Email) IsEmailValid(em string) bool {
	if len(em) < 3 && len(em) > 254 {
		return false
	}
	if !emailRegex.MatchString(em) {
		return false
	}
	parts := strings.Split(em, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}
	return true
}
