package injector

import (
	"path"

	corev1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm/pkg/configurator"
	"github.com/openservicemesh/osm/pkg/constants"
)

const (
	grpcBootstrapConfigFile = "bootstrap.json"
	grpcProxyConfigPath     = "/etc/xds-client"
	grpcProxySecretsPath     = "/etc/xds-client-secrets"
)

func getGrpcSidecarContainerSpec(pod *corev1.Pod, cfg configurator.Configurator, originalHealthProbes healthProbes, podOS string) corev1.Container {
	return corev1.Container{
		Name:            constants.GrpcClientContainerName,
		Image:           constants.GrpcClientImage,
		ImagePullPolicy: corev1.PullIfNotPresent,
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      grpcBootstrapConfigVolume,
				MountPath: grpcProxyConfigPath,
			},
			{
				Name:      grpcBootstrapSecretsVolume,
				MountPath: grpcProxySecretsPath,
			},
		},
		Command:   []string{"sleep", "100000000"},
		Resources: cfg.GetProxyResources(),
		Env: []corev1.EnvVar{
			{
				Name:  "GRPC_XDS_BOOTSTRAP",
				Value: path.Join(grpcProxyConfigPath, grpcBootstrapConfigFile),
			},
			{
				Name:  "GRPC_GO_LOG_VERBOSITY_LEVEL",
				Value: "99",
			},
			{
				Name:  "GRPC_GO_LOG_SEVERITY_LEVEL",
				Value: "info",
			},
		},
	}
}

func getGrpcContainerPorts(originalHealthProbes healthProbes) []corev1.ContainerPort {
	containerPorts := []corev1.ContainerPort{
		{
			Name:          constants.EnvoyAdminPortName,
			ContainerPort: constants.EnvoyAdminPort,
		},
		{
			Name:          constants.EnvoyInboundListenerPortName,
			ContainerPort: constants.EnvoyInboundListenerPort,
		},
		{
			Name:          constants.EnvoyInboundPrometheusListenerPortName,
			ContainerPort: constants.EnvoyPrometheusInboundListenerPort,
		},
	}

	if originalHealthProbes.liveness != nil {
		livenessPort := corev1.ContainerPort{
			// Name must be no more than 15 characters
			Name:          "liveness-port",
			ContainerPort: livenessProbePort,
		}
		containerPorts = append(containerPorts, livenessPort)
	}

	if originalHealthProbes.readiness != nil {
		readinessPort := corev1.ContainerPort{
			// Name must be no more than 15 characters
			Name:          "readiness-port",
			ContainerPort: readinessProbePort,
		}
		containerPorts = append(containerPorts, readinessPort)
	}

	if originalHealthProbes.startup != nil {
		startupPort := corev1.ContainerPort{
			// Name must be no more than 15 characters
			Name:          "startup-port",
			ContainerPort: startupProbePort,
		}
		containerPorts = append(containerPorts, startupPort)
	}

	return containerPorts
}
