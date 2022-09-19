// Code generated by go-swagger; DO NOT EDIT.

package eks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new eks API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for eks API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	ListEKSAMITypes(params *ListEKSAMITypesParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSAMITypesOK, error)

	ListEKSCapacityTypes(params *ListEKSCapacityTypesParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSCapacityTypesOK, error)

	ListEKSClusterRoles(params *ListEKSClusterRolesParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSClusterRolesOK, error)

	ListEKSInstanceTypesNoCredentials(params *ListEKSInstanceTypesNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSInstanceTypesNoCredentialsOK, error)

	ListEKSNodeRoles(params *ListEKSNodeRolesParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSNodeRolesOK, error)

	ListEKSRegions(params *ListEKSRegionsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSRegionsOK, error)

	ListEKSSecurityGroups(params *ListEKSSecurityGroupsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSSecurityGroupsOK, error)

	ListEKSSubnets(params *ListEKSSubnetsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSSubnetsOK, error)

	ListEKSSubnetsNoCredentials(params *ListEKSSubnetsNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSSubnetsNoCredentialsOK, error)

	ListEKSVPCS(params *ListEKSVPCSParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSVPCSOK, error)

	ListEKSVPCsNoCredentials(params *ListEKSVPCsNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSVPCsNoCredentialsOK, error)

	ListEKSVersions(params *ListEKSVersionsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSVersionsOK, error)

	ValidateEKSCredentials(params *ValidateEKSCredentialsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ValidateEKSCredentialsOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
ListEKSAMITypes gets the e k s a m i types for node group
*/
func (a *Client) ListEKSAMITypes(params *ListEKSAMITypesParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSAMITypesOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListEKSAMITypesParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "listEKSAMITypes",
		Method:             "GET",
		PathPattern:        "/api/v2/eks/amitypes",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListEKSAMITypesReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListEKSAMITypesOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListEKSAMITypesDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
ListEKSCapacityTypes gets the e k s capacity types for node group
*/
func (a *Client) ListEKSCapacityTypes(params *ListEKSCapacityTypesParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSCapacityTypesOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListEKSCapacityTypesParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "listEKSCapacityTypes",
		Method:             "GET",
		PathPattern:        "/api/v2/eks/capacitytypes",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListEKSCapacityTypesReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListEKSCapacityTypesOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListEKSCapacityTypesDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
ListEKSClusterRoles lists e k s cluster service roles
*/
func (a *Client) ListEKSClusterRoles(params *ListEKSClusterRolesParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSClusterRolesOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListEKSClusterRolesParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "listEKSClusterRoles",
		Method:             "GET",
		PathPattern:        "/api/v2/providers/eks/clusterroles",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListEKSClusterRolesReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListEKSClusterRolesOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListEKSClusterRolesDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
ListEKSInstanceTypesNoCredentials gets the e k s instance types for node group based on architecture
*/
func (a *Client) ListEKSInstanceTypesNoCredentials(params *ListEKSInstanceTypesNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSInstanceTypesNoCredentialsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListEKSInstanceTypesNoCredentialsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "listEKSInstanceTypesNoCredentials",
		Method:             "GET",
		PathPattern:        "/api/v2/projects/{project_id}/kubernetes/clusters/{cluster_id}/providers/eks/instancetypes",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListEKSInstanceTypesNoCredentialsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListEKSInstanceTypesNoCredentialsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListEKSInstanceTypesNoCredentialsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
ListEKSNodeRoles lists e k s node i a m roles
*/
func (a *Client) ListEKSNodeRoles(params *ListEKSNodeRolesParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSNodeRolesOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListEKSNodeRolesParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "listEKSNodeRoles",
		Method:             "GET",
		PathPattern:        "/api/v2/providers/eks/noderoles",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListEKSNodeRolesReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListEKSNodeRolesOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListEKSNodeRolesDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
ListEKSRegions lists e k s regions
*/
func (a *Client) ListEKSRegions(params *ListEKSRegionsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSRegionsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListEKSRegionsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "listEKSRegions",
		Method:             "GET",
		PathPattern:        "/api/v2/providers/eks/regions",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListEKSRegionsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListEKSRegionsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListEKSRegionsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
ListEKSSecurityGroups lists e k s securitygroup list
*/
func (a *Client) ListEKSSecurityGroups(params *ListEKSSecurityGroupsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSSecurityGroupsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListEKSSecurityGroupsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "listEKSSecurityGroups",
		Method:             "GET",
		PathPattern:        "/api/v2/providers/eks/securitygroups",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListEKSSecurityGroupsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListEKSSecurityGroupsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListEKSSecurityGroupsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
ListEKSSubnets lists e k s subnet list
*/
func (a *Client) ListEKSSubnets(params *ListEKSSubnetsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSSubnetsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListEKSSubnetsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "listEKSSubnets",
		Method:             "GET",
		PathPattern:        "/api/v2/providers/eks/subnets",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListEKSSubnetsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListEKSSubnetsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListEKSSubnetsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
ListEKSSubnetsNoCredentials gets the e k s subnets for node group
*/
func (a *Client) ListEKSSubnetsNoCredentials(params *ListEKSSubnetsNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSSubnetsNoCredentialsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListEKSSubnetsNoCredentialsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "listEKSSubnetsNoCredentials",
		Method:             "GET",
		PathPattern:        "/api/v2/projects/{project_id}/kubernetes/clusters/{cluster_id}/providers/eks/subnets",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListEKSSubnetsNoCredentialsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListEKSSubnetsNoCredentialsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListEKSSubnetsNoCredentialsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
ListEKSVPCS Lists EKS vpc's
*/
func (a *Client) ListEKSVPCS(params *ListEKSVPCSParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSVPCSOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListEKSVPCSParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "listEKSVPCS",
		Method:             "GET",
		PathPattern:        "/api/v2/providers/eks/vpcs",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListEKSVPCSReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListEKSVPCSOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListEKSVPCSDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
ListEKSVPCsNoCredentials gets the e k s vpc s for node group
*/
func (a *Client) ListEKSVPCsNoCredentials(params *ListEKSVPCsNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSVPCsNoCredentialsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListEKSVPCsNoCredentialsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "listEKSVPCsNoCredentials",
		Method:             "GET",
		PathPattern:        "/api/v2/projects/{project_id}/kubernetes/clusters/{cluster_id}/providers/eks/vpcs",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListEKSVPCsNoCredentialsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListEKSVPCsNoCredentialsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListEKSVPCsNoCredentialsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
ListEKSVersions Lists EKS versions
*/
func (a *Client) ListEKSVersions(params *ListEKSVersionsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSVersionsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListEKSVersionsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "listEKSVersions",
		Method:             "GET",
		PathPattern:        "/api/v2/providers/eks/versions",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListEKSVersionsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListEKSVersionsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListEKSVersionsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
ValidateEKSCredentials Validates EKS credentials
*/
func (a *Client) ValidateEKSCredentials(params *ValidateEKSCredentialsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ValidateEKSCredentialsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewValidateEKSCredentialsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "validateEKSCredentials",
		Method:             "GET",
		PathPattern:        "/api/v2/providers/eks/validatecredentials",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ValidateEKSCredentialsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ValidateEKSCredentialsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ValidateEKSCredentialsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
