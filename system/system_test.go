package system_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/syntasso/kratix/test/kubeutils"
)

var _ = Describe("System", func() {
	var platform *kubeutils.Cluster
	var worker *kubeutils.Cluster

	BeforeEach(func() {
		platform = &kubeutils.Cluster{
			Name:    "platform",
			Context: "kind-platform",
		}

		worker = &kubeutils.Cluster{
			Name:    "worker",
			Context: "kind-worker",
		}

		kubeutils.SetTimeoutAndInterval(3*time.Minute, 1*time.Second)
		SetDefaultEventuallyTimeout(3 * time.Minute)
		SetDefaultEventuallyPollingInterval(1 * time.Second)
	})

	AfterEach(func() {
		platform.Kubectl("delete", "-f", "assets/promise.yaml")
	})

	It("should produce the expected output", func() {
		By("applying the promise.yaml", func() {
			platform.Kubectl("apply", "-f", "assets/promise.yaml")

			By("generating the expected config map on the destination", func() {
				Eventually(func() string {
					return worker.Kubectl("get", "configmap", "-n", "default")
				}).Should(ContainSubstring("config-config"))

				data := worker.Kubectl("get", "configmap", "config-config", "-n", "default", "-o", "yaml")
				Expect(data).To(ContainSubstring("workflowAction: configure"))
				Expect(data).To(ContainSubstring("workflowType: promise"))
				Expect(data).To(ContainSubstring("promiseName: config"))
				Expect(data).To(ContainSubstring("pipelineName: promise"))
			})
		})

		By("applying the resource.yaml", func() {
			platform.Kubectl("apply", "-f", "assets/example-resource.yaml")

			By("generating the expected config map on the destination", func() {
				Eventually(func() string {
					return worker.Kubectl("get", "configmap", "-n", "default")
				}).Should(ContainSubstring("go-sdk-resource-config"))

				data := worker.Kubectl("get", "configmap", "go-sdk-resource-config", "-n", "default", "-o", "yaml")
				Expect(data).To(ContainSubstring("workflowAction: configure"))
				Expect(data).To(ContainSubstring("workflowType: resource"))
				Expect(data).To(ContainSubstring("promiseName: config"))
				Expect(data).To(ContainSubstring("pipelineName: instance"))
				Expect(data).To(ContainSubstring("field0: value0"))
				Expect(data).To(ContainSubstring("field1: value1"))
			})

			By("producing the expected status", func() {
				resourceStatus := platform.Kubectl("get", "configs", "go-sdk-resource", "-n", "default", "-o", "yaml")
				Expect(resourceStatus).To(ContainSubstring("publishedDirectly: true"))
				Expect(resourceStatus).To(ContainSubstring("viaFile: true"))
			})
		})
	})
})
