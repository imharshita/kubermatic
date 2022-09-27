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

	"k8c.io/kubermatic/v2/pkg/test/e2e/utils/apiclient/models"
)

// NewCreateEKSClusterRoleParams creates a new CreateEKSClusterRoleParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewCreateEKSClusterRoleParams() *CreateEKSClusterRoleParams {
	return &CreateEKSClusterRoleParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewCreateEKSClusterRoleParamsWithTimeout creates a new CreateEKSClusterRoleParams object
// with the ability to set a timeout on a request.
func NewCreateEKSClusterRoleParamsWithTimeout(timeout time.Duration) *CreateEKSClusterRoleParams {
	return &CreateEKSClusterRoleParams{
		timeout: timeout,
	}
}

// NewCreateEKSClusterRoleParamsWithContext creates a new CreateEKSClusterRoleParams object
// with the ability to set a context for a request.
func NewCreateEKSClusterRoleParamsWithContext(ctx context.Context) *CreateEKSClusterRoleParams {
	return &CreateEKSClusterRoleParams{
		Context: ctx,
	}
}

// NewCreateEKSClusterRoleParamsWithHTTPClient creates a new CreateEKSClusterRoleParams object
// with the ability to set a custom HTTPClient for a request.
func NewCreateEKSClusterRoleParamsWithHTTPClient(client *http.Client) *CreateEKSClusterRoleParams {
	return &CreateEKSClusterRoleParams{
		HTTPClient: client,
	}
}

/*
CreateEKSClusterRoleParams contains all the parameters to send to the API endpoint

	for the create e k s cluster role operation.

	Typically these are written to a http.Request.
*/
type CreateEKSClusterRoleParams struct {

	// AccessKeyID.
	AccessKeyID *string

	// Body.
	Body *models.ClusterServiceRole

	// Credential.
	Credential *string

	// SecretAccessKey.
	SecretAccessKey *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the create e k s cluster role params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateEKSClusterRoleParams) WithDefaults() *CreateEKSClusterRoleParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the create e k s cluster role params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateEKSClusterRoleParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the create e k s cluster role params
func (o *CreateEKSClusterRoleParams) WithTimeout(timeout time.Duration) *CreateEKSClusterRoleParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create e k s cluster role params
func (o *CreateEKSClusterRoleParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create e k s cluster role params
func (o *CreateEKSClusterRoleParams) WithContext(ctx context.Context) *CreateEKSClusterRoleParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create e k s cluster role params
func (o *CreateEKSClusterRoleParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create e k s cluster role params
func (o *CreateEKSClusterRoleParams) WithHTTPClient(client *http.Client) *CreateEKSClusterRoleParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create e k s cluster role params
func (o *CreateEKSClusterRoleParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAccessKeyID adds the accessKeyID to the create e k s cluster role params
func (o *CreateEKSClusterRoleParams) WithAccessKeyID(accessKeyID *string) *CreateEKSClusterRoleParams {
	o.SetAccessKeyID(accessKeyID)
	return o
}

// SetAccessKeyID adds the accessKeyId to the create e k s cluster role params
func (o *CreateEKSClusterRoleParams) SetAccessKeyID(accessKeyID *string) {
	o.AccessKeyID = accessKeyID
}

// WithBody adds the body to the create e k s cluster role params
func (o *CreateEKSClusterRoleParams) WithBody(body *models.ClusterServiceRole) *CreateEKSClusterRoleParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the create e k s cluster role params
func (o *CreateEKSClusterRoleParams) SetBody(body *models.ClusterServiceRole) {
	o.Body = body
}

// WithCredential adds the credential to the create e k s cluster role params
func (o *CreateEKSClusterRoleParams) WithCredential(credential *string) *CreateEKSClusterRoleParams {
	o.SetCredential(credential)
	return o
}

// SetCredential adds the credential to the create e k s cluster role params
func (o *CreateEKSClusterRoleParams) SetCredential(credential *string) {
	o.Credential = credential
}

// WithSecretAccessKey adds the secretAccessKey to the create e k s cluster role params
func (o *CreateEKSClusterRoleParams) WithSecretAccessKey(secretAccessKey *string) *CreateEKSClusterRoleParams {
	o.SetSecretAccessKey(secretAccessKey)
	return o
}

// SetSecretAccessKey adds the secretAccessKey to the create e k s cluster role params
func (o *CreateEKSClusterRoleParams) SetSecretAccessKey(secretAccessKey *string) {
	o.SecretAccessKey = secretAccessKey
}

// WriteToRequest writes these params to a swagger request
func (o *CreateEKSClusterRoleParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.AccessKeyID != nil {

		// header param AccessKeyID
		if err := r.SetHeaderParam("AccessKeyID", *o.AccessKeyID); err != nil {
			return err
		}
	}
	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if o.Credential != nil {

		// header param Credential
		if err := r.SetHeaderParam("Credential", *o.Credential); err != nil {
			return err
		}
	}

	if o.SecretAccessKey != nil {

		// header param SecretAccessKey
		if err := r.SetHeaderParam("SecretAccessKey", *o.SecretAccessKey); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
