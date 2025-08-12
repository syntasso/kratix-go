package kratix

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Status", func() {
	var status *StatusImpl

	BeforeEach(func() {
		status = &StatusImpl{
			data: map[string]any{
				"phase": "Ready",
				"conditions": map[string]any{
					"available": map[string]any{
						"status": "True",
						"reason": "ReconciliationSucceeded",
					},
					"degraded": map[string]any{
						"status": "False",
						"reason": "AsExpected",
					},
				},
				"observedGeneration": int64(1),
				"replicas":           int64(3),
				"pods": []any{
					map[string]any{
						"name":   "pod-1",
						"status": "Running",
						"containers": []any{
							map[string]any{
								"name":  "app",
								"ready": true,
							},
							map[string]any{
								"name":  "sidecar",
								"ready": false,
							},
						},
					},
					map[string]any{
						"name":   "pod-2",
						"status": "Pending",
						"containers": []any{
							map[string]any{
								"name":  "app",
								"ready": false,
							},
						},
					},
				},
				"events": []any{
					"Started",
					"Running",
					"Completed",
				},
			},
		}
	})

	Describe("Get", func() {
		It("retrieves top-level values", func() {
			Expect(status.Get("phase")).To(Equal("Ready"))
			Expect(status.Get("observedGeneration")).To(Equal(1))
			Expect(status.Get("replicas")).To(Equal(3))
		})

		It("retrieves nested values", func() {
			Expect(status.Get("conditions.available.status")).To(Equal("True"))
			Expect(status.Get("conditions.available.reason")).To(Equal("ReconciliationSucceeded"))
			Expect(status.Get("conditions.degraded.status")).To(Equal("False"))
		})

		It("retrieves array values using index notation", func() {
			Expect(status.Get("pods[0].name")).To(Equal("pod-1"))
			Expect(status.Get("pods[0].status")).To(Equal("Running"))
			Expect(status.Get("pods[1].name")).To(Equal("pod-2"))
			Expect(status.Get("pods[1].status")).To(Equal("Pending"))
		})

		It("retrieves deeply nested array values", func() {
			Expect(status.Get("pods[0].containers[0].name")).To(Equal("app"))
			Expect(status.Get("pods[0].containers[0].ready")).To(Equal(true))
			Expect(status.Get("pods[0].containers[1].name")).To(Equal("sidecar"))
			Expect(status.Get("pods[0].containers[1].ready")).To(Equal(false))
		})

		It("retrieves simple array values", func() {
			Expect(status.Get("events[0]")).To(Equal("Started"))
			Expect(status.Get("events[1]")).To(Equal("Running"))
			Expect(status.Get("events[2]")).To(Equal("Completed"))
		})

		It("returns nil for non-existent paths", func() {
			Expect(status.Get("nonExistent")).To(BeNil())
			Expect(status.Get("conditions.nonExistent")).To(BeNil())
			Expect(status.Get("conditions.available.nonExistent")).To(BeNil())
		})

		It("returns nil for invalid paths", func() {
			Expect(status.Get("conditions.available.status.invalid")).To(BeNil())
		})

		It("returns nil for array index out of bounds", func() {
			Expect(status.Get("pods[5]")).To(BeNil())
			Expect(status.Get("pods[0].containers[10]")).To(BeNil())
			Expect(status.Get("events[10]")).To(BeNil())
		})

		It("returns nil when accessing array with non-index", func() {
			Expect(status.Get("pods.name")).To(BeNil())
			Expect(status.Get("events.status")).To(BeNil())
		})

		It("handles empty path", func() {
			Expect(status.Get("")).To(BeNil())
		})
	})

	Describe("Set", func() {
		It("sets top-level values", func() {
			err := status.Set("phase", "Pending")
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Get("phase")).To(Equal("Pending"))
		})

		It("sets nested values", func() {
			err := status.Set("conditions.available.status", "False")
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Get("conditions.available.status")).To(Equal("False"))
		})

		It("creates new nested paths", func() {
			err := status.Set("newField.nested.deep", "value")
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Get("newField.nested.deep")).To(Equal("value"))
		})

		It("overwrites existing nested values", func() {
			err := status.Set("conditions.available.reason", "NewReason")
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Get("conditions.available.reason")).To(Equal("NewReason"))
		})

		It("sets array values using index notation", func() {
			err := status.Set("pods[0].name", "updated-pod-1")
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Get("pods[0].name")).To(Equal("updated-pod-1"))
		})

		It("sets deeply nested array values", func() {
			err := status.Set("pods[0].containers[1].ready", true)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Get("pods[0].containers[1].ready")).To(Equal(true))
		})

		It("sets simple array values", func() {
			err := status.Set("events[1]", "Updated")
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Get("events[1]")).To(Equal("Updated"))
		})

		It("extends arrays when setting beyond current length", func() {
			err := status.Set("pods[5].name", "new-pod")
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Get("pods[0].name")).To(Equal("pod-1"))
			Expect(status.Get("pods[5].name")).To(Equal("new-pod"))
			err = status.Set("pods[1].containers[0].name", "updated-app")
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Get("pods[1].containers[0].name")).To(Equal("updated-app"))
		})

		It("creates new arrays when they don't exist", func() {
			err := status.Set("newArray[2].value", "test")
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Get("newArray[2].value")).To(Equal("test"))
			// Array should be created with proper length
			Expect(status.Get("newArray")).ToNot(BeNil())
		})

		It("handles different value types", func() {
			err := status.Set("stringValue", "test")
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Get("stringValue")).To(Equal("test"))

			err = status.Set("intValue", 42)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Get("intValue")).To(Equal(42))

			err = status.Set("boolValue", true)
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Get("boolValue")).To(Equal(true))

			err = status.Set("sliceValue", []string{"a", "b", "c"})
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Get("sliceValue")).To(Equal([]any{"a", "b", "c"}))
		})

		It("returns error for empty path", func() {
			err := status.Set("", "value")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("path cannot be empty"))
		})

		It("initializes data map if nil", func() {
			emptyStatus := &StatusImpl{}
			err := emptyStatus.Set("test", "value")
			Expect(err).ToNot(HaveOccurred())
			Expect(emptyStatus.Get("test")).To(Equal("value"))
		})

		It("handles complex nested array creation", func() {
			err := status.Set("complex[0].nested[1].deep[2].value", "deep-value")
			Expect(err).ToNot(HaveOccurred())
			Expect(status.Get("complex[0].nested[1].deep[2].value")).To(Equal("deep-value"))
		})
	})

	Describe("Remove", func() {
		It("removes top-level values", func() {
			Expect(status.Remove("phase")).To(Succeed())
			Expect(status.Get("phase")).To(BeNil())
		})

		It("removes nested values", func() {
			Expect(status.Remove("conditions.available.status")).To(Succeed())
			Expect(status.Get("conditions.available.status")).To(BeNil())
			Expect(status.Get("conditions.available.reason")).To(Equal("ReconciliationSucceeded"))
		})

		It("doesn't error for non-existing paths", func() {
			Expect(status.Remove("nonexisting")).To(Succeed())
		})

		It("errors when the path is invalid", func() {
			Expect(status.Remove("replicas.something")).To(MatchError(
				ContainSubstring("expected an object"),
			))
			Expect(status.Remove("")).To(MatchError("path cannot be empty"))
		})

		It("removes entire nested objects", func() {
			Expect(status.Remove("conditions.available")).To(Succeed())
			Expect(status.Get("conditions.available")).To(BeNil())
			Expect(status.Get("conditions.degraded")).ToNot(BeNil())
		})
	})
})
