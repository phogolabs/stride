// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/phogolabs/stride/codedom"
	"github.com/phogolabs/stride/service"
)

type SpecResolver struct {
	ResolveStub        func(*openapi3.Swagger) (*codedom.SpecDescriptor, error)
	resolveMutex       sync.RWMutex
	resolveArgsForCall []struct {
		arg1 *openapi3.Swagger
	}
	resolveReturns struct {
		result1 *codedom.SpecDescriptor
		result2 error
	}
	resolveReturnsOnCall map[int]struct {
		result1 *codedom.SpecDescriptor
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *SpecResolver) Resolve(arg1 *openapi3.Swagger) (*codedom.SpecDescriptor, error) {
	fake.resolveMutex.Lock()
	ret, specificReturn := fake.resolveReturnsOnCall[len(fake.resolveArgsForCall)]
	fake.resolveArgsForCall = append(fake.resolveArgsForCall, struct {
		arg1 *openapi3.Swagger
	}{arg1})
	fake.recordInvocation("Resolve", []interface{}{arg1})
	fake.resolveMutex.Unlock()
	if fake.ResolveStub != nil {
		return fake.ResolveStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.resolveReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *SpecResolver) ResolveCallCount() int {
	fake.resolveMutex.RLock()
	defer fake.resolveMutex.RUnlock()
	return len(fake.resolveArgsForCall)
}

func (fake *SpecResolver) ResolveCalls(stub func(*openapi3.Swagger) (*codedom.SpecDescriptor, error)) {
	fake.resolveMutex.Lock()
	defer fake.resolveMutex.Unlock()
	fake.ResolveStub = stub
}

func (fake *SpecResolver) ResolveArgsForCall(i int) *openapi3.Swagger {
	fake.resolveMutex.RLock()
	defer fake.resolveMutex.RUnlock()
	argsForCall := fake.resolveArgsForCall[i]
	return argsForCall.arg1
}

func (fake *SpecResolver) ResolveReturns(result1 *codedom.SpecDescriptor, result2 error) {
	fake.resolveMutex.Lock()
	defer fake.resolveMutex.Unlock()
	fake.ResolveStub = nil
	fake.resolveReturns = struct {
		result1 *codedom.SpecDescriptor
		result2 error
	}{result1, result2}
}

func (fake *SpecResolver) ResolveReturnsOnCall(i int, result1 *codedom.SpecDescriptor, result2 error) {
	fake.resolveMutex.Lock()
	defer fake.resolveMutex.Unlock()
	fake.ResolveStub = nil
	if fake.resolveReturnsOnCall == nil {
		fake.resolveReturnsOnCall = make(map[int]struct {
			result1 *codedom.SpecDescriptor
			result2 error
		})
	}
	fake.resolveReturnsOnCall[i] = struct {
		result1 *codedom.SpecDescriptor
		result2 error
	}{result1, result2}
}

func (fake *SpecResolver) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.resolveMutex.RLock()
	defer fake.resolveMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *SpecResolver) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ service.SpecResolver = new(SpecResolver)
