package kratix

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/syntasso/kratix/api/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ = Describe("PromiseImpl", func() {
	var promise PromiseImpl
	var promiseObject *v1alpha1.Promise

	BeforeEach(func() {
		promiseObject = &v1alpha1.Promise{}
		promiseObject.SetName("test-promise")

		promise = PromiseImpl{
			ResourceImpl: ResourceImpl{
				obj: unstructured.Unstructured{},
			},
			promise: promiseObject,
		}
	})

	Describe("GetPromise", func() {
		It("returns the underlying Promise object", func() {
			Expect(promise.GetPromise()).To(Equal(promiseObject))
		})
	})
})