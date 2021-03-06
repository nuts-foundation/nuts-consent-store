// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/consent.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	pkg "github.com/nuts-foundation/nuts-consent-store/pkg"
	reflect "reflect"
	time "time"
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
func (m *MockConsentStoreClient) ConsentAuth(context context.Context, custodian, subject, actor, dataClass string, checkpoint *time.Time) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConsentAuth", context, custodian, subject, actor, dataClass, checkpoint)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConsentAuth indicates an expected call of ConsentAuth
func (mr *MockConsentStoreClientMockRecorder) ConsentAuth(context, custodian, subject, actor, dataClass, checkpoint interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConsentAuth", reflect.TypeOf((*MockConsentStoreClient)(nil).ConsentAuth), context, custodian, subject, actor, dataClass, checkpoint)
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

// QueryConsent mocks base method
func (m *MockConsentStoreClient) QueryConsent(context context.Context, actor, custodian, subject *string, validAt *time.Time) ([]pkg.PatientConsent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryConsent", context, actor, custodian, subject, validAt)
	ret0, _ := ret[0].([]pkg.PatientConsent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryConsent indicates an expected call of QueryConsent
func (mr *MockConsentStoreClientMockRecorder) QueryConsent(context, actor, custodian, subject, validAt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryConsent", reflect.TypeOf((*MockConsentStoreClient)(nil).QueryConsent), context, actor, custodian, subject, validAt)
}

// DeleteConsentRecordByHash mocks base method
func (m *MockConsentStoreClient) DeleteConsentRecordByHash(context context.Context, consentRecordHash string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteConsentRecordByHash", context, consentRecordHash)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteConsentRecordByHash indicates an expected call of DeleteConsentRecordByHash
func (mr *MockConsentStoreClientMockRecorder) DeleteConsentRecordByHash(context, consentRecordHash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteConsentRecordByHash", reflect.TypeOf((*MockConsentStoreClient)(nil).DeleteConsentRecordByHash), context, consentRecordHash)
}

// FindConsentRecordByHash mocks base method
func (m *MockConsentStoreClient) FindConsentRecordByHash(context context.Context, consentRecordHash string, latest bool) (pkg.ConsentRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindConsentRecordByHash", context, consentRecordHash, latest)
	ret0, _ := ret[0].(pkg.ConsentRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindConsentRecordByHash indicates an expected call of FindConsentRecordByHash
func (mr *MockConsentStoreClientMockRecorder) FindConsentRecordByHash(context, consentRecordHash, latest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindConsentRecordByHash", reflect.TypeOf((*MockConsentStoreClient)(nil).FindConsentRecordByHash), context, consentRecordHash, latest)
}
