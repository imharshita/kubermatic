// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NamespaceSpec namespace spec
//
// swagger:model NamespaceSpec
type NamespaceSpec struct {

	// Annotations of the namespace
	Annotations map[string]string `json:"annotations,omitempty"`

	// Create defines whether the namespace should be created if it does not exist.
	Create bool `json:"create,omitempty"`

	// Labels of the namespace
	Labels map[string]string `json:"labels,omitempty"`

	// Name is the namespace to deploy the Application into
	Name string `json:"name,omitempty"`
}

// Validate validates this namespace spec
func (m *NamespaceSpec) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this namespace spec based on context it is used
func (m *NamespaceSpec) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *NamespaceSpec) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *NamespaceSpec) UnmarshalBinary(b []byte) error {
	var res NamespaceSpec
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
