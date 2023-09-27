package auth

import (
	apiv3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	wranglerv3 "github.com/rancher/rancher/pkg/generated/controllers/management.cattle.io/v3"
	"github.com/rancher/rancher/pkg/types/config"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	prtbServiceAccountControllerName = "prtb-service-account-controller"
	serviceAccountAnnotation         = "management.cattle.io/serviceAccount"
)

type prtbServiceAccountController struct {
	prtbClient wranglerv3.ProjectRoleTemplateBindingClient
}

// newPRTBServiceAccountController creates a new ProjectRoleTemplateBinding account controller.
func newPRTBServiceAccountController(mgmt *config.ManagementContext) *prtbServiceAccountController {
	return &prtbServiceAccountController{
		prtbClient: mgmt.Wrangler.Mgmt.ProjectRoleTemplateBinding(),
	}
}

// sync updates a ProjectRoleTemplateBinding on the management cluster and sets the management.cattle.io/serviceAccount
// annotation = name of the ServiceAccount subject. This ensures the correct Service Account gets permissions defined
// in the PRTB.
func (c prtbServiceAccountController) sync(_ string, prtb *apiv3.ProjectRoleTemplateBinding) (runtime.Object, error) {
	if prtb == nil {
		return prtb, nil
	}
	if prtb.ServiceAccount != "" {
		if _, ok := prtb.Annotations[serviceAccountAnnotation]; ok {
			return prtb, nil
		}
		copied := prtb.DeepCopy()
		if copied.Annotations == nil {
			copied.Annotations = make(map[string]string)
		}
		copied.Annotations[serviceAccountAnnotation] = copied.ServiceAccount
		upd, err := c.prtbClient.Update(copied)
		if err != nil {
			return nil, err
		}
		return upd, nil
	}
	return prtb, nil
}
