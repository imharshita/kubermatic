/*
Copyright 2022 The Kubermatic Kubernetes Platform contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package validation

import (
	"context"
	"errors"

	kubermaticv1 "k8c.io/kubermatic/v2/pkg/apis/kubermatic/v1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// validator for validating Resource Quota CRD.
type validator struct {
	client ctrlruntimeclient.Client
}

// NewValidator returns a new Resource Quota validator.
func NewValidator(client ctrlruntimeclient.Client) *validator {
	return &validator{
		client: client,
	}
}

var _ admission.CustomValidator = &validator{}

func (v *validator) ValidateCreate(ctx context.Context, obj runtime.Object) error {
	quota, ok := obj.(*kubermaticv1.ResourceQuota)
	if !ok {
		return errors.New("object is not a Resource Quota")
	}

	if quota != nil {
		return validateResourceQuota(ctx, quota, v.client)
	}
	return nil
}

func (v *validator) ValidateUpdate(_ context.Context, _, _ runtime.Object) error {
	return nil
}

func (v *validator) ValidateDelete(_ context.Context, _ runtime.Object) error {
	return nil
}
