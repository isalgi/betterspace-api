package request

import (
	officeimages "backend/businesses/office_images"

	"github.com/go-playground/validator/v10"
)

type OfficeImage struct {
	URL      string `json:"url" form:"url" validate:"required"`
	OfficeID uint `json:"office_id" form:"office_id" validate:"required"`
}

func (req *OfficeImage) ToDomain() *officeimages.Domain {
	return &officeimages.Domain{
		URL: req.URL,
		OfficeID: req.OfficeID,
	}
}

func (req *OfficeImage) Validate() error {
	validate := validator.New()

	err := validate.Struct(req)

	return err
}