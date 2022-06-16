//go:build ee

/*
                  Kubermatic Enterprise Read-Only License
                         Version 1.0 ("KERO-1.0”)
                     Copyright © 2022 Kubermatic GmbH

   1.	You may only view, read and display for studying purposes the source
      code of the software licensed under this license, and, to the extent
      explicitly provided under this license, the binary code.
   2.	Any use of the software which exceeds the foregoing right, including,
      without limitation, its execution, compilation, copying, modification
      and distribution, is expressly prohibited.
   3.	THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND,
      EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
      MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
      IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
      CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
      TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
      SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

   END OF TERMS AND CONDITIONS
*/

package resourcequota

import (
	"context"
	"fmt"

	kubermaticv1 "k8c.io/kubermatic/v2/pkg/apis/kubermatic/v1"

	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func ValidateResourceQuota(ctx context.Context,
	incomingQuota *kubermaticv1.ResourceQuota,
	client ctrlruntimeclient.Client) *field.Error {
	errs := &field.Error{}
	fieldPath := field.NewPath("unique", "reqource quotas")

	currentQuotaList := &kubermaticv1.ResourceQuotaList{}
	if err := client.List(ctx, currentQuotaList); err != nil {
		return field.InternalError(fieldPath, fmt.Errorf("failed to list quotas: %w", err))
	}

	for _, currentQuota := range currentQuotaList.Items {
		currentSubject := currentQuota.Spec.Subject
		incomingSuject := incomingQuota.Spec.Subject
		if currentSubject.Name == incomingSuject.Name && currentSubject.Kind == incomingSuject.Kind {
			return field.Required(fieldPath, fmt.Sprintf("Subject Name %q and Project %q must be unique", incomingSuject.Name, incomingSuject.Kind))
		}
	}

	return errs
}
