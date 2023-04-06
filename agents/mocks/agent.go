// Code generated by MockGen. DO NOT EDIT.
// Source: agent.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	agents "github.com/natexcvi/go-llm/agents"
)

// MockAgent is a mock of Agent interface.
type MockAgent struct {
	ctrl     *gomock.Controller
	recorder *MockAgentMockRecorder
}

// MockAgentMockRecorder is the mock recorder for MockAgent.
type MockAgentMockRecorder struct {
	mock *MockAgent
}

// NewMockAgent creates a new mock instance.
func NewMockAgent(ctrl *gomock.Controller) *MockAgent {
	mock := &MockAgent{ctrl: ctrl}
	mock.recorder = &MockAgentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAgent) EXPECT() *MockAgentMockRecorder {
	return m.recorder
}

// Run mocks base method.
func (m *MockAgent) Run(input agents.T) (agents.S, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", input)
	ret0, _ := ret[0].(agents.S)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Run indicates an expected call of Run.
func (mr *MockAgentMockRecorder) Run(input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockAgent)(nil).Run), input)
}
