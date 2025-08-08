package kratix

import (
	"github.com/syntasso/kratix/api/v1alpha1"
)

type PromiseAccessor interface {
	ResourceAccessor
	GetPromise() *v1alpha1.Promise
}

type Promise struct {
	Resource
	promise *v1alpha1.Promise
}

var _ PromiseAccessor = (*Promise)(nil)

func (p *Promise) GetPromise() *v1alpha1.Promise {
	return p.promise
}
