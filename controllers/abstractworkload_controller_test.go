package controllers

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"itamar.marom/abstractworkload/api/v1alpha1"
	v12 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"time"
)

var _ = Describe("CronJob controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		AbstractWorkloadName           = "test-stateless"
		AbstractWorkloadNamespace      = "default"
		AbstractWorkloadContainerImage = "nginx:latest"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Second * 5
	)
	var (
		deploymentLookupKey             = types.NamespacedName{Name: AbstractWorkloadName, Namespace: AbstractWorkloadNamespace}
		abstractWorkloadLookupKey       = types.NamespacedName{Namespace: AbstractWorkloadNamespace, Name: AbstractWorkloadName}
		AbstractWorkloadReplica   int32 = 2
		abstractWorkload                = &v1alpha1.AbstractWorkload{
			ObjectMeta: metav1.ObjectMeta{
				Name:      AbstractWorkloadName,
				Namespace: AbstractWorkloadNamespace,
			},
			Spec: v1alpha1.AbstractWorkloadSpec{
				Replicas:       &AbstractWorkloadReplica,
				ContainerImage: AbstractWorkloadContainerImage,
				WorkloadType:   v1alpha1.StrStateless,
			},
		}
	)

	Context("Lifecycle of stateless AbstractWorkload", func() {
		It("Should manage Deployment object and update the AbstractWorkload status", func() {
			By("By creating a new AbstractWorkload")
			ctx := context.Background()
			Expect(k8sClient.Create(ctx, abstractWorkload)).Should(Succeed())

			By("Checking the created Deployment")
			createdDeployment := &v12.Deployment{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, deploymentLookupKey, createdDeployment)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(createdDeployment.Spec.Replicas).Should(Equal(abstractWorkload.Spec.Replicas))
			Expect(createdDeployment.Spec.Template.Spec.Containers[0].Image).Should(Equal(AbstractWorkloadContainerImage))

			//By("Checking the AbstractWorkload's status")
			//statusAbstractWorkload := &v1alpha1.AbstractWorkload{}
			//Eventually(func() string {
			//	if err := k8sClient.Get(ctx, abstractWorkloadLookupKey, statusAbstractWorkload); err != nil {
			//		return ""
			//	}
			//	return statusAbstractWorkload.Status.Workload.Kind
			//}, duration, interval).Should(Equal(createdDeployment.Kind))

			By("Delete the Deployment and see that it being created again")
			Expect(k8sClient.Delete(ctx, createdDeployment)).Should(Succeed())
			newCreatedDeployment := &v12.Deployment{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, deploymentLookupKey, newCreatedDeployment)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(newCreatedDeployment.Spec.Replicas).Should(Equal(abstractWorkload.Spec.Replicas))
			Expect(newCreatedDeployment.Spec.Template.Spec.Containers[0].Image).Should(Equal(AbstractWorkloadContainerImage))

			By("Changing replicas number of AbstractWorkload and check Deployment")
			replicaAbstractWorkload := &v1alpha1.AbstractWorkload{}
			replicaDeployment := &v12.Deployment{}

			Expect(k8sClient.Get(ctx, abstractWorkloadLookupKey, replicaAbstractWorkload)).Should(Succeed())
			*replicaAbstractWorkload.Spec.Replicas = 1
			Expect(k8sClient.Update(ctx, replicaAbstractWorkload)).Should(Succeed())

			Eventually(func() (*int32, error) {
				err := k8sClient.Get(ctx, deploymentLookupKey, replicaDeployment)
				if err != nil {
					var badReplicas *int32
					*badReplicas = -1
					return badReplicas, err
				}
				return replicaDeployment.Spec.Replicas, nil
			}, duration, interval).Should(Equal(replicaAbstractWorkload.Spec.Replicas))

			//By("Deleting the AbstractWorkload and check Deployment is deleted")
			//deleteAbstractWorkload := &v1alpha1.AbstractWorkload{}
			//deleteDeployment := &v12.Deployment{}
			//
			//Expect(k8sClient.Get(ctx, abstractWorkloadLookupKey, deleteAbstractWorkload)).Should(Succeed())
			//Expect(k8sClient.Delete(ctx, deleteAbstractWorkload)).Should(Succeed())
			//Eventually(func() bool {
			//	err := k8sClient.Get(ctx, abstractWorkloadLookupKey, deleteAbstractWorkload)
			//	if err != nil {
			//		return true
			//	}
			//	return false
			//}, timeout, interval).Should(BeTrue())
			//Eventually(func() bool {
			//	err := k8sClient.Get(ctx, deploymentLookupKey, deleteDeployment)
			//	if err != nil {
			//		return true
			//	}
			//	return false
			//}, timeout, interval).Should(BeTrue())
		})
	})
})
