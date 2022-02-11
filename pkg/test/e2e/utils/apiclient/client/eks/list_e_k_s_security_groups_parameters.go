// Code generated by go-swagger; DO NOT EDIT.

package eks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewListEKSSecurityGroupsParams creates a new ListEKSSecurityGroupsParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewListEKSSecurityGroupsParams() *ListEKSSecurityGroupsParams {
	return &ListEKSSecurityGroupsParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewListEKSSecurityGroupsParamsWithTimeout creates a new ListEKSSecurityGroupsParams object
// with the ability to set a timeout on a request.
func NewListEKSSecurityGroupsParamsWithTimeout(timeout time.Duration) *ListEKSSecurityGroupsParams {
	return &ListEKSSecurityGroupsParams{
		timeout: timeout,
	}
}

// NewListEKSSecurityGroupsParamsWithContext creates a new ListEKSSecurityGroupsParams object
// with the ability to set a context for a request.
func NewListEKSSecurityGroupsParamsWithContext(ctx context.Context) *ListEKSSecurityGroupsParams {
	return &ListEKSSecurityGroupsParams{
		Context: ctx,
	}
}

// NewListEKSSecurityGroupsParamsWithHTTPClient creates a new ListEKSSecurityGroupsParams object
// with the ability to set a custom HTTPClient for a request.
func NewListEKSSecurityGroupsParamsWithHTTPClient(client *http.Client) *ListEKSSecurityGroupsParams {
	return &ListEKSSecurityGroupsParams{
		HTTPClient: client,
	}
}

/* ListEKSSecurityGroupsParams contains all the parameters to send to the API endpoint
   for the list e k s security groups operation.

   Typically these are written to a http.Request.
*/
type ListEKSSecurityGroupsParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the list e k s security groups params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ListEKSSecurityGroupsParams) WithDefaults() *ListEKSSecurityGroupsParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the list e k s security groups params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ListEKSSecurityGroupsParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the list e k s security groups params
func (o *ListEKSSecurityGroupsParams) WithTimeout(timeout time.Duration) *ListEKSSecurityGroupsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the list e k s security groups params
func (o *ListEKSSecurityGroupsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the list e k s security groups params
func (o *ListEKSSecurityGroupsParams) WithContext(ctx context.Context) *ListEKSSecurityGroupsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the list e k s security groups params
func (o *ListEKSSecurityGroupsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the list e k s security groups params
func (o *ListEKSSecurityGroupsParams) WithHTTPClient(client *http.Client) *ListEKSSecurityGroupsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the list e k s security groups params
func (o *ListEKSSecurityGroupsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *ListEKSSecurityGroupsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
