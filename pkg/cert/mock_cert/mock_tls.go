// Code generated by MockGen. DO NOT EDIT.
// Source: ./tls.go
//
// Generated by this command:
//
//	mockgen -source=./tls.go -destination=./mock_cert/mock_tls.go -package=mock_cert
//

// Package mock_cert is a generated GoMock package.
package mock_cert

import (
	x509 "crypto/x509"
	pkix "crypto/x509/pkix"
	reflect "reflect"
	time "time"

	cert "github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	gomock "go.uber.org/mock/gomock"
)

// MockTLSCertificate is a mock of TLSCertificate interface.
type MockTLSCertificate struct {
	ctrl     *gomock.Controller
	recorder *MockTLSCertificateMockRecorder
}

// MockTLSCertificateMockRecorder is the mock recorder for MockTLSCertificate.
type MockTLSCertificateMockRecorder struct {
	mock *MockTLSCertificate
}

// NewMockTLSCertificate creates a new mock instance.
func NewMockTLSCertificate(ctrl *gomock.Controller) *MockTLSCertificate {
	mock := &MockTLSCertificate{ctrl: ctrl}
	mock.recorder = &MockTLSCertificateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTLSCertificate) EXPECT() *MockTLSCertificateMockRecorder {
	return m.recorder
}

// CreateCertificateAndSign mocks base method.
func (m *MockTLSCertificate) CreateCertificateAndSign(certificateSubject pkix.Name, subjectPublicKey cert.KeyData, validDuration time.Duration) (cert.TLSCertificate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCertificateAndSign", certificateSubject, subjectPublicKey, validDuration)
	ret0, _ := ret[0].(cert.TLSCertificate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCertificateAndSign indicates an expected call of CreateCertificateAndSign.
func (mr *MockTLSCertificateMockRecorder) CreateCertificateAndSign(certificateSubject, subjectPublicKey, validDuration any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCertificateAndSign", reflect.TypeOf((*MockTLSCertificate)(nil).CreateCertificateAndSign), certificateSubject, subjectPublicKey, validDuration)
}

// GetCACertificate mocks base method.
func (m *MockTLSCertificate) GetCACertificate() (*x509.Certificate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCACertificate")
	ret0, _ := ret[0].(*x509.Certificate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCACertificate indicates an expected call of GetCACertificate.
func (mr *MockTLSCertificateMockRecorder) GetCACertificate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCACertificate", reflect.TypeOf((*MockTLSCertificate)(nil).GetCACertificate))
}

// GetCertificate mocks base method.
func (m *MockTLSCertificate) GetCertificate() *x509.Certificate {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCertificate")
	ret0, _ := ret[0].(*x509.Certificate)
	return ret0
}

// GetCertificate indicates an expected call of GetCertificate.
func (mr *MockTLSCertificateMockRecorder) GetCertificate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCertificate", reflect.TypeOf((*MockTLSCertificate)(nil).GetCertificate))
}

// GetCertificateChains mocks base method.
func (m *MockTLSCertificate) GetCertificateChains(pemEncoded bool) [][]byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCertificateChains", pemEncoded)
	ret0, _ := ret[0].([][]byte)
	return ret0
}

// GetCertificateChains indicates an expected call of GetCertificateChains.
func (mr *MockTLSCertificateMockRecorder) GetCertificateChains(pemEncoded any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCertificateChains", reflect.TypeOf((*MockTLSCertificate)(nil).GetCertificateChains), pemEncoded)
}

// GetCertificateInPEM mocks base method.
func (m *MockTLSCertificate) GetCertificateInPEM() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCertificateInPEM")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCertificateInPEM indicates an expected call of GetCertificateInPEM.
func (mr *MockTLSCertificateMockRecorder) GetCertificateInPEM() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCertificateInPEM", reflect.TypeOf((*MockTLSCertificate)(nil).GetCertificateInPEM))
}

// GetCertificateSubjectSerialNumber mocks base method.
func (m *MockTLSCertificate) GetCertificateSubjectSerialNumber() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCertificateSubjectSerialNumber")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetCertificateSubjectSerialNumber indicates an expected call of GetCertificateSubjectSerialNumber.
func (mr *MockTLSCertificateMockRecorder) GetCertificateSubjectSerialNumber() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCertificateSubjectSerialNumber", reflect.TypeOf((*MockTLSCertificate)(nil).GetCertificateSubjectSerialNumber))
}

// GetNodeID mocks base method.
func (m *MockTLSCertificate) GetNodeID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodeID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetNodeID indicates an expected call of GetNodeID.
func (mr *MockTLSCertificateMockRecorder) GetNodeID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodeID", reflect.TypeOf((*MockTLSCertificate)(nil).GetNodeID))
}

// GetParentTLSCertificate mocks base method.
func (m *MockTLSCertificate) GetParentTLSCertificate() cert.TLSCertificate {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetParentTLSCertificate")
	ret0, _ := ret[0].(cert.TLSCertificate)
	return ret0
}

// GetParentTLSCertificate indicates an expected call of GetParentTLSCertificate.
func (mr *MockTLSCertificateMockRecorder) GetParentTLSCertificate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetParentTLSCertificate", reflect.TypeOf((*MockTLSCertificate)(nil).GetParentTLSCertificate))
}

// GetPrivateKey mocks base method.
func (m *MockTLSCertificate) GetPrivateKey() cert.KeyData {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPrivateKey")
	ret0, _ := ret[0].(cert.KeyData)
	return ret0
}

// GetPrivateKey indicates an expected call of GetPrivateKey.
func (mr *MockTLSCertificateMockRecorder) GetPrivateKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPrivateKey", reflect.TypeOf((*MockTLSCertificate)(nil).GetPrivateKey))
}

// GetPublicKey mocks base method.
func (m *MockTLSCertificate) GetPublicKey() cert.KeyData {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPublicKey")
	ret0, _ := ret[0].(cert.KeyData)
	return ret0
}

// GetPublicKey indicates an expected call of GetPublicKey.
func (mr *MockTLSCertificateMockRecorder) GetPublicKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPublicKey", reflect.TypeOf((*MockTLSCertificate)(nil).GetPublicKey))
}

// IsCA mocks base method.
func (m *MockTLSCertificate) IsCA() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsCA")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsCA indicates an expected call of IsCA.
func (mr *MockTLSCertificateMockRecorder) IsCA() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsCA", reflect.TypeOf((*MockTLSCertificate)(nil).IsCA))
}

// SaveCertificateToFile mocks base method.
func (m *MockTLSCertificate) SaveCertificateToFile(certificateFilePath, privateKeyFilePath string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveCertificateToFile", certificateFilePath, privateKeyFilePath)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveCertificateToFile indicates an expected call of SaveCertificateToFile.
func (mr *MockTLSCertificateMockRecorder) SaveCertificateToFile(certificateFilePath, privateKeyFilePath any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveCertificateToFile", reflect.TypeOf((*MockTLSCertificate)(nil).SaveCertificateToFile), certificateFilePath, privateKeyFilePath)
}

// ValidateSignature mocks base method.
func (m *MockTLSCertificate) ValidateSignature(underValidateCertificate *x509.Certificate) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateSignature", underValidateCertificate)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateSignature indicates an expected call of ValidateSignature.
func (mr *MockTLSCertificateMockRecorder) ValidateSignature(underValidateCertificate any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateSignature", reflect.TypeOf((*MockTLSCertificate)(nil).ValidateSignature), underValidateCertificate)
}
