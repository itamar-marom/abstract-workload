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

package controllers

import (
	"context"
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	examplesv1alpha1 "itamar.marom/abstractworkload/api/v1alpha1"
)

const (
	kindDeployment  = "Deployment"
	kindStatefulSet = "StatefulSet"
)

// AbstractWorkloadReconciler reconciles a AbstractWorkload object
type AbstractWorkloadReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=examples.itamar.marom,resources=abstractworkloads,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=examples.itamar.marom,resources=abstractworkloads/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=examples.itamar.marom,resources=abstractworkloads/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AbstractWorkload object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *AbstractWorkloadReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Add a uuid for each reconciliation
	log := log.FromContext(ctx).WithValues("reconcileID", uuid.NewUUID())

	// Add the controller logger to the context
	ctx = ctrl.LoggerInto(ctx, log)

	obj := &examplesv1alpha1.AbstractWorkload{}
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		log.Error(err, "unable to fetch AbstractWorkload")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var workloadKind string
	var workloadAPIVersion string
	labels := map[string]string{
		"name": req.NamespacedName.Name,
		"type": obj.Spec.WorkloadType.String(),
	}
	switch obj.Spec.WorkloadType.String() {
	case examplesv1alpha1.StrStateless:
		log.Info("Stateless application, creating Deployment")
		deployment := &v1.Deployment{}
		if err := r.Get(ctx, req.NamespacedName, deployment); err != nil {
			deployment = &v1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      req.NamespacedName.Name,
					Namespace: req.NamespacedName.Namespace,
					Labels:    labels,
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion: obj.APIVersion,
						Kind:       obj.Kind,
						Name:       obj.Name,
						UID:        obj.UID,
					}},
				},
				Spec: v1.DeploymentSpec{
					Replicas: obj.Spec.Replicas,
					Selector: &metav1.LabelSelector{MatchLabels: labels},
					Template: v12.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{Labels: labels},
						Spec: v12.PodSpec{
							Containers: []v12.Container{{
								Name:  req.NamespacedName.Name,
								Image: obj.Spec.ContainerImage,
							}},
						},
					},
				},
			}

			if err := r.Client.Create(ctx, deployment); err != nil {
				log.Error(err, "Creating Deployment object")
				return ctrl.Result{}, err
			}
		} else {
			deployment.Spec.Replicas = obj.Spec.Replicas
			deployment.Spec.Template.Spec.Containers = []v12.Container{{
				Name:  req.NamespacedName.Name,
				Image: obj.Spec.ContainerImage,
			}}

			if err := r.Client.Update(ctx, deployment); err != nil {
				log.Error(err, "Updating Deployment object")
				return ctrl.Result{}, err
			}
		}

		workloadKind = deployment.Kind
		workloadAPIVersion = deployment.APIVersion
	case examplesv1alpha1.StrStateful:
		log.Info("Stateful application, deploying StatefulSet")
		statefulSet := &v1.StatefulSet{}
		if err := r.Get(ctx, req.NamespacedName, statefulSet); err != nil {
			statefulSet = &v1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      req.NamespacedName.Name,
					Namespace: req.NamespacedName.Namespace,
					Labels:    labels,
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion: obj.APIVersion,
						Kind:       obj.Kind,
						Name:       obj.Name,
						UID:        obj.UID,
					}},
				},
				Spec: v1.StatefulSetSpec{
					Replicas: obj.Spec.Replicas,
					Selector: &metav1.LabelSelector{MatchLabels: labels},
					Template: v12.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{Labels: labels},
						Spec: v12.PodSpec{
							Containers: []v12.Container{{
								Name:  req.NamespacedName.Name,
								Image: obj.Spec.ContainerImage,
							}},
						},
					},
					VolumeClaimTemplates: []v12.PersistentVolumeClaim{{
						ObjectMeta: metav1.ObjectMeta{Name: req.NamespacedName.Name},
						Spec: v12.PersistentVolumeClaimSpec{
							AccessModes: []v12.PersistentVolumeAccessMode{v12.ReadWriteOnce},
							Resources: v12.ResourceRequirements{
								Requests: map[v12.ResourceName]resource.Quantity{
									v12.ResourceStorage: *resource.NewQuantity(1, resource.BinarySI),
								},
							},
							StorageClassName: nil,
						},
					}},
				},
			}

			if err := r.Client.Create(ctx, statefulSet); err != nil {
				log.Error(err, "Creating StatefulSet object")
				return ctrl.Result{}, err
			}
		} else {
			statefulSet.Spec.Replicas = obj.Spec.Replicas
			statefulSet.Spec.Template.Spec.Containers = []v12.Container{{
				Name:  req.NamespacedName.Name,
				Image: obj.Spec.ContainerImage,
			}}

			if err := r.Client.Update(ctx, statefulSet); err != nil {
				log.Error(err, "Updating StatefulSet object")
				return ctrl.Result{}, err
			}
		}

		workloadKind = statefulSet.Kind
		workloadAPIVersion = statefulSet.APIVersion
	}

	obj.Status.Workload = examplesv1alpha1.CrossNamespaceObjectReference{
		APIVersion: workloadAPIVersion,
		Kind:       workloadKind,
		Name:       req.NamespacedName.Name,
		Namespace:  req.NamespacedName.Namespace,
	}

	if err := r.Client.Status().Update(ctx, obj); err != nil {
		log.Error(err, "Updating workload status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *AbstractWorkloadReconciler) getAWForChildObject(workload client.Object) []reconcile.Request {
	requests := []reconcile.Request{{
		NamespacedName: types.NamespacedName{
			Name:      workload.GetName(),
			Namespace: workload.GetNamespace(),
		}}}
	return requests
}

// SetupWithManager sets up the controller with the Manager.
func (r *AbstractWorkloadReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplesv1alpha1.AbstractWorkload{}).
		Owns(&v1.Deployment{}).
		Owns(&v1.StatefulSet{}).
		Watches(
			&source.Kind{Type: &v1.Deployment{}},
			handler.EnqueueRequestsFromMapFunc(r.getAWForChildObject)).
		Watches(
			&source.Kind{Type: &v1.StatefulSet{}},
			handler.EnqueueRequestsFromMapFunc(r.getAWForChildObject)).
		Complete(r)
}
