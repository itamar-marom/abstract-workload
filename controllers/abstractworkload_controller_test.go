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
	var AbstractWorkloadReplica int32 = 2
	const (
		AbstractWorkloadName           = "test-stateless"
		AbstractWorkloadNamespace      = "default"
		AbstractWorkloadContainerImage = "nginx:latest"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Second * 5
	)

	Context("Lifecycle of stateless AbstractWorkload", func() {
		It("Should manage Deployment object and update the AbstractWorkload status", func() {
			By("By creating a new AbstractWorkload")
			ctx := context.Background()
			abstractWorkload := &v1alpha1.AbstractWorkload{
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
			Expect(k8sClient.Create(ctx, abstractWorkload)).Should(Succeed())

			By("Checking the created Deployment")
			deploymentLookupKey := types.NamespacedName{Name: AbstractWorkloadName, Namespace: AbstractWorkloadNamespace}
			createdDeployment := &v12.Deployment{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, deploymentLookupKey, createdDeployment)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			//Expect(&createdDeployment.Spec.Replicas).Should(Equal(&AbstractWorkloadReplica))
			Expect(createdDeployment.Spec.Template.Spec.Containers[0].Image).Should(Equal(AbstractWorkloadContainerImage))

			By("Checking the AbstractWorkload's status")
			abstractWorkloadLookupKey := types.NamespacedName{Namespace: AbstractWorkloadNamespace, Name: AbstractWorkloadName}
			updatedAbstractWorkload := &v1alpha1.AbstractWorkload{}
			Consistently(func() (string, error) {
				err := k8sClient.Get(ctx, abstractWorkloadLookupKey, updatedAbstractWorkload)
				if err != nil {
					return "", err
				}
				return updatedAbstractWorkload.Status.Workload.Kind, nil
			}, duration, interval).Should(Equal(createdDeployment.Kind))
			//Expect(createdAbstractWorkload.Status.Workload.Kind).Should(Equal(createdDeployment.Kind))

			//By("Delete the Deployment and see that it being created again")
			//err := k8sClient.Delete(ctx, createdDeployment)
			//Expect(err).ToNot(HaveOccurred(), "failed to delete Deployment")
			//newCreatedDeployment := &v12.Deployment{}
			//Eventually(func() bool {
			//	err := k8sClient.Get(ctx, deploymentLookupKey, newCreatedDeployment)
			//	if err != nil {
			//		return false
			//	}
			//	return true
			//}, timeout, interval).Should(BeTrue())
			//Expect(createdDeployment.Spec.Replicas).Should(Equal(AbstractWorkloadReplica))
			//Expect(createdDeployment.Spec.Template.Spec.Containers[0].Image).Should(Equal(AbstractWorkloadContainerImage))
			//
			//By("Changing replicas number of AbstractWorkload and check Deployment")
			//*createdAbstractWorkload.Spec.Replicas = 1
			//err = k8sClient.Update(ctx, createdAbstractWorkload)
			//Expect(err).ToNot(HaveOccurred(), "failed to update AbstractWorkload")
			//Consistently(func() (*int32, error) {
			//	err := k8sClient.Get(ctx, deploymentLookupKey, newCreatedDeployment)
			//	if err != nil {
			//		var badReplicas *int32
			//		*badReplicas = -1
			//		return badReplicas, err
			//	}
			//	return newCreatedDeployment.Spec.Replicas, nil
			//}, duration, interval).Should(Equal(*createdAbstractWorkload.Spec.Replicas))
			//
			//By("Deleting the AbstractWorkload and check Deployment is deleted")
			//err = k8sClient.Delete(ctx, createdAbstractWorkload)
			//Expect(err).ToNot(HaveOccurred(), "failed to delete AbstractWorkload")
			//Eventually(func() bool {
			//	err := k8sClient.Get(ctx, deploymentLookupKey, newCreatedDeployment)
			//	if err != nil {
			//		return true
			//	}
			//	return false
			//}, timeout, interval).Should(BeTrue())
		})
	})
})
