package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/thanhpk/randstr"
)

// ShortenedLink is used by pop to map your shortened_links database table to your go code.
type ShortenedLink struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	URL       string    `json:"url" db:"url"`
	ShortCode string    `json:"short_code" db:"short_code"`
	Hits      int       `json:"hits" db:"hits"`
}

// String is not required by pop and may be deleted
func (s ShortenedLink) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// ShortenedLinks is not required by pop and may be deleted
type ShortenedLinks []ShortenedLink

// String is not required by pop and may be deleted
func (s ShortenedLinks) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *ShortenedLink) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.URLIsPresent{Field: s.URL, Name: "URL"},
		&validators.StringIsPresent{Field: s.ShortCode, Name: "ShortCode"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *ShortenedLink) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *ShortenedLink) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	s.Validate(tx)

	return validate.NewErrors(), nil
}

func (s *ShortenedLink) BeforeValidate(tx *pop.Connection) error {
	if s.ShortCode == "" {
		s.ShortCode = randstr.String(6)
	}
	return nil
}
