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
			kratix.WithObjectPath("assets/input/object.yaml"),
			kratix.WithOutputDir("assets/output"),
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

			By("getting the name of the resource", func() {
				name := resource.GetName()
				Expect(name).To(Equal("my-resource"))
			})
		})
	})
})
