// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// AKSMachineDeploymentCloudSpec a k s machine deployment cloud spec
//
// swagger:model AKSMachineDeploymentCloudSpec
type AKSMachineDeploymentCloudSpec struct {

	// name
	Name string `json:"name,omitempty"`

	// basics settings
	BasicsSettings *AgentPoolBasics `json:"basicsSettings,omitempty"`

	// configuration
	Configuration *AgentPoolConfig `json:"configuration,omitempty"`

	// optional settings
	OptionalSettings *AgentPoolOptionalSettings `json:"optionalSettings,omitempty"`
}

// Validate validates this a k s machine deployment cloud spec
func (m *AKSMachineDeploymentCloudSpec) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBasicsSettings(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateConfiguration(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOptionalSettings(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AKSMachineDeploymentCloudSpec) validateBasicsSettings(formats strfmt.Registry) error {
	if swag.IsZero(m.BasicsSettings) { // not required
		return nil
	}

	if m.BasicsSettings != nil {
		if err := m.BasicsSettings.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("basicsSettings")
			}
			return err
		}
	}

	return nil
}

func (m *AKSMachineDeploymentCloudSpec) validateConfiguration(formats strfmt.Registry) error {
	if swag.IsZero(m.Configuration) { // not required
		return nil
	}

	if m.Configuration != nil {
		if err := m.Configuration.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("configuration")
			}
			return err
		}
	}

	return nil
}

func (m *AKSMachineDeploymentCloudSpec) validateOptionalSettings(formats strfmt.Registry) error {
	if swag.IsZero(m.OptionalSettings) { // not required
		return nil
	}

	if m.OptionalSettings != nil {
		if err := m.OptionalSettings.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("optionalSettings")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this a k s machine deployment cloud spec based on the context it is used
func (m *AKSMachineDeploymentCloudSpec) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateBasicsSettings(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateConfiguration(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateOptionalSettings(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AKSMachineDeploymentCloudSpec) contextValidateBasicsSettings(ctx context.Context, formats strfmt.Registry) error {

	if m.BasicsSettings != nil {
		if err := m.BasicsSettings.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("basicsSettings")
			}
			return err
		}
	}

	return nil
}

func (m *AKSMachineDeploymentCloudSpec) contextValidateConfiguration(ctx context.Context, formats strfmt.Registry) error {

	if m.Configuration != nil {
		if err := m.Configuration.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("configuration")
			}
			return err
		}
	}

	return nil
}

func (m *AKSMachineDeploymentCloudSpec) contextValidateOptionalSettings(ctx context.Context, formats strfmt.Registry) error {

	if m.OptionalSettings != nil {
		if err := m.OptionalSettings.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("optionalSettings")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *AKSMachineDeploymentCloudSpec) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AKSMachineDeploymentCloudSpec) UnmarshalBinary(b []byte) error {
	var res AKSMachineDeploymentCloudSpec
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
