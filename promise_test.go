package kratix

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/syntasso/kratix/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ = Describe("The Promise interface", func() {
	var promise PromiseImpl
	var promiseObject *v1alpha1.Promise

	BeforeEach(func() {
		promiseObject = &v1alpha1.Promise{
			ObjectMeta: metav1.ObjectMeta{
				Name:   "test-promise",
				Labels: map[string]string{"app": "test"},
				Annotations: map[string]string{
					"version": "v1",
				},
			},
			TypeMeta: metav1.TypeMeta{
				Kind:       "Promise",
				APIVersion: "platform.kratix.io/v1alpha1",
			},
			Spec: v1alpha1.PromiseSpec{
				RequiredPromises: []v1alpha1.RequiredPromise{
					{Name: "sub-promise", Version: "v1"},
				},
				DestinationSelectors: []v1alpha1.PromiseScheduling{
					{MatchLabels: map[string]string{"app": "test"}},
				},
			},
			Status: v1alpha1.PromiseStatus{
				Status:  "Super Ready",
				Version: "v1",
			},
		}

		uPromise, err := promiseObject.ToUnstructured()
		Expect(err).ToNot(HaveOccurred())

		promise = PromiseImpl{
			ResourceImpl: ResourceImpl{
				obj: *uPromise,
			},
			promise: promiseObject,
		}
	})

	It("can access the underlying Promise object", func() {
		Expect(promise.GetPromise()).To(Equal(promiseObject))
	})

	It("can access GetValue method", func() {
		val, err := promise.GetValue("spec.requiredPromises")
		Expect(err).ToNot(HaveOccurred())
		Expect(val).To(Equal([]any{
			map[string]any{
				"name":    "sub-promise",
				"version": "v1",
			},
		}))
	})

	It("can access GetStatus method", func() {
		status, err := promise.GetStatus()
		Expect(err).ToNot(HaveOccurred())
		Expect(status.Get("status")).To(Equal("Super Ready"))
		Expect(status.Get("version")).To(Equal("v1"))
	})

	It("can access GetName method", func() {
		Expect(promise.GetName()).To(Equal("test-promise"))
	})

	It("can access GetNamespace method", func() {
		Expect(promise.GetNamespace()).To(Equal(""))
	})

	It("can access GetGroupVersionKind method", func() {
		gvk := promise.GetGroupVersionKind()
		expectedGVK := schema.GroupVersionKind{
			Group:   "platform.kratix.io",
			Version: "v1alpha1",
			Kind:    "Promise",
		}
		Expect(gvk).To(Equal(expectedGVK))
	})

	It("can access GetLabels method", func() {
		Expect(promise.GetLabels()).To(Equal(map[string]string{"app": "test"}))
	})

	It("can access GetAnnotations method", func() {
		Expect(promise.GetAnnotations()).To(Equal(map[string]string{"version": "v1"}))
	})

	It("returns error for non-existent path in GetValue", func() {
		_, err := promise.GetValue("spec.nonExistent")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("path spec.nonExistent not found"))
	})
})
