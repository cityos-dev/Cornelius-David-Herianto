package uuid

import (
	"github.com/google/uuid"
)

// Utils to generate a new uuid
type Utils interface {
	Generate() (string, error)
	IsValidUUID(u string) bool
}

// utils returns struct
type utils struct{}

// NewUtils returns new utils instance
func NewUtils() Utils {
	return &utils{}
}

// Generate returns back the uuid in string or error
func (s *utils) Generate() (string, error) {
	tempUUID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return tempUUID.String(), nil
}

// IsValidUUID validate the provided uuid
func (s *utils) IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
