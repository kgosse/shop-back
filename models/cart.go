package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Cart struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Status    string    `json:"status" db:"status"`
	Quantity  int       `json:"quantity" db:"quantity"`
	ProductID int       `json:"product_id" db:"product_id"`
	UserID    int       `json:"user_id" db:"user_id"`
}

// String is not required by pop and may be deleted
func (c Cart) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Carts is not required by pop and may be deleted
type Carts []Cart

// String is not required by pop and may be deleted
func (c Carts) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *Cart) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: c.Status, Name: "Status"},
		&validators.IntIsPresent{Field: c.Quantity, Name: "Quantity"},
		&validators.IntIsPresent{Field: c.ProductID, Name: "ProductID"},
		&validators.IntIsPresent{Field: c.UserID, Name: "UserID"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *Cart) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *Cart) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
