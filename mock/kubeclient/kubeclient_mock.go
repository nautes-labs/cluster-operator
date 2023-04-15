// Code generated by MockGen. DO NOT EDIT.
// Source: internal/pkg/vault/vault-store.go

// Package mock_vault is a generated GoMock package.
package mock_vault

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1alpha1 "github.com/nautes-labs/pkg/api/v1alpha1"
	v1 "github.com/nautes-labs/vault-proxy/api/vaultproxy/v1"
	v10 "k8s.io/api/apps/v1"
	v11 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MockKubernetesClient is a mock of KubernetesClient interface.
type MockKubernetesClient struct {
	ctrl     *gomock.Controller
	recorder *MockKubernetesClientMockRecorder
}

// MockKubernetesClientMockRecorder is the mock recorder for MockKubernetesClient.
type MockKubernetesClientMockRecorder struct {
	mock *MockKubernetesClient
}

// NewMockKubernetesClient creates a new mock instance.
func NewMockKubernetesClient(ctrl *gomock.Controller) *MockKubernetesClient {
	mock := &MockKubernetesClient{ctrl: ctrl}
	mock.recorder = &MockKubernetesClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKubernetesClient) EXPECT() *MockKubernetesClientMockRecorder {
	return m.recorder
}

// GetCluster mocks base method.
func (m *MockKubernetesClient) GetCluster(ctx context.Context, name, namespace string) (*v1alpha1.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCluster", ctx, name, namespace)
	ret0, _ := ret[0].(*v1alpha1.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCluster indicates an expected call of GetCluster.
func (mr *MockKubernetesClientMockRecorder) GetCluster(ctx, name, namespace interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCluster", reflect.TypeOf((*MockKubernetesClient)(nil).GetCluster), ctx, name, namespace)
}

// GetSecret mocks base method.
func (m *MockKubernetesClient) GetSecret(ctx context.Context, name, namespace string) (*v11.Secret, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecret", ctx, name, namespace)
	ret0, _ := ret[0].(*v11.Secret)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecret indicates an expected call of GetSecret.
func (mr *MockKubernetesClientMockRecorder) GetSecret(ctx, name, namespace interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecret", reflect.TypeOf((*MockKubernetesClient)(nil).GetSecret), ctx, name, namespace)
}

// GetServiceAccount mocks base method.
func (m *MockKubernetesClient) GetServiceAccount(ctx context.Context, name, namespace string) (*v11.ServiceAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServiceAccount", ctx, name, namespace)
	ret0, _ := ret[0].(*v11.ServiceAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetServiceAccount indicates an expected call of GetServiceAccount.
func (mr *MockKubernetesClientMockRecorder) GetServiceAccount(ctx, name, namespace interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServiceAccount", reflect.TypeOf((*MockKubernetesClient)(nil).GetServiceAccount), ctx, name, namespace)
}

// ListStatefulSets mocks base method.
func (m *MockKubernetesClient) ListStatefulSets(ctx context.Context, namespace string, opts v12.ListOptions) (*v10.StatefulSetList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListStatefulSets", ctx, namespace, opts)
	ret0, _ := ret[0].(*v10.StatefulSetList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListStatefulSets indicates an expected call of ListStatefulSets.
func (mr *MockKubernetesClientMockRecorder) ListStatefulSets(ctx, namespace, opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListStatefulSets", reflect.TypeOf((*MockKubernetesClient)(nil).ListStatefulSets), ctx, namespace, opts)
}

// MockClusterAccount is a mock of ClusterAccount interface.
type MockClusterAccount struct {
	ctrl     *gomock.Controller
	recorder *MockClusterAccountMockRecorder
}

// MockClusterAccountMockRecorder is the mock recorder for MockClusterAccount.
type MockClusterAccountMockRecorder struct {
	mock *MockClusterAccount
}

// NewMockClusterAccount creates a new mock instance.
func NewMockClusterAccount(ctrl *gomock.Controller) *MockClusterAccount {
	mock := &MockClusterAccount{ctrl: ctrl}
	mock.recorder = &MockClusterAccountMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClusterAccount) EXPECT() *MockClusterAccountMockRecorder {
	return m.recorder
}

// GetAccount mocks base method.
func (m *MockClusterAccount) GetAccount() *v1.ClusterAccount {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount")
	ret0, _ := ret[0].(*v1.ClusterAccount)
	return ret0
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockClusterAccountMockRecorder) GetAccount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockClusterAccount)(nil).GetAccount))
}