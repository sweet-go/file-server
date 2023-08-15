package model

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type DeleteRule string

const (
	ManualDelete DeleteRule = "MANUAL_DELETE"
)

func ParseStringToDeleteRule(s string) (DeleteRule, error) {
	switch strings.ToUpper(s) {
	case "MANUAL_DELETE":
		return ManualDelete, nil
	default:
		return "", errors.New("invalid delete rule")
	}
}

type DeletableMedia struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	DeleteRule DeleteRule     `json:"delete_rule"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty"`
}

type DeleteMediaInput struct {
	ID string `json:"id" validate:"required"`
}

func (input *DeleteMediaInput) Validate() error {
	return validator.Struct(input)
}

type DeletableMediaUsecase interface {
	DeleteMedia(ctx context.Context, input *DeleteMediaInput) error
}

type DeletableMediaRepository interface {
	Create(ctx context.Context, input *DeletableMedia) error
	FindByID(ctx context.Context, id string) (*DeletableMedia, error)
	Delete(ctx context.Context, id string) error
}
