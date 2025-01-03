// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/lib/installer.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockInstaller is a mock of Installer interface.
type MockInstaller struct {
	ctrl     *gomock.Controller
	recorder *MockInstallerMockRecorder
}

// MockInstallerMockRecorder is the mock recorder for MockInstaller.
type MockInstallerMockRecorder struct {
	mock *MockInstaller
}

// NewMockInstaller creates a new mock instance.
func NewMockInstaller(ctrl *gomock.Controller) *MockInstaller {
	mock := &MockInstaller{ctrl: ctrl}
	mock.recorder = &MockInstallerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInstaller) EXPECT() *MockInstallerMockRecorder {
	return m.recorder
}

// Install mocks base method.
func (m *MockInstaller) Install(constraint string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Install", constraint)
}

// Install indicates an expected call of Install.
func (mr *MockInstallerMockRecorder) Install(constraint interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Install", reflect.TypeOf((*MockInstaller)(nil).Install), constraint)
}

// MockResolver is a mock of Resolver interface.
type MockResolver struct {
	ctrl     *gomock.Controller
	recorder *MockResolverMockRecorder
}

// MockResolverMockRecorder is the mock recorder for MockResolver.
type MockResolverMockRecorder struct {
	mock *MockResolver
}

// NewMockResolver creates a new mock instance.
func NewMockResolver(ctrl *gomock.Controller) *MockResolver {
	mock := &MockResolver{ctrl: ctrl}
	mock.recorder = &MockResolverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResolver) EXPECT() *MockResolverMockRecorder {
	return m.recorder
}

// AddNewVersion mocks base method.
func (m *MockResolver) AddNewVersion(version, destination string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNewVersion", version, destination)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNewVersion indicates an expected call of AddNewVersion.
func (mr *MockResolverMockRecorder) AddNewVersion(version, destination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewVersion", reflect.TypeOf((*MockResolver)(nil).AddNewVersion), version, destination)
}

// Implementation mocks base method.
func (m *MockResolver) Implementation() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Implementation")
	ret0, _ := ret[0].(string)
	return ret0
}

// Implementation indicates an expected call of Implementation.
func (mr *MockResolverMockRecorder) Implementation() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Implementation", reflect.TypeOf((*MockResolver)(nil).Implementation))
}

// ListVersions mocks base method.
func (m *MockResolver) ListVersions() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListVersions")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListVersions indicates an expected call of ListVersions.
func (mr *MockResolverMockRecorder) ListVersions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListVersions", reflect.TypeOf((*MockResolver)(nil).ListVersions))
}

// Name mocks base method.
func (m *MockResolver) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockResolverMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockResolver)(nil).Name))
}
