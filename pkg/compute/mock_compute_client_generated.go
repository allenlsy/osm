// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/openservicemesh/osm/pkg/compute (interfaces: Interface)

// Package compute is a generated GoMock package.
package compute

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1alpha2 "github.com/openservicemesh/osm/pkg/apis/config/v1alpha2"
	v1alpha1 "github.com/openservicemesh/osm/pkg/apis/policy/v1alpha1"
	endpoint "github.com/openservicemesh/osm/pkg/endpoint"
	envoy "github.com/openservicemesh/osm/pkg/envoy"
	identity "github.com/openservicemesh/osm/pkg/identity"
	service "github.com/openservicemesh/osm/pkg/service"
	types "k8s.io/apimachinery/pkg/types"
)

// MockInterface is a mock of Interface interface.
type MockInterface struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceMockRecorder
}

// MockInterfaceMockRecorder is the mock recorder for MockInterface.
type MockInterfaceMockRecorder struct {
	mock *MockInterface
}

// NewMockInterface creates a new mock instance.
func NewMockInterface(ctrl *gomock.Controller) *MockInterface {
	mock := &MockInterface{ctrl: ctrl}
	mock.recorder = &MockInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInterface) EXPECT() *MockInterfaceMockRecorder {
	return m.recorder
}

// GetHostnamesForService mocks base method.
func (m *MockInterface) GetHostnamesForService(arg0 service.MeshService, arg1 bool) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHostnamesForService", arg0, arg1)
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetHostnamesForService indicates an expected call of GetHostnamesForService.
func (mr *MockInterfaceMockRecorder) GetHostnamesForService(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHostnamesForService", reflect.TypeOf((*MockInterface)(nil).GetHostnamesForService), arg0, arg1)
}

// GetIngressBackendPolicy mocks base method.
func (m *MockInterface) GetIngressBackendPolicy(arg0 service.MeshService) *v1alpha1.IngressBackend {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIngressBackendPolicy", arg0)
	ret0, _ := ret[0].(*v1alpha1.IngressBackend)
	return ret0
}

// GetIngressBackendPolicy indicates an expected call of GetIngressBackendPolicy.
func (mr *MockInterfaceMockRecorder) GetIngressBackendPolicy(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIngressBackendPolicy", reflect.TypeOf((*MockInterface)(nil).GetIngressBackendPolicy), arg0)
}

// GetMeshConfig mocks base method.
func (m *MockInterface) GetMeshConfig() v1alpha2.MeshConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMeshConfig")
	ret0, _ := ret[0].(v1alpha2.MeshConfig)
	return ret0
}

// GetMeshConfig indicates an expected call of GetMeshConfig.
func (mr *MockInterfaceMockRecorder) GetMeshConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMeshConfig", reflect.TypeOf((*MockInterface)(nil).GetMeshConfig))
}

// GetOSMNamespace mocks base method.
func (m *MockInterface) GetOSMNamespace() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOSMNamespace")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetOSMNamespace indicates an expected call of GetOSMNamespace.
func (mr *MockInterfaceMockRecorder) GetOSMNamespace() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOSMNamespace", reflect.TypeOf((*MockInterface)(nil).GetOSMNamespace))
}

// GetResolvableEndpointsForService mocks base method.
func (m *MockInterface) GetResolvableEndpointsForService(arg0 service.MeshService) []endpoint.Endpoint {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResolvableEndpointsForService", arg0)
	ret0, _ := ret[0].([]endpoint.Endpoint)
	return ret0
}

// GetResolvableEndpointsForService indicates an expected call of GetResolvableEndpointsForService.
func (mr *MockInterfaceMockRecorder) GetResolvableEndpointsForService(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResolvableEndpointsForService", reflect.TypeOf((*MockInterface)(nil).GetResolvableEndpointsForService), arg0)
}

// GetServicesForServiceIdentity mocks base method.
func (m *MockInterface) GetServicesForServiceIdentity(arg0 identity.ServiceIdentity) []service.MeshService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServicesForServiceIdentity", arg0)
	ret0, _ := ret[0].([]service.MeshService)
	return ret0
}

// GetServicesForServiceIdentity indicates an expected call of GetServicesForServiceIdentity.
func (mr *MockInterfaceMockRecorder) GetServicesForServiceIdentity(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServicesForServiceIdentity", reflect.TypeOf((*MockInterface)(nil).GetServicesForServiceIdentity), arg0)
}

// GetTargetPortForServicePort mocks base method.
func (m *MockInterface) GetTargetPortForServicePort(arg0 types.NamespacedName, arg1 uint16) (uint16, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTargetPortForServicePort", arg0, arg1)
	ret0, _ := ret[0].(uint16)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTargetPortForServicePort indicates an expected call of GetTargetPortForServicePort.
func (mr *MockInterfaceMockRecorder) GetTargetPortForServicePort(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTargetPortForServicePort", reflect.TypeOf((*MockInterface)(nil).GetTargetPortForServicePort), arg0, arg1)
}

// GetUpstreamTrafficSettingByHost mocks base method.
func (m *MockInterface) GetUpstreamTrafficSettingByHost(arg0 string) *v1alpha1.UpstreamTrafficSetting {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUpstreamTrafficSettingByHost", arg0)
	ret0, _ := ret[0].(*v1alpha1.UpstreamTrafficSetting)
	return ret0
}

// GetUpstreamTrafficSettingByHost indicates an expected call of GetUpstreamTrafficSettingByHost.
func (mr *MockInterfaceMockRecorder) GetUpstreamTrafficSettingByHost(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUpstreamTrafficSettingByHost", reflect.TypeOf((*MockInterface)(nil).GetUpstreamTrafficSettingByHost), arg0)
}

// GetUpstreamTrafficSettingByNamespace mocks base method.
func (m *MockInterface) GetUpstreamTrafficSettingByNamespace(arg0 *types.NamespacedName) *v1alpha1.UpstreamTrafficSetting {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUpstreamTrafficSettingByNamespace", arg0)
	ret0, _ := ret[0].(*v1alpha1.UpstreamTrafficSetting)
	return ret0
}

// GetUpstreamTrafficSettingByNamespace indicates an expected call of GetUpstreamTrafficSettingByNamespace.
func (mr *MockInterfaceMockRecorder) GetUpstreamTrafficSettingByNamespace(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUpstreamTrafficSettingByNamespace", reflect.TypeOf((*MockInterface)(nil).GetUpstreamTrafficSettingByNamespace), arg0)
}

// GetUpstreamTrafficSettingByService mocks base method.
func (m *MockInterface) GetUpstreamTrafficSettingByService(arg0 *service.MeshService) *v1alpha1.UpstreamTrafficSetting {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUpstreamTrafficSettingByService", arg0)
	ret0, _ := ret[0].(*v1alpha1.UpstreamTrafficSetting)
	return ret0
}

// GetUpstreamTrafficSettingByService indicates an expected call of GetUpstreamTrafficSettingByService.
func (mr *MockInterfaceMockRecorder) GetUpstreamTrafficSettingByService(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUpstreamTrafficSettingByService", reflect.TypeOf((*MockInterface)(nil).GetUpstreamTrafficSettingByService), arg0)
}

// IsMetricsEnabled mocks base method.
func (m *MockInterface) IsMetricsEnabled(arg0 *envoy.Proxy) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsMetricsEnabled", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsMetricsEnabled indicates an expected call of IsMetricsEnabled.
func (mr *MockInterfaceMockRecorder) IsMetricsEnabled(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsMetricsEnabled", reflect.TypeOf((*MockInterface)(nil).IsMetricsEnabled), arg0)
}

// ListEgressPolicies mocks base method.
func (m *MockInterface) ListEgressPolicies(arg0 identity.K8sServiceAccount) []*v1alpha1.Egress {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEgressPolicies", arg0)
	ret0, _ := ret[0].([]*v1alpha1.Egress)
	return ret0
}

// ListEgressPolicies indicates an expected call of ListEgressPolicies.
func (mr *MockInterfaceMockRecorder) ListEgressPolicies(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEgressPolicies", reflect.TypeOf((*MockInterface)(nil).ListEgressPolicies), arg0)
}

// ListEndpointsForIdentity mocks base method.
func (m *MockInterface) ListEndpointsForIdentity(arg0 identity.ServiceIdentity) []endpoint.Endpoint {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEndpointsForIdentity", arg0)
	ret0, _ := ret[0].([]endpoint.Endpoint)
	return ret0
}

// ListEndpointsForIdentity indicates an expected call of ListEndpointsForIdentity.
func (mr *MockInterfaceMockRecorder) ListEndpointsForIdentity(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEndpointsForIdentity", reflect.TypeOf((*MockInterface)(nil).ListEndpointsForIdentity), arg0)
}

// ListEndpointsForService mocks base method.
func (m *MockInterface) ListEndpointsForService(arg0 service.MeshService) []endpoint.Endpoint {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEndpointsForService", arg0)
	ret0, _ := ret[0].([]endpoint.Endpoint)
	return ret0
}

// ListEndpointsForService indicates an expected call of ListEndpointsForService.
func (mr *MockInterfaceMockRecorder) ListEndpointsForService(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEndpointsForService", reflect.TypeOf((*MockInterface)(nil).ListEndpointsForService), arg0)
}

// ListRetryPolicies mocks base method.
func (m *MockInterface) ListRetryPolicies(arg0 identity.K8sServiceAccount) []*v1alpha1.Retry {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRetryPolicies", arg0)
	ret0, _ := ret[0].([]*v1alpha1.Retry)
	return ret0
}

// ListRetryPolicies indicates an expected call of ListRetryPolicies.
func (mr *MockInterfaceMockRecorder) ListRetryPolicies(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRetryPolicies", reflect.TypeOf((*MockInterface)(nil).ListRetryPolicies), arg0)
}

// ListServiceIdentitiesForService mocks base method.
func (m *MockInterface) ListServiceIdentitiesForService(arg0 service.MeshService) []identity.ServiceIdentity {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListServiceIdentitiesForService", arg0)
	ret0, _ := ret[0].([]identity.ServiceIdentity)
	return ret0
}

// ListServiceIdentitiesForService indicates an expected call of ListServiceIdentitiesForService.
func (mr *MockInterfaceMockRecorder) ListServiceIdentitiesForService(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListServiceIdentitiesForService", reflect.TypeOf((*MockInterface)(nil).ListServiceIdentitiesForService), arg0)
}

// ListServices mocks base method.
func (m *MockInterface) ListServices() []service.MeshService {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListServices")
	ret0, _ := ret[0].([]service.MeshService)
	return ret0
}

// ListServices indicates an expected call of ListServices.
func (mr *MockInterfaceMockRecorder) ListServices() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListServices", reflect.TypeOf((*MockInterface)(nil).ListServices))
}

// ListServicesForProxy mocks base method.
func (m *MockInterface) ListServicesForProxy(arg0 *envoy.Proxy) ([]service.MeshService, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListServicesForProxy", arg0)
	ret0, _ := ret[0].([]service.MeshService)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListServicesForProxy indicates an expected call of ListServicesForProxy.
func (mr *MockInterfaceMockRecorder) ListServicesForProxy(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListServicesForProxy", reflect.TypeOf((*MockInterface)(nil).ListServicesForProxy), arg0)
}

// UpdateIngressBackendStatus mocks base method.
func (m *MockInterface) UpdateIngressBackendStatus(arg0 *v1alpha1.IngressBackend) (*v1alpha1.IngressBackend, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateIngressBackendStatus", arg0)
	ret0, _ := ret[0].(*v1alpha1.IngressBackend)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateIngressBackendStatus indicates an expected call of UpdateIngressBackendStatus.
func (mr *MockInterfaceMockRecorder) UpdateIngressBackendStatus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateIngressBackendStatus", reflect.TypeOf((*MockInterface)(nil).UpdateIngressBackendStatus), arg0)
}

// UpdateUpstreamTrafficSettingStatus mocks base method.
func (m *MockInterface) UpdateUpstreamTrafficSettingStatus(arg0 *v1alpha1.UpstreamTrafficSetting) (*v1alpha1.UpstreamTrafficSetting, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUpstreamTrafficSettingStatus", arg0)
	ret0, _ := ret[0].(*v1alpha1.UpstreamTrafficSetting)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUpstreamTrafficSettingStatus indicates an expected call of UpdateUpstreamTrafficSettingStatus.
func (mr *MockInterfaceMockRecorder) UpdateUpstreamTrafficSettingStatus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUpstreamTrafficSettingStatus", reflect.TypeOf((*MockInterface)(nil).UpdateUpstreamTrafficSettingStatus), arg0)
}
