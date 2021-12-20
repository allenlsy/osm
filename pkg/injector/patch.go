package injector

import (
	"fmt"

	"github.com/openservicemesh/osm/pkg/constants"
)

func getOSMControllerFQDN(osmNamespace string) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local", constants.OSMControllerName, osmNamespace)
}
