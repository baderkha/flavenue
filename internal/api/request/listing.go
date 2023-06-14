package request

import (
	"regexp"

	"github.com/baderkha/easy-gin/v1/easygin"
	"github.com/baderkha/flavenue/internal/api/model"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ easygin.IRequest = &AddListing{}

type BaseRequest struct {
}

func (b *BaseRequest) ValidationErrorFormat(err error) any {
	return nil
}

type AddListing struct {
	model.Listing
	BaseRequest
}

func (a *AddListing) Validate() error {
	return validation.ValidateStruct(&a,
		// Street cannot be empty, and the length must between 5 and 50
		validation.Field(&a.Street1, validation.Required, validation.Length(5, 50)),
		validation.Field(&a.Street2, validation.Required, validation.Length(5, 50)),
		validation.Field(&a.Street3, validation.Required, validation.Length(5, 50)),
		// City cannot be empty, and the length must between 5 and 50
		validation.Field(&a.City, validation.Required, validation.Length(5, 50)),
		// State cannot be empty, and must be a string consisting of two letters in upper case
		validation.Field(&a.State, validation.Required, validation.Match(regexp.MustCompile("^[A-Z]{2}$"))),
		// State cannot be empty, and must be a string consisting of five digits
		validation.Field(&a.Country, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{5}$"))),
		validation.Field(&a.PostalCode, validation.Required),
		validation.Field(&a.Description, validation.Required),
		validation.Field(&a.Label, validation.Required),
	)
}
