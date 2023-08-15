// Package model contains struct / interface to define internal use only datatype
package model

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// DeleteRule is a type to define how media should be deleted
type DeleteRule string

// list of available delete rule
const (
	ManualDelete DeleteRule = "MANUAL_DELETE"
)

// ParseStringToDeleteRule is a function to parse string to DeleteRule
func ParseStringToDeleteRule(s string) (DeleteRule, error) {
	switch strings.ToUpper(s) {
	case "MANUAL_DELETE":
		return ManualDelete, nil
	default:
		return "", errors.New("invalid delete rule")
	}
}

// DeletableMedia is a struct to define media that can be deleted
type DeletableMedia struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	DeleteRule DeleteRule     `json:"delete_rule"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty"`
}

// DeleteMediaInput is a struct to define input to delete media
type DeleteMediaInput struct {
	ID string `json:"id" validate:"required"`
}

// Validate is a function to validate DeleteMediaInput
func (input *DeleteMediaInput) Validate() error {
	return validator.Struct(input)
}

// DeletableMediaUsecase is an interface to define usecase for deletable media
type DeletableMediaUsecase interface {
	DeleteMedia(ctx context.Context, input *DeleteMediaInput) error
}

// DeletableMediaRepository is an interface to define repository for deletable media
type DeletableMediaRepository interface {
	Create(ctx context.Context, input *DeletableMedia) error
	FindByID(ctx context.Context, id string) (*DeletableMedia, error)
	Delete(ctx context.Context, id string) error
}
