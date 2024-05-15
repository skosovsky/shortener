package config

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidServerAddress = errors.New("invalid server address")
	ErrInvalidServerPort    = errors.New("invalid server port")
)

type Value interface {
	String() string
	Set(string) error
}

func (s *Server) String() string {
	return s.Address
}

func (s *Server) Set(flagValue string) error {
	const flagParts = 2

	flagValues := strings.Split(flagValue, ":")

	if len(flagValues) != flagParts {
		return fmt.Errorf("parsing address error - %s: %w", flagValue, ErrInvalidServerAddress)
	}

	// for separate store host and port
	// port, err := strconv.Atoi(flagValues[1])
	// if err != nil {
	//	return fmt.Errorf("parsing port error - %s: %w", flagValue, ErrInvalidServerPort)
	// }
	//
	// s.Port = port
	// s.Host = flagValues[0]

	s.Address = flagValues[0] + ":" + flagValues[1]

	return nil
}
