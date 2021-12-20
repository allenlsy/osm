package injector

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/openservicemesh/osm/pkg/certificate"
	"github.com/openservicemesh/osm/pkg/constants"
	"github.com/openservicemesh/osm/pkg/envoy"
	"github.com/openservicemesh/osm/pkg/metricsstore"
	"github.com/openservicemesh/osm/pkg/version"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (

	certFile = "cert.oem"
	keyFile = "key.pem"
	rootCAFile = "ca.pem"
)

func (wh *mutatingWebhook) createPatchProxyless(pod *corev1.Pod, req *admissionv1.AdmissionRequest, proxyUUID uuid.UUID) ([]byte, error) {
	namespace := req.Namespace

	// Issue a certificate for the proxy sidecar - used for Envoy to connect to XDS (not Envoy-to-Envoy connections)
	cn := envoy.NewXDSCertCommonName(proxyUUID, envoy.KindSidecar, pod.Spec.ServiceAccountName, namespace)
	log.Debug().Msgf("Patching POD spec: service-account=%s, namespace=%s with certificate CN=%s", pod.Spec.ServiceAccountName, namespace, cn)
	startTime := time.Now()
	bootstrapCertificate, err := wh.certManager.IssueCertificate(cn, constants.XDSCertificateValidityPeriod)
	if err != nil {
		log.Error().Err(err).Msgf("Error issuing bootstrap certificate for Envoy with CN=%s", cn)
		return nil, err
	}
	elapsed := time.Since(startTime)

	metricsstore.DefaultMetricsStore.CertIssuedCount.Inc()
	metricsstore.DefaultMetricsStore.CertIssuedTime.
		WithLabelValues().Observe(elapsed.Seconds())
	originalHealthProbes := rewriteHealthProbes(pod)

	// Create the bootstrap configuration for the gRPC client for the given pod
	grpcBootstrapConfigName := fmt.Sprintf("grpc-bootstrap-config-%s", proxyUUID)
	grpcBootstrapSecretsName := fmt.Sprintf("grpc-bootstrap-secrets-%s", proxyUUID)

	// The webhook has a side effect (making out-of-band changes) of creating k8s secret
	// corresponding to the gRPC bootstrap config. Such a side effect needs to be skipped
	// when the request is a DryRun.
	// Ref: https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#side-effects
	if req.DryRun != nil && *req.DryRun {
		log.Debug().Msgf("Skipping gRPC bootstrap config creation for dry-run request: service-account=%s, namespace=%s", pod.Spec.ServiceAccountName, namespace)
	} else if err = wh.createGrpcBootstrapConfig(grpcBootstrapConfigName, grpcBootstrapSecretsName, namespace, wh.osmNamespace, bootstrapCertificate, proxyUUID); err != nil {
		log.Error().Err(err).Msgf("Failed to create gRPC bootstrap config for pod: service-account=%s, namespace=%s, certificate CN=%s", pod.Spec.ServiceAccountName, namespace, cn)
		return nil, err
	}

	// Create volume for envoy TLS secret
	pod.Spec.Volumes = append(
		pod.Spec.Volumes,
		corev1.Volume{
			Name: grpcBootstrapConfigVolume,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: grpcBootstrapConfigName,
					},
				},
			},
		},
		corev1.Volume{
			Name: grpcBootstrapSecretsVolume,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: grpcBootstrapSecretsName,
				},
			},
		},
	)

	// On Windows we cannot use init containers to program HNS because it requires elevated privileges
	// As a result we assume that the HNS redirection policies are already programmed via a CNI plugin.
	// Skip adding the init container and only patch the pod spec with sidecar container.
	podOS := pod.Spec.NodeSelector["kubernetes.io/os"]
	if err := wh.verifyPrerequisites(podOS); err != nil {
		return nil, err
	}
	if !strings.EqualFold(podOS, constants.OSWindows) {
		// Build outbound port exclusion list
		podOutboundPortExclusionList, _ := wh.getPortExclusionListForPod(pod, namespace, outboundPortExclusionListAnnotation)
		globalOutboundPortExclusionList := wh.configurator.GetOutboundPortExclusionList()
		outboundPortExclusionList := mergePortExclusionLists(podOutboundPortExclusionList, globalOutboundPortExclusionList)

		// Build inbound port exclusion list
		podInboundPortExclusionList, _ := wh.getPortExclusionListForPod(pod, namespace, inboundPortExclusionListAnnotation)
		globalInboundPortExclusionList := wh.configurator.GetInboundPortExclusionList()
		inboundPortExclusionList := mergePortExclusionLists(podInboundPortExclusionList, globalInboundPortExclusionList)

		// Add the Init Container
		initContainer := getInitContainerSpec(constants.InitContainerName, wh.configurator, wh.configurator.GetOutboundIPRangeExclusionList(), outboundPortExclusionList, inboundPortExclusionList, wh.configurator.IsPrivilegedInitContainer())
		pod.Spec.InitContainers = append(pod.Spec.InitContainers, initContainer)
	}

	// Add the gRPC sidecar
	sidecar := getGrpcSidecarContainerSpec(pod, wh.configurator, originalHealthProbes, podOS)
	pod.Spec.Containers = append(pod.Spec.Containers, sidecar)

	// This will append a label to the pod, which points to the unique Envoy ID used in the
	// xDS certificate for that Envoy. This label will help xDS match the actual pod to the Envoy that
	// connects to xDS (with the certificate's CN matching this label).
	if pod.Labels == nil {
		pod.Labels = make(map[string]string)
	}
	pod.Labels[constants.GRPCUniqueIDLabelName] = proxyUUID.String()

	return json.Marshal(makePatches(req, pod))
}

func (wh *mutatingWebhook) createGrpcBootstrapConfig(configName, secretsName, namespace, osmNamespace string, cert certificate.Certificater, proxyUUID uuid.UUID) error {
	proxyUUIDStr := proxyUUID.String()
	grpcBootstrapOptions := GrpcBootstrapOptions{
		NodeID: proxyUUIDStr,
		ServerURI: fmt.Sprintf("%s:%d", getOSMControllerFQDN(osmNamespace), constants.ADSServerPort),
		Version: TransportV3,
		CertificateProviders: map[string]json.RawMessage{
			proxyUUIDStr: DefaultFileWatcherConfig(
				path.Join(grpcProxySecretsPath, certFile),
				path.Join(grpcProxySecretsPath, keyFile),
				path.Join(grpcProxySecretsPath, rootCAFile),
			),
		},
	}
	bootstrapConfigBytes, err := GrpcBootstrapContents(grpcBootstrapOptions)
	if err != nil {
		log.Error().Err(err).Msg("Error creating grpc bootstrap JSON")
		return err
	}

	bootstrapCM := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: configName,
			Labels: map[string]string{
				constants.OSMAppNameLabelKey:     constants.OSMAppNameLabelValue,
				constants.OSMAppInstanceLabelKey: wh.meshName,
				constants.OSMAppVersionLabelKey:  version.Version,
			},
		},
		Data: map[string]string{
			grpcBootstrapConfigFile: string(bootstrapConfigBytes),
		},
	}
	if existing, err := wh.kubeClient.CoreV1().ConfigMaps(namespace).Get(context.Background(), configName, metav1.GetOptions{}); err == nil {
		log.Debug().Msgf("Updating bootstrap config gRPC: name=%s, namespace=%s", configName, namespace)
		existing.Data = bootstrapCM.Data
		_, err = wh.kubeClient.CoreV1().ConfigMaps(namespace).Update(context.Background(), existing, metav1.UpdateOptions{})
		if err != nil {
			log.Error().Err(err).Msgf("Error updating bootstrap config gRPC")
			return err
		}
	} else {
		log.Debug().Msgf("Creating bootstrap config for gRPC: name=%s, namespace=%s", configName, namespace)
		_, err = wh.kubeClient.CoreV1().ConfigMaps(namespace).Create(context.Background(), &bootstrapCM, metav1.CreateOptions{})
		if err != nil {
			log.Error().Err(err).Msgf("Error creating bootstrap config gRPC")
			return err
		}
	}

	bootstrapSecrets := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretsName,
			Labels: map[string]string{
				constants.OSMAppNameLabelKey:     constants.OSMAppNameLabelValue,
				constants.OSMAppInstanceLabelKey: wh.meshName,
				constants.OSMAppVersionLabelKey:  version.Version,
			},
		},
		Data: map[string][]byte{
			certFile: cert.GetCertificateChain(),
			keyFile: cert.GetPrivateKey(),
			rootCAFile: cert.GetIssuingCA(),
		},
	}
	if existing, err := wh.kubeClient.CoreV1().Secrets(namespace).Get(context.Background(), secretsName, metav1.GetOptions{}); err == nil {
		log.Debug().Msgf("Updating bootstrap secrets gRPC: name=%s, namespace=%s", secretsName, namespace)
		existing.Data = bootstrapSecrets.Data
		_, err = wh.kubeClient.CoreV1().Secrets(namespace).Update(context.Background(), existing, metav1.UpdateOptions{})
		if err != nil {
			log.Error().Err(err).Msgf("Error updating bootstrap secrets gRPC")
			return err
		}
	}  else {
		log.Debug().Msgf("Creating bootstrap secrets for gRPC: name=%s, namespace=%s", secretsName, namespace)
		_, err = wh.kubeClient.CoreV1().Secrets(namespace).Create(context.Background(), &bootstrapSecrets, metav1.CreateOptions{})
		if err != nil {
			log.Error().Err(err).Msgf("Error creating bootstrap secrets gRPC")
			return err
		}
	}

	return nil
}

// DefaultFileWatcherConfig is a helper function to create a default certificate
// provider plugin configuration. The test is expected to have setup the files
// appropriately before this configuration is used to instantiate providers.
func DefaultFileWatcherConfig(certPath, keyPath, caPath string) json.RawMessage {
	return json.RawMessage(fmt.Sprintf(`{
			"plugin_name": "file_watcher",
			"config": {
				"certificate_file": %q,
				"private_key_file": %q,
				"ca_certificate_file": %q,
				"refresh_interval": "600s"
			}
		}`, certPath, keyPath, caPath))
}
