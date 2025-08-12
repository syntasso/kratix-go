package kratix

import (
	"github.com/syntasso/kratix/api/v1alpha1"
)

type Promise interface {
	Resource
	GetPromise() *v1alpha1.Promise
}

type PromiseImpl struct {
	ResourceImpl
	promise *v1alpha1.Promise
}

var _ Promise = (*PromiseImpl)(nil)

func (p *PromiseImpl) GetPromise() *v1alpha1.Promise {
	return p.promise
}
