// Code generated by go-swagger; DO NOT EDIT.

package aks

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

// NewListAKSVMSizesParams creates a new ListAKSVMSizesParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewListAKSVMSizesParams() *ListAKSVMSizesParams {
	return &ListAKSVMSizesParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewListAKSVMSizesParamsWithTimeout creates a new ListAKSVMSizesParams object
// with the ability to set a timeout on a request.
func NewListAKSVMSizesParamsWithTimeout(timeout time.Duration) *ListAKSVMSizesParams {
	return &ListAKSVMSizesParams{
		timeout: timeout,
	}
}

// NewListAKSVMSizesParamsWithContext creates a new ListAKSVMSizesParams object
// with the ability to set a context for a request.
func NewListAKSVMSizesParamsWithContext(ctx context.Context) *ListAKSVMSizesParams {
	return &ListAKSVMSizesParams{
		Context: ctx,
	}
}

// NewListAKSVMSizesParamsWithHTTPClient creates a new ListAKSVMSizesParams object
// with the ability to set a custom HTTPClient for a request.
func NewListAKSVMSizesParamsWithHTTPClient(client *http.Client) *ListAKSVMSizesParams {
	return &ListAKSVMSizesParams{
		HTTPClient: client,
	}
}

/* ListAKSVMSizesParams contains all the parameters to send to the API endpoint
   for the list a k s VM sizes operation.

   Typically these are written to a http.Request.
*/
type ListAKSVMSizesParams struct {

	// ClientID.
	ClientID *string

	// ClientSecret.
	ClientSecret *string

	// Credential.
	Credential *string

	/* Location.

	   Location - Resource location
	*/
	Location *string

	// SubscriptionID.
	SubscriptionID *string

	// TenantID.
	TenantID *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the list a k s VM sizes params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ListAKSVMSizesParams) WithDefaults() *ListAKSVMSizesParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the list a k s VM sizes params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ListAKSVMSizesParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the list a k s VM sizes params
func (o *ListAKSVMSizesParams) WithTimeout(timeout time.Duration) *ListAKSVMSizesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the list a k s VM sizes params
func (o *ListAKSVMSizesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the list a k s VM sizes params
func (o *ListAKSVMSizesParams) WithContext(ctx context.Context) *ListAKSVMSizesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the list a k s VM sizes params
func (o *ListAKSVMSizesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the list a k s VM sizes params
func (o *ListAKSVMSizesParams) WithHTTPClient(client *http.Client) *ListAKSVMSizesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the list a k s VM sizes params
func (o *ListAKSVMSizesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithClientID adds the clientID to the list a k s VM sizes params
func (o *ListAKSVMSizesParams) WithClientID(clientID *string) *ListAKSVMSizesParams {
	o.SetClientID(clientID)
	return o
}

// SetClientID adds the clientId to the list a k s VM sizes params
func (o *ListAKSVMSizesParams) SetClientID(clientID *string) {
	o.ClientID = clientID
}

// WithClientSecret adds the clientSecret to the list a k s VM sizes params
func (o *ListAKSVMSizesParams) WithClientSecret(clientSecret *string) *ListAKSVMSizesParams {
	o.SetClientSecret(clientSecret)
	return o
}

// SetClientSecret adds the clientSecret to the list a k s VM sizes params
func (o *ListAKSVMSizesParams) SetClientSecret(clientSecret *string) {
	o.ClientSecret = clientSecret
}

// WithCredential adds the credential to the list a k s VM sizes params
func (o *ListAKSVMSizesParams) WithCredential(credential *string) *ListAKSVMSizesParams {
	o.SetCredential(credential)
	return o
}

// SetCredential adds the credential to the list a k s VM sizes params
func (o *ListAKSVMSizesParams) SetCredential(credential *string) {
	o.Credential = credential
}

// WithLocation adds the location to the list a k s VM sizes params
func (o *ListAKSVMSizesParams) WithLocation(location *string) *ListAKSVMSizesParams {
	o.SetLocation(location)
	return o
}

// SetLocation adds the location to the list a k s VM sizes params
func (o *ListAKSVMSizesParams) SetLocation(location *string) {
	o.Location = location
}

// WithSubscriptionID adds the subscriptionID to the list a k s VM sizes params
func (o *ListAKSVMSizesParams) WithSubscriptionID(subscriptionID *string) *ListAKSVMSizesParams {
	o.SetSubscriptionID(subscriptionID)
	return o
}

// SetSubscriptionID adds the subscriptionId to the list a k s VM sizes params
func (o *ListAKSVMSizesParams) SetSubscriptionID(subscriptionID *string) {
	o.SubscriptionID = subscriptionID
}

// WithTenantID adds the tenantID to the list a k s VM sizes params
func (o *ListAKSVMSizesParams) WithTenantID(tenantID *string) *ListAKSVMSizesParams {
	o.SetTenantID(tenantID)
	return o
}

// SetTenantID adds the tenantId to the list a k s VM sizes params
func (o *ListAKSVMSizesParams) SetTenantID(tenantID *string) {
	o.TenantID = tenantID
}

// WriteToRequest writes these params to a swagger request
func (o *ListAKSVMSizesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.ClientID != nil {

		// header param ClientID
		if err := r.SetHeaderParam("ClientID", *o.ClientID); err != nil {
			return err
		}
	}

	if o.ClientSecret != nil {

		// header param ClientSecret
		if err := r.SetHeaderParam("ClientSecret", *o.ClientSecret); err != nil {
			return err
		}
	}

	if o.Credential != nil {

		// header param Credential
		if err := r.SetHeaderParam("Credential", *o.Credential); err != nil {
			return err
		}
	}

	if o.Location != nil {

		// header param Location
		if err := r.SetHeaderParam("Location", *o.Location); err != nil {
			return err
		}
	}

	if o.SubscriptionID != nil {

		// header param SubscriptionID
		if err := r.SetHeaderParam("SubscriptionID", *o.SubscriptionID); err != nil {
			return err
		}
	}

	if o.TenantID != nil {

		// header param TenantID
		if err := r.SetHeaderParam("TenantID", *o.TenantID); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
