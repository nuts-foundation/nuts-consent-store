// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/consent.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	pkg "github.com/nuts-foundation/nuts-consent-store/pkg"
	reflect "reflect"
)

// MockConsentStoreClient is a mock of ConsentStoreClient interface
type MockConsentStoreClient struct {
	ctrl     *gomock.Controller
	recorder *MockConsentStoreClientMockRecorder
}

// MockConsentStoreClientMockRecorder is the mock recorder for MockConsentStoreClient
type MockConsentStoreClientMockRecorder struct {
	mock *MockConsentStoreClient
}

// NewMockConsentStoreClient creates a new mock instance
func NewMockConsentStoreClient(ctrl *gomock.Controller) *MockConsentStoreClient {
	mock := &MockConsentStoreClient{ctrl: ctrl}
	mock.recorder = &MockConsentStoreClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConsentStoreClient) EXPECT() *MockConsentStoreClientMockRecorder {
	return m.recorder
}

// ConsentAuth mocks base method
func (m *MockConsentStoreClient) ConsentAuth(context context.Context, consentRule pkg.PatientConsent, resourceType string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConsentAuth", context, consentRule, resourceType)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConsentAuth indicates an expected call of ConsentAuth
func (mr *MockConsentStoreClientMockRecorder) ConsentAuth(context, consentRule, resourceType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConsentAuth", reflect.TypeOf((*MockConsentStoreClient)(nil).ConsentAuth), context, consentRule, resourceType)
}

// RecordConsent mocks base method
func (m *MockConsentStoreClient) RecordConsent(context context.Context, consent []pkg.PatientConsent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecordConsent", context, consent)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecordConsent indicates an expected call of RecordConsent
func (mr *MockConsentStoreClientMockRecorder) RecordConsent(context, consent interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecordConsent", reflect.TypeOf((*MockConsentStoreClient)(nil).RecordConsent), context, consent)
}

// QueryConsentForActor mocks base method
func (m *MockConsentStoreClient) QueryConsentForActor(context context.Context, actor, query string) ([]pkg.PatientConsent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryConsentForActor", context, actor, query)
	ret0, _ := ret[0].([]pkg.PatientConsent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryConsentForActor indicates an expected call of QueryConsentForActor
func (mr *MockConsentStoreClientMockRecorder) QueryConsentForActor(context, actor, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryConsentForActor", reflect.TypeOf((*MockConsentStoreClient)(nil).QueryConsentForActor), context, actor, query)
}

// QueryConsentForActorAndSubject mocks base method
func (m *MockConsentStoreClient) QueryConsentForActorAndSubject(context context.Context, actor, subject string) ([]pkg.PatientConsent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryConsentForActorAndSubject", context, actor, subject)
	ret0, _ := ret[0].([]pkg.PatientConsent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryConsentForActorAndSubject indicates an expected call of QueryConsentForActorAndSubject
func (mr *MockConsentStoreClientMockRecorder) QueryConsentForActorAndSubject(context, actor, subject interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryConsentForActorAndSubject", reflect.TypeOf((*MockConsentStoreClient)(nil).QueryConsentForActorAndSubject), context, actor, subject)
}