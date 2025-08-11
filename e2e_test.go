package kratix_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	kratix "github.com/syntasso/go-sdk"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ = Describe("E2E Tests", func() {
	var (
		sdk         *kratix.KratixSDK
		outputDir   string
		metadataDir string
	)

	BeforeEach(func() {
		var err error
		outputDir, err = os.MkdirTemp("", "kratix-e2e-test")
		Expect(err).ToNot(HaveOccurred())

		metadataDir, err = os.MkdirTemp("", "kratix-e2e-test")
		Expect(err).ToNot(HaveOccurred())

		sdk = kratix.New(
			kratix.WithInputDir("assets/input"),
			kratix.WithInputObject("resource.yaml"),
			kratix.WithOutputDir(outputDir),
			kratix.WithMetadataDir(metadataDir),
		)

		copyAssetsToMetadata(metadataDir)
	})

	AfterEach(func() {
		os.RemoveAll(outputDir)
		os.RemoveAll(metadataDir)
	})

	Describe("Accessing the environment variables", func() {
		BeforeEach(func() {
			os.Setenv("KRATIX_WORKFLOW_ACTION", "create")
			os.Setenv("KRATIX_WORKFLOW_TYPE", "resource")
			os.Setenv("KRATIX_PROMISE_NAME", "my-promise")
			os.Setenv("KRATIX_PIPELINE_NAME", "my-pipeline")
		})

		It("returns the correct values", func() {
			Expect(sdk.WorkflowAction()).To(Equal("create"))
			Expect(sdk.WorkflowType()).To(Equal("resource"))
			Expect(sdk.PromiseName()).To(Equal("my-promise"))
			Expect(sdk.PipelineName()).To(Equal("my-pipeline"))
		})
	})

	Describe("An example workflow", func() {
		It("may use all the functions in the SDK", func() {
			By("reading files", func() {
				By("reading the resource input", func() {
					resource, err := sdk.ReadResourceInput()
					Expect(err).ToNot(HaveOccurred())
					Expect(resource).ToNot(BeNil())

					By("accessing resource properties", func() {
						Expect(resource.GetName()).To(Equal("my-resource"))
						Expect(resource.GetNamespace()).To(Equal("my-namespace"))

						gvk := resource.GetGroupVersionKind()
						Expect(gvk.Kind).To(Equal("MyResource"))
						Expect(gvk.Version).To(Equal("v1"))

						Expect(resource.GetLabels()).To(SatisfyAll(
							HaveKeyWithValue("app", "my-app"),
							HaveKeyWithValue("environment", "test"),
						))
						Expect(resource.GetAnnotations()).To(SatisfyAll(
							HaveKeyWithValue("description", "Example resource for e2e testing"),
						))
					})

					By("accessing nested spec values", func() {
						testCases := []struct {
							path     string
							expected any
						}{
							{path: "metadata.name", expected: "my-resource"},
							{path: "spec.replicas", expected: float64(3)},
							{path: "spec.dbConfig.size", expected: "large"},
							{path: "spec.dbConfig.type", expected: "postgres"},
						}

						for _, tc := range testCases {
							value, err := resource.GetValue(tc.path)
							Expect(err).ToNot(HaveOccurred())
							Expect(value).To(Equal(tc.expected))
						}
					})

					By("accessing status information", func() {
						status, err := resource.GetStatus()
						Expect(err).ToNot(HaveOccurred())
						Expect(status.Get("phase")).To(Equal("Running"))
						Expect(status.Get("conditions")).To(Equal([]any{
							map[string]any{
								"type":               "Ready",
								"status":             "True",
								"lastTransitionTime": "2024-01-01T12:00:00Z",
							},
						}))
					})
				})

				By("reading the status file", func() {
					status, err := sdk.ReadStatus()
					Expect(err).ToNot(HaveOccurred())
					Expect(status.Get("message")).To(Equal("input status"))
					Expect(status.Get("map")).To(Equal(map[string]any{
						"key":        "value",
						"anotherKey": "another-value",
					}))
				})

				By("reading the destination selectors", func() {
					destinationSelectors, err := sdk.ReadDestinationSelectors()
					Expect(err).ToNot(HaveOccurred())
					Expect(destinationSelectors).ToNot(BeNil())
					Expect(destinationSelectors[0].Directory).To(Equal("foo/bar"))
					Expect(destinationSelectors[0].MatchLabels).To(Equal(map[string]string{"app": "my-app"}))
				})
			})

			By("writing files", func() {
				By("writing to the output directory, creating any nested directories", func() {
					err := sdk.WriteOutput("foo/bar/output.yaml", []byte("output"))
					Expect(err).ToNot(HaveOccurred())
					content := readFileContent(outputDir, "foo/bar/output.yaml")
					Expect(string(content)).To(Equal("output"))
				})

				By("writing to the destination selectors yaml", func() {
					err := sdk.WriteDestinationSelectors([]kratix.DestinationSelector{{
						Directory:   "foo/bar",
						MatchLabels: map[string]string{"app": "new-app"},
					}})
					Expect(err).ToNot(HaveOccurred())

					content := readFileContent(metadataDir, "destination_selectors.yaml")
					Expect(content).To(MatchYAML("[{directory: foo/bar, matchLabels: {app: new-app}}]"))
				})

				By("writing to the status yaml", func() {
					status := &kratix.Status{}
					status.Set("nested.field", "nested-value")
					status.Set("message", "status from metadata")

					Expect(sdk.WriteStatus(status)).To(Succeed())
					content := readFileContent(metadataDir, "status.yaml")
					Expect(content).To(MatchYAML(`{nested: {field: nested-value}, message: "status from metadata"}`))
				})
			})
		})
	})

	When("the input object is a promise", func() {
		BeforeEach(func() {
			sdk = kratix.New(
				kratix.WithInputDir("assets/input"),
				kratix.WithInputObject("promise.yaml"),
				kratix.WithOutputDir(outputDir),
				kratix.WithMetadataDir(metadataDir),
			)
		})

		It("can read the promise input", func() {
			var promise kratix.Promise

			By("reading the promise input", func() {
				var err error
				promise, err = sdk.ReadPromiseInput()
				Expect(err).ToNot(HaveOccurred())
				Expect(promise).ToNot(BeNil())
			})

			By("accessing promise properties", func() {
				Expect(promise.GetName()).To(Equal("my-promise"))
				Expect(promise.GetNamespace()).To(Equal(""))
				Expect(promise.GetGroupVersionKind()).To(Equal(schema.GroupVersionKind{
					Group:   "platform.kratix.io",
					Version: "v1alpha1",
					Kind:    "Promise",
				}))
				Expect(promise.GetLabels()).To(SatisfyAll(
					HaveKeyWithValue("kratix.io/promise-version", "v0.1.0"),
				))
				Expect(promise.GetAnnotations()).To(SatisfyAll(
					HaveKeyWithValue("some-annotation", "some-value"),
				))
			})

			By("accessing promise status", func() {
				status, err := promise.GetStatus()
				Expect(err).ToNot(HaveOccurred())
				Expect(status.Get("workflowsSucceeded")).To(Equal(1))
			})

			By("accessing the promise object", func() {
				kratixPromise := promise.GetPromise()
				Expect(kratixPromise.Spec.RequiredPromises).To(HaveLen(1))
				Expect(kratixPromise.Spec.Workflows.Resource.Configure).To(HaveLen(1))
			})
		})
	})
})

// readFileContent is a helper function to read file content and verify no error occurred
func readFileContent(baseDir, relativePath string) []byte {
	GinkgoHelper()
	fullpath := filepath.Join(baseDir, relativePath)
	content, err := os.ReadFile(fullpath)
	Expect(err).ToNot(HaveOccurred())
	return content
}

func copyAssetsToMetadata(metadataDir string) {
	GinkgoHelper()
	for _, file := range []string{"destination_selectors.yaml", "status.yaml"} {
		input, err := os.ReadFile(filepath.Join("assets/metadata", file))
		Expect(err).ToNot(HaveOccurred())
		err = os.WriteFile(filepath.Join(metadataDir, file), input, 0644)
		Expect(err).ToNot(HaveOccurred())
	}
}
