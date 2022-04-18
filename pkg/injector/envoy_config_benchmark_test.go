package injector

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/openservicemesh/osm/pkg/certificate/providers/tresor"
	"github.com/openservicemesh/osm/pkg/configurator"
	"github.com/openservicemesh/osm/pkg/constants"
	"github.com/openservicemesh/osm/pkg/logger"
)

func BenchmarkGetEnvoyConfigYAML(b *testing.B) {
	b.StopTimer()
	logger.SetLogLevel("error")

	probes := healthProbes{
		liveness:  &healthProbe{path: "/liveness", port: 81, isHTTP: true},
		readiness: &healthProbe{path: "/readiness", port: 82, isHTTP: true},
		startup:   &healthProbe{path: "/startup", port: 83, isHTTP: true},
	}
	cert := tresor.NewFakeCertificate()
	config := envoyBootstrapConfigMeta{
		NodeID:   cert.GetCommonName().String(),
		RootCert: cert.GetIssuingCA(),
		Cert:     cert.GetCertificateChain(),
		Key:      cert.GetPrivateKey(),

		EnvoyAdminPort: 15000,

		XDSClusterName: constants.OSMControllerName,
		XDSHost:        "osm-controller.b.svc.cluster.local",
		XDSPort:        15128,

		OriginalHealthProbes:  probes,
	}

	mockCtrl := gomock.NewController(b)
	mockConfigurator := configurator.NewMockConfigurator(mockCtrl)

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		getEnvoyConfigYAML(config, mockConfigurator)
		b.StopTimer()
	}
}