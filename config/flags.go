package config

import (
	"errors"
	"fmt"
	"strconv"
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
	return s.Host + ":" + strconv.Itoa(s.Port)
}

func (s *Server) Set(flagValue string) error {
	const flagParts = 2

	flagValues := strings.Split(flagValue, ":")

	if len(flagValues) != flagParts {
		return fmt.Errorf("parsing address error - %s: %w", flagValue, ErrInvalidServerAddress)
	}

	port, err := strconv.Atoi(flagValues[1])
	if err != nil {
		return fmt.Errorf("parsing port error - %s: %w", flagValue, ErrInvalidServerPort)
	}

	s.Port = port
	s.Host = flagValues[0]

	return nil
}
