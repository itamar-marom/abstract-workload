/*
Copyright 2022.

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

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var abstractworkloadlog = logf.Log.WithName("abstractworkload-resource")

var defaultReplicas int32 = 1

func (r *AbstractWorkload) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-examples-itamar-marom-v1alpha1-abstractworkload,mutating=true,failurePolicy=fail,sideEffects=None,groups=examples.itamar.marom,resources=abstractworkloads,verbs=create;update,versions=v1alpha1,name=mabstractworkload.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &AbstractWorkload{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *AbstractWorkload) Default() {
	abstractworkloadlog.Info("default", "name", r.Name)

	if r.Spec.Replicas == nil {
		abstractworkloadlog.Info("setting default replicas", "name", r.Name)
		var ptrDefaultReplicas = &defaultReplicas
		r.Spec.Replicas = ptrDefaultReplicas
	}

	if r.Spec.WorkloadType == "" {
		abstractworkloadlog.Info("setting default workloadType", "name", r.Name)
		r.Spec.WorkloadType = StrStateless
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-examples-itamar-marom-v1alpha1-abstractworkload,mutating=false,failurePolicy=fail,sideEffects=None,groups=examples.itamar.marom,resources=abstractworkloads,verbs=create;update,versions=v1alpha1,name=vabstractworkload.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &AbstractWorkload{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *AbstractWorkload) ValidateCreate() error {
	abstractworkloadlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *AbstractWorkload) ValidateUpdate(old runtime.Object) error {
	abstractworkloadlog.Info("validate update", "name", r.Name)

	aw := old.(*AbstractWorkload)
	if r.Spec.WorkloadType != aw.Spec.WorkloadType {
		return field.Invalid(field.NewPath("spec").Child("workloadType"), r.Spec.WorkloadType, "workloadType is immutable")
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *AbstractWorkload) ValidateDelete() error {
	abstractworkloadlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
