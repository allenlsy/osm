package ads

import (
	"context"
	"fmt"
	"testing"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	configv1alpha2 "github.com/openservicemesh/osm/pkg/apis/config/v1alpha2"
	"github.com/openservicemesh/osm/pkg/catalog"
	"github.com/openservicemesh/osm/pkg/certificate"
	"github.com/openservicemesh/osm/pkg/certificate/providers/tresor"
	"github.com/openservicemesh/osm/pkg/configurator"
	"github.com/openservicemesh/osm/pkg/constants"
	"github.com/openservicemesh/osm/pkg/envoy"
	"github.com/openservicemesh/osm/pkg/envoy/registry"
	configFake "github.com/openservicemesh/osm/pkg/gen/client/config/clientset/versioned/fake"
	"github.com/openservicemesh/osm/pkg/k8s"
	"github.com/openservicemesh/osm/pkg/logger"
	"github.com/openservicemesh/osm/pkg/service"
	"github.com/openservicemesh/osm/pkg/tests"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)
func BenchmarkSendXDSResponse(b *testing.B) {
	b.StopTimer()
	logger.SetLogLevel("error")

	var (
		mockCtrl         *gomock.Controller
		mockConfigurator *configurator.MockConfigurator
		mockCertManager  *certificate.MockManager
	)

	mockCtrl = gomock.NewController(b)
	mockConfigurator = configurator.NewMockConfigurator(mockCtrl)
	mockCertManager = certificate.NewMockManager(mockCtrl)

	// --- setup
	namespace := tests.Namespace

	kubeClient := testclient.NewSimpleClientset()
	configClient := configFake.NewSimpleClientset()
	proxyService := service.MeshService{
		Name:      tests.BookstoreV1ServiceName,
		Namespace: namespace,
	}
	proxySvcAccount := tests.BookstoreServiceAccount

	mockConfigurator.EXPECT().GetCertKeyBitSize().Return(2048).AnyTimes()

	proxyUUID := uuid.New()
	labels := map[string]string{constants.EnvoyUniqueIDLabelName: proxyUUID.String()}
	mc := catalog.NewFakeMeshCatalog(kubeClient, configClient)
	proxyRegistry := registry.NewProxyRegistry(registry.ExplicitProxyServiceMapper(func(*envoy.Proxy) ([]service.MeshService, error) {
		return nil, nil
	}), nil)


	pod := tests.NewPodFixture(namespace, fmt.Sprintf("pod-0-%s", proxyUUID), tests.BookstoreServiceAccountName, tests.PodLabels)
	pod.Labels[constants.EnvoyUniqueIDLabelName] = proxyUUID.String()
	_, err := kubeClient.CoreV1().Pods(namespace).Create(context.TODO(), &pod, metav1.CreateOptions{})
	if err != nil {
		b.Fatalf("Failed to create pod: %v", err)
	}

	svc := tests.NewServiceFixture(proxyService.Name, namespace, labels)
	_, err = kubeClient.CoreV1().Services(namespace).Create(context.TODO(), svc, metav1.CreateOptions{})
	if err != nil {
		b.Fatalf("Failed to create service: %v", err)
	}


	certManager := tresor.NewFake(nil)
	certCommonName := envoy.NewXDSCertCommonName(proxyUUID, envoy.KindSidecar, proxySvcAccount.Name, proxySvcAccount.Namespace)
	certDuration := 1 * time.Hour
	certPEM, _ := certManager.IssueCertificate(certCommonName, certDuration)
	cert, _ := certificate.DecodePEMCertificate(certPEM.GetCertificateChain())
	server, _ := tests.NewFakeXDSServer(cert, nil, nil)
	kubectrlMock := k8s.NewMockController(mockCtrl)
	certSerialNumber := certificate.SerialNumber("123456")
	proxy, _ := envoy.NewProxy(certCommonName, certSerialNumber, nil)

	mockConfigurator.EXPECT().IsEgressEnabled().Return(false).AnyTimes()
	mockConfigurator.EXPECT().IsTracingEnabled().Return(false).AnyTimes()
	mockConfigurator.EXPECT().IsPermissiveTrafficPolicyMode().Return(false).AnyTimes()
	mockConfigurator.EXPECT().GetServiceCertValidityPeriod().Return(certDuration).AnyTimes()
	mockConfigurator.EXPECT().GetCertKeyBitSize().Return(2048).AnyTimes()
	mockConfigurator.EXPECT().IsDebugServerEnabled().Return(true).AnyTimes()
	mockConfigurator.EXPECT().GetFeatureFlags().Return(configv1alpha2.FeatureFlags{
		EnableWASMStats:    false,
		EnableEgressPolicy: false,
	}).AnyTimes()
	mockConfigurator.EXPECT().GetMeshConfig().AnyTimes()

	s := NewADSServer(mc, proxyRegistry, true, tests.Namespace, mockConfigurator, mockCertManager, kubectrlMock, nil)

	mockCertManager.EXPECT().IssueCertificate(gomock.Any(), certDuration).Return(certPEM, nil).AnyTimes()

	// Set subscribed resources for SDS
	proxy.SetSubscribedResources(envoy.TypeSDS, mapset.NewSetWith("service-cert:default/bookstore", "root-cert-for-mtls-inbound:default/bookstore"))

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		// s.getTypeResources(proxy, nil)
		s.sendResponse(proxy, &server, nil, mockConfigurator, envoy.XDSResponseOrder...)
		b.StopTimer()
	}
}