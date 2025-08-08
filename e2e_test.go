package kratix_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	kratix "github.com/syntasso/go-sdk"
)

var _ = Describe("E2E Tests", func() {
	var sdk *kratix.KratixSDK
	var outputDir string

	BeforeEach(func() {
		var err error
		outputDir, err = os.MkdirTemp("", "kratix-e2e-test")
		Expect(err).ToNot(HaveOccurred())

		sdk = kratix.New(
			kratix.WithObjectPath("assets/input/resource.yaml"),
			kratix.WithOutputDir(outputDir),
		)
	})

	Describe("An example resource workflow", func() {
		It("may use all the functions in the SDK", func() {
			var resource kratix.ResourceAccessor
			By("reading the resource input", func() {
				var err error
				resource, err = sdk.ReadResourceInput()
				Expect(err).ToNot(HaveOccurred())
				Expect(resource).ToNot(BeNil())
			})

			By("accessing resource properties", func() {
				name := resource.GetName()
				Expect(name).To(Equal("my-resource"))
				
				namespace := resource.GetNamespace()
				Expect(namespace).To(Equal("my-namespace"))
				
				gvk := resource.GetGroupVersionKind()
				Expect(gvk.Kind).To(Equal("MyResource"))
				Expect(gvk.Version).To(Equal("v1"))
				
				labels := resource.GetLabels()
				Expect(labels["app"]).To(Equal("my-app"))
				Expect(labels["environment"]).To(Equal("test"))
			})

			By("accessing nested spec values", func() {
				replicas, err := resource.GetValue("spec.replicas")
				Expect(err).ToNot(HaveOccurred())
				Expect(replicas).To(BeNumerically("==", 3))
				
				dbSize, err := resource.GetValue("spec.dbConfig.size")
				Expect(err).ToNot(HaveOccurred())
				Expect(dbSize).To(Equal("large"))
				
				dbType, err := resource.GetValue("spec.dbConfig.type")
				Expect(err).ToNot(HaveOccurred())
				Expect(dbType).To(Equal("postgres"))
			})

			By("accessing status information", func() {
				status, err := resource.GetStatus("")
				Expect(err).ToNot(HaveOccurred())
				
				phase := status.Get("phase")
				Expect(phase).To(Equal("Running"))
			})
		})
	})
})
