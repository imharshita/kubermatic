// Code generated by go-swagger; DO NOT EDIT.

package eks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"k8c.io/kubermatic/v2/pkg/test/e2e/utils/apiclient/models"
)

// CreateEKSClusterRoleReader is a Reader for the CreateEKSClusterRole structure.
type CreateEKSClusterRoleReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CreateEKSClusterRoleReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewCreateEKSClusterRoleOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewCreateEKSClusterRoleUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewCreateEKSClusterRoleForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewCreateEKSClusterRoleDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewCreateEKSClusterRoleOK creates a CreateEKSClusterRoleOK with default headers values
func NewCreateEKSClusterRoleOK() *CreateEKSClusterRoleOK {
	return &CreateEKSClusterRoleOK{}
}

/*
CreateEKSClusterRoleOK describes a response with status code 200, with default header values.

EKSClusterRole
*/
type CreateEKSClusterRoleOK struct {
	Payload *models.EKSClusterRole
}

// IsSuccess returns true when this create e k s cluster role o k response has a 2xx status code
func (o *CreateEKSClusterRoleOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this create e k s cluster role o k response has a 3xx status code
func (o *CreateEKSClusterRoleOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create e k s cluster role o k response has a 4xx status code
func (o *CreateEKSClusterRoleOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this create e k s cluster role o k response has a 5xx status code
func (o *CreateEKSClusterRoleOK) IsServerError() bool {
	return false
}

// IsCode returns true when this create e k s cluster role o k response a status code equal to that given
func (o *CreateEKSClusterRoleOK) IsCode(code int) bool {
	return code == 200
}

func (o *CreateEKSClusterRoleOK) Error() string {
	return fmt.Sprintf("[POST /api/v2/providers/eks/clusterroles][%d] createEKSClusterRoleOK  %+v", 200, o.Payload)
}

func (o *CreateEKSClusterRoleOK) String() string {
	return fmt.Sprintf("[POST /api/v2/providers/eks/clusterroles][%d] createEKSClusterRoleOK  %+v", 200, o.Payload)
}

func (o *CreateEKSClusterRoleOK) GetPayload() *models.EKSClusterRole {
	return o.Payload
}

func (o *CreateEKSClusterRoleOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.EKSClusterRole)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateEKSClusterRoleUnauthorized creates a CreateEKSClusterRoleUnauthorized with default headers values
func NewCreateEKSClusterRoleUnauthorized() *CreateEKSClusterRoleUnauthorized {
	return &CreateEKSClusterRoleUnauthorized{}
}

/*
CreateEKSClusterRoleUnauthorized describes a response with status code 401, with default header values.

EmptyResponse is a empty response
*/
type CreateEKSClusterRoleUnauthorized struct {
}

// IsSuccess returns true when this create e k s cluster role unauthorized response has a 2xx status code
func (o *CreateEKSClusterRoleUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create e k s cluster role unauthorized response has a 3xx status code
func (o *CreateEKSClusterRoleUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create e k s cluster role unauthorized response has a 4xx status code
func (o *CreateEKSClusterRoleUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this create e k s cluster role unauthorized response has a 5xx status code
func (o *CreateEKSClusterRoleUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this create e k s cluster role unauthorized response a status code equal to that given
func (o *CreateEKSClusterRoleUnauthorized) IsCode(code int) bool {
	return code == 401
}

func (o *CreateEKSClusterRoleUnauthorized) Error() string {
	return fmt.Sprintf("[POST /api/v2/providers/eks/clusterroles][%d] createEKSClusterRoleUnauthorized ", 401)
}

func (o *CreateEKSClusterRoleUnauthorized) String() string {
	return fmt.Sprintf("[POST /api/v2/providers/eks/clusterroles][%d] createEKSClusterRoleUnauthorized ", 401)
}

func (o *CreateEKSClusterRoleUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewCreateEKSClusterRoleForbidden creates a CreateEKSClusterRoleForbidden with default headers values
func NewCreateEKSClusterRoleForbidden() *CreateEKSClusterRoleForbidden {
	return &CreateEKSClusterRoleForbidden{}
}

/*
CreateEKSClusterRoleForbidden describes a response with status code 403, with default header values.

EmptyResponse is a empty response
*/
type CreateEKSClusterRoleForbidden struct {
}

// IsSuccess returns true when this create e k s cluster role forbidden response has a 2xx status code
func (o *CreateEKSClusterRoleForbidden) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create e k s cluster role forbidden response has a 3xx status code
func (o *CreateEKSClusterRoleForbidden) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create e k s cluster role forbidden response has a 4xx status code
func (o *CreateEKSClusterRoleForbidden) IsClientError() bool {
	return true
}

// IsServerError returns true when this create e k s cluster role forbidden response has a 5xx status code
func (o *CreateEKSClusterRoleForbidden) IsServerError() bool {
	return false
}

// IsCode returns true when this create e k s cluster role forbidden response a status code equal to that given
func (o *CreateEKSClusterRoleForbidden) IsCode(code int) bool {
	return code == 403
}

func (o *CreateEKSClusterRoleForbidden) Error() string {
	return fmt.Sprintf("[POST /api/v2/providers/eks/clusterroles][%d] createEKSClusterRoleForbidden ", 403)
}

func (o *CreateEKSClusterRoleForbidden) String() string {
	return fmt.Sprintf("[POST /api/v2/providers/eks/clusterroles][%d] createEKSClusterRoleForbidden ", 403)
}

func (o *CreateEKSClusterRoleForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewCreateEKSClusterRoleDefault creates a CreateEKSClusterRoleDefault with default headers values
func NewCreateEKSClusterRoleDefault(code int) *CreateEKSClusterRoleDefault {
	return &CreateEKSClusterRoleDefault{
		_statusCode: code,
	}
}

/*
CreateEKSClusterRoleDefault describes a response with status code -1, with default header values.

errorResponse
*/
type CreateEKSClusterRoleDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the create e k s cluster role default response
func (o *CreateEKSClusterRoleDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this create e k s cluster role default response has a 2xx status code
func (o *CreateEKSClusterRoleDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this create e k s cluster role default response has a 3xx status code
func (o *CreateEKSClusterRoleDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this create e k s cluster role default response has a 4xx status code
func (o *CreateEKSClusterRoleDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this create e k s cluster role default response has a 5xx status code
func (o *CreateEKSClusterRoleDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this create e k s cluster role default response a status code equal to that given
func (o *CreateEKSClusterRoleDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *CreateEKSClusterRoleDefault) Error() string {
	return fmt.Sprintf("[POST /api/v2/providers/eks/clusterroles][%d] createEKSClusterRole default  %+v", o._statusCode, o.Payload)
}

func (o *CreateEKSClusterRoleDefault) String() string {
	return fmt.Sprintf("[POST /api/v2/providers/eks/clusterroles][%d] createEKSClusterRole default  %+v", o._statusCode, o.Payload)
}

func (o *CreateEKSClusterRoleDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *CreateEKSClusterRoleDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
