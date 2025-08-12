package kratix

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ = Describe("ResourceImpl", func() {
	var resource ResourceImpl
	var labels map[string]string
	var annotations map[string]string

	BeforeEach(func() {
		resourceObject := unstructured.Unstructured{}
		spec := map[string]any{
			"spec": map[string]any{
				"dbConfig": map[string]any{
					"size": "small",
				},
			},
		}
		resourceObject.SetUnstructuredContent(spec)
		resourceObject.SetName("my-resource")
		resourceObject.SetNamespace("default")
		resourceObject.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   "mygroup.example",
			Version: "v1",
			Kind:    "mykind",
		})

		annotations = map[string]string{
			"app.kubernetes.io/managed-by": "kratix",
			"app.kubernetes.io/name":       "my-resource",
		}
		resourceObject.SetAnnotations(annotations)

		labels = map[string]string{
			"sdk.io/type":    "resource",
			"sdk.io/promise": "mykind",
		}
		resourceObject.SetLabels(labels)

		resource = ResourceImpl{
			obj: resourceObject,
		}
	})

	Describe("GetName", func() {
		It("returns the name of the underlying object", func() {
			Expect(resource.GetName()).To(Equal("my-resource"))
		})
	})

	Describe("GetNamespace", func() {
		It("returns the namespace of the underlying object", func() {
			Expect(resource.GetNamespace()).To(Equal("default"))
		})
	})

	Describe("GetGroupVersionKind", func() {
		It("returns the GroupVersionKind of the underlying object", func() {
			expectedGVK := schema.GroupVersionKind{
				Group:   "mygroup.example",
				Version: "v1",
				Kind:    "mykind",
			}
			Expect(resource.GetGroupVersionKind()).To(Equal(expectedGVK))
		})
	})

	Describe("GetLabels", func() {
		It("returns the labels set on the underlying object", func() {
			Expect(resource.GetLabels()).To(Equal(labels))
		})
	})

	Describe("GetAnnotations", func() {
		It("returns the labels set on the underlying object", func() {
			Expect(resource.GetAnnotations()).To(Equal(annotations))
		})
	})

	Describe("GetValue", func() {
		When("the value exists on the resource", func() {
			It("returns the value of the provided keys", func() {
				Expect(resource.GetValue("spec.dbConfig.size")).To(Equal("small"))
				Expect(resource.GetValue("spec.dbConfig")).To(Equal(map[string]any{
					"size": "small",
				}))
			})

			It("can handle dot-prefixed keys", func() {
				Expect(resource.GetValue(".spec.dbConfig.size")).To(Equal("small"))
			})
		})
	})

	Describe("ToUnstructured", func() {
		It("returns the underlying unstructured object", func() {
			Expect(resource.ToUnstructured()).To(Equal(resource.obj))
		})
	})
})
