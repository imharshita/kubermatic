// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// DatacenterSpecAnexia DatacenterSpecAnexia describes a anexia datacenter.
//
// swagger:model DatacenterSpecAnexia
type DatacenterSpecAnexia struct {

	// LocationID the location of the region
	LocationID string `json:"locationID,omitempty"`
}

// Validate validates this datacenter spec anexia
func (m *DatacenterSpecAnexia) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this datacenter spec anexia based on context it is used
func (m *DatacenterSpecAnexia) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DatacenterSpecAnexia) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DatacenterSpecAnexia) UnmarshalBinary(b []byte) error {
	var res DatacenterSpecAnexia
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
