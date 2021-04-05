package accesstoken

import "time"

// this interface provides the scheme to generate
// JWT access token
type Generator interface {
	Generate(indent string, expire time.Duration) (token string, err error)
}

// this interface provides the scheme to validate
// JWT access token
type Validator interface {
	Validate(token string) (indent string, err error)
}