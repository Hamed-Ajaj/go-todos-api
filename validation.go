package main

import (
	"errors"
	"strings"
)

func (r *CreateTodoRequest) Validate() error {
	if strings.TrimSpace(r.Title) == "" {
		return errors.New("title is required")
	}

	if len(r.Title) < 3 {
		return errors.New("title must be at least 3 characters")
	}

	if len(r.Title) > 200 {
		return errors.New("title too long")
	}
	return nil
}
