// Package model contains struct / interface to define internal use only datatype
package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// DeleteRule is a type to define how media should be deleted
type DeleteRule string

// list of available delete rule
const (
	ManualDelete    DeleteRule = "MANUAL_DELETE"
	ScheduledDelete DeleteRule = "SCHEDULED_DELETE"
)

// ParseStringToDeleteRule is a function to parse string to DeleteRule
func ParseStringToDeleteRule(s string) (DeleteRule, error) {
	switch strings.ToUpper(s) {
	case "MANUAL_DELETE":
		return ManualDelete, nil
	case "SCHEDULED_DELETE":
		return ScheduledDelete, nil
	default:
		return "", errors.New("invalid delete rule")
	}
}

// DeletableRuleMetadata is a struct to define metadata for deletable rule.
// You can add more metadata here to be used based on its delete rule.
// Any update made here, must be backward compatible, or at least doesn't break the existing data.
// Implements gorm's SerializerInterface. Can be used directly in gorm's model, and can be represented as regular text / JSONB on database.
type DeletableRuleMetadata struct {
	DeleteAfter *time.Time `json:"delete_after,omitempty"`
}

// Scan is a function to scan database value to DeletableRuleMetadata
func (drm *DeletableRuleMetadata) Scan(_ context.Context, _ *schema.Field, _ reflect.Value, dbValue interface{}) (err error) {
	if dbValue == nil {
		return
	}

	var bytes []byte
	switch v := dbValue.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSONB value: %#v", dbValue)
	}

	if err = json.Unmarshal(bytes, drm); err != nil {
		return
	}

	return
}

// Value is a function to convert DeletableRuleMetadata to json
func (drm DeletableRuleMetadata) Value(_ context.Context, _ *schema.Field, _ reflect.Value, fieldValue interface{}) (interface{}, error) {
	return json.Marshal(fieldValue)
}

// DeletableMedia is a struct to define media that can be deleted
type DeletableMedia struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	DeleteRule DeleteRule             `json:"delete_rule"`
	CreatedAt  time.Time              `json:"created_at"`
	DeletedAt  gorm.DeletedAt         `json:"deleted_at,omitempty"`
	Metadata   *DeletableRuleMetadata `json:"metadata,omitempty"`
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
	FindForScheduledDelete(ctx context.Context, limit int) ([]DeletableMedia, error)
}
