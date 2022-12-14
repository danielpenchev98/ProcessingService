// Code generated by MockGen. DO NOT EDIT.
// Source: content.go

// Package pkg_mocks is a generated GoMock package.
package pkg_mocks

import (
	script "danielpenchev98/http-job-processing-service/pkg/script"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockContentBuilderCreator is a mock of ContentBuilderCreator interface.
type MockContentBuilderCreator struct {
	ctrl     *gomock.Controller
	recorder *MockContentBuilderCreatorMockRecorder
}

// MockContentBuilderCreatorMockRecorder is the mock recorder for MockContentBuilderCreator.
type MockContentBuilderCreatorMockRecorder struct {
	mock *MockContentBuilderCreator
}

// NewMockContentBuilderCreator creates a new mock instance.
func NewMockContentBuilderCreator(ctrl *gomock.Controller) *MockContentBuilderCreator {
	mock := &MockContentBuilderCreator{ctrl: ctrl}
	mock.recorder = &MockContentBuilderCreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContentBuilderCreator) EXPECT() *MockContentBuilderCreatorMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockContentBuilderCreator) Create() script.ContentBuilder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create")
	ret0, _ := ret[0].(script.ContentBuilder)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockContentBuilderCreatorMockRecorder) Create() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockContentBuilderCreator)(nil).Create))
}
