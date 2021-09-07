package eds

import (
	"strconv"
	"strings"

	xds_discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/pkg/errors"

	"github.com/openservicemesh/osm/pkg/catalog"
	"github.com/openservicemesh/osm/pkg/certificate"
	"github.com/openservicemesh/osm/pkg/configurator"
	"github.com/openservicemesh/osm/pkg/endpoint"
	"github.com/openservicemesh/osm/pkg/envoy"
	"github.com/openservicemesh/osm/pkg/envoy/registry"
	"github.com/openservicemesh/osm/pkg/errcode"
	"github.com/openservicemesh/osm/pkg/identity"
	"github.com/openservicemesh/osm/pkg/service"
)

// NewResponse creates a new Endpoint Discovery Response.
func NewResponse(meshCatalog catalog.MeshCataloger, proxy *envoy.Proxy, request *xds_discovery.DiscoveryRequest, _ configurator.Configurator, _ certificate.Manager, _ *registry.ProxyRegistry) ([]types.Resource, error) {
	// If request comes through and requests specific endpoints, just attempt to answer those
	if request != nil && len(request.ResourceNames) > 0 {
		return fulfillEDSRequest(meshCatalog, proxy, request)
	}

	// Otherwise, generate all endpoint configuration for this proxy
	return generateEDSConfig(meshCatalog, proxy)
}

// fulfillEDSRequest replies only to requested EDS endpoints on Discovery Request
func fulfillEDSRequest(meshCatalog catalog.MeshCataloger, proxy *envoy.Proxy, request *xds_discovery.DiscoveryRequest) ([]types.Resource, error) {
	proxyIdentity, err := envoy.GetServiceIdentityFromProxyCertificate(proxy.GetCertificateCommonName())
	if err != nil {
		log.Error().Err(err).Msgf("Error looking up identity for proxy %s", proxy.String())
		return nil, err
	}

	if request == nil {
		return nil, errors.Errorf("Endpoint discovery request for proxy %s cannot be nil", proxyIdentity)
	}

	var rdsResources []types.Resource
	for _, cluster := range request.ResourceNames {
		meshSvc, err := clusterToMeshSvc(cluster)
		if err != nil {
			log.Error().Err(err).Msgf("Error retrieving MeshService from Cluster %s", cluster)
			continue
		}
		endpoints := meshCatalog.ListAllowedUpstreamEndpointsForService(proxyIdentity, meshSvc)
		if len(endpoints) == 0 {
			log.Error().Str(errcode.Kind, errcode.GetErrCodeWithMetric(errcode.ErrEndpointsNotFound)).
				Msgf("Endpoints not found for upstream cluster %s for proxy identity %s, skipping cluster in EDS response", cluster, proxyIdentity)
			continue
		}
		log.Trace().Msgf("Endpoints for upstream cluster %s for downstream proxy identity %s: %v", cluster, proxyIdentity, endpoints)
		loadAssignment := newClusterLoadAssignment(meshSvc, endpoints)
		rdsResources = append(rdsResources, loadAssignment)
	}

	return rdsResources, nil
}

// generateEDSConfig generates all endpoints expected for a given proxy
func generateEDSConfig(meshCatalog catalog.MeshCataloger, proxy *envoy.Proxy) ([]types.Resource, error) {
	proxyIdentity, err := envoy.GetServiceIdentityFromProxyCertificate(proxy.GetCertificateCommonName())
	if err != nil {
		log.Error().Err(err).Msgf("Error looking up identity for proxy %s", proxy.String())
		return nil, err
	}

	allowedEndpoints, err := getUpstreamEndpointsForProxyIdentity(meshCatalog, proxyIdentity)
	if err != nil {
		log.Error().Err(err).Msgf("Error looking up endpoints for proxy %s", proxy.String())
		return nil, err
	}

	var edsResources []types.Resource
	for svc, endpoints := range allowedEndpoints {
		loadAssignment := newClusterLoadAssignment(svc, endpoints)
		edsResources = append(edsResources, loadAssignment)
	}

	return edsResources, nil
}

// clusterToMeshSvc returns the MeshService associated with the given cluster name
func clusterToMeshSvc(cluster string) (service.MeshService, error) {
	splitFunc := func(r rune) bool {
		return r == '/' || r == '|'
	}

	chunks := strings.FieldsFunc(cluster, splitFunc)
	if len(chunks) != 3 {
		return service.MeshService{}, errors.Errorf("Invalid cluster name. Expected: <namespace>/<name>|<port>, got: %s", cluster)
	}

	port, err := strconv.ParseUint(chunks[2], 10, 16)
	if err != nil {
		return service.MeshService{}, errors.Errorf("Invalid cluster port %s, expected int value: %s", chunks[2], err)
	}

	return service.MeshService{
		Namespace: chunks[0],
		Name:      chunks[1],

		// The port always maps to MeshServer.TargetPort and not MeshService.Port because
		// endpoints of a service are derived from it's TargetPort and not Port.
		TargetPort: uint16(port),
	}, nil
}

// getUpstreamEndpointsForProxyIdentity returns only those service endpoints that belong to the allowed upstream service accounts for the proxy
// Note: ServiceIdentity must be in the format "name.namespace" [https://github.com/openservicemesh/osm/issues/3188]
func getUpstreamEndpointsForProxyIdentity(meshCatalog catalog.MeshCataloger, proxyIdentity identity.ServiceIdentity) (map[service.MeshService][]endpoint.Endpoint, error) {
	allowedServicesEndpoints := make(map[service.MeshService][]endpoint.Endpoint)

	for _, dstSvc := range meshCatalog.ListOutboundServicesForIdentity(proxyIdentity) {
		endpoints := meshCatalog.ListAllowedUpstreamEndpointsForService(proxyIdentity, dstSvc)
		if len(endpoints) == 0 {
			log.Error().Str(errcode.Kind, errcode.GetErrCodeWithMetric(errcode.ErrEndpointsNotFound)).
				Msgf("Endpoints not found for upstream MeshService %s for proxy identity %s, skipping cluster in EDS response", dstSvc, proxyIdentity)
			continue
		}
		allowedServicesEndpoints[dstSvc] = endpoints
	}
	log.Trace().Msgf("Allowed outbound service endpoints for proxy with identity %s: %v", proxyIdentity, allowedServicesEndpoints)
	return allowedServicesEndpoints, nil
}
