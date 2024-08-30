// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package v1

import (
	"fmt"
	"strings"
)

const (
	// EdgeLabelAuthn is a EdgeLabel of type authn.
	EdgeLabelAuthn EdgeLabel = "authn"
	// EdgeLabelAuthz is a EdgeLabel of type authz.
	EdgeLabelAuthz EdgeLabel = "authz"
	// EdgeLabelAuthSecret is a EdgeLabel of type auth_secret.
	EdgeLabelAuthSecret EdgeLabel = "auth_secret"
	// EdgeLabelBackupVia is a EdgeLabel of type backup_via.
	EdgeLabelBackupVia EdgeLabel = "backup_via"
	// EdgeLabelCatalog is a EdgeLabel of type catalog.
	EdgeLabelCatalog EdgeLabel = "catalog"
	// EdgeLabelCertIssuer is a EdgeLabel of type cert_issuer.
	EdgeLabelCertIssuer EdgeLabel = "cert_issuer"
	// EdgeLabelConfig is a EdgeLabel of type config.
	EdgeLabelConfig EdgeLabel = "config"
	// EdgeLabelConnectVia is a EdgeLabel of type connect_via.
	EdgeLabelConnectVia EdgeLabel = "connect_via"
	// EdgeLabelExposedBy is a EdgeLabel of type exposed_by.
	EdgeLabelExposedBy EdgeLabel = "exposed_by"
	// EdgeLabelEvent is a EdgeLabel of type event.
	EdgeLabelEvent EdgeLabel = "event"
	// EdgeLabelLocatedOn is a EdgeLabel of type located_on.
	EdgeLabelLocatedOn EdgeLabel = "located_on"
	// EdgeLabelMonitoredBy is a EdgeLabel of type monitored_by.
	EdgeLabelMonitoredBy EdgeLabel = "monitored_by"
	// EdgeLabelOcmBind is a EdgeLabel of type ocm_bind.
	EdgeLabelOcmBind EdgeLabel = "ocm_bind"
	// EdgeLabelOffshoot is a EdgeLabel of type offshoot.
	EdgeLabelOffshoot EdgeLabel = "offshoot"
	// EdgeLabelOps is a EdgeLabel of type ops.
	EdgeLabelOps EdgeLabel = "ops"
	// EdgeLabelPolicy is a EdgeLabel of type policy.
	EdgeLabelPolicy EdgeLabel = "policy"
	// EdgeLabelRecommendedFor is a EdgeLabel of type recommended_for.
	EdgeLabelRecommendedFor EdgeLabel = "recommended_for"
	// EdgeLabelRestoreInto is a EdgeLabel of type restore_into.
	EdgeLabelRestoreInto EdgeLabel = "restore_into"
	// EdgeLabelScaledBy is a EdgeLabel of type scaled_by.
	EdgeLabelScaledBy EdgeLabel = "scaled_by"
	// EdgeLabelSource is a EdgeLabel of type source.
	EdgeLabelSource EdgeLabel = "source"
	// EdgeLabelStorage is a EdgeLabel of type storage.
	EdgeLabelStorage EdgeLabel = "storage"
	// EdgeLabelView is a EdgeLabel of type view.
	EdgeLabelView EdgeLabel = "view"
)

var ErrInvalidEdgeLabel = fmt.Errorf("not a valid EdgeLabel, try [%s]", strings.Join(_EdgeLabelNames, ", "))

var _EdgeLabelNames = []string{
	string(EdgeLabelAuthn),
	string(EdgeLabelAuthz),
	string(EdgeLabelAuthSecret),
	string(EdgeLabelBackupVia),
	string(EdgeLabelCatalog),
	string(EdgeLabelCertIssuer),
	string(EdgeLabelConfig),
	string(EdgeLabelConnectVia),
	string(EdgeLabelExposedBy),
	string(EdgeLabelEvent),
	string(EdgeLabelLocatedOn),
	string(EdgeLabelMonitoredBy),
	string(EdgeLabelOcmBind),
	string(EdgeLabelOffshoot),
	string(EdgeLabelOps),
	string(EdgeLabelPolicy),
	string(EdgeLabelRecommendedFor),
	string(EdgeLabelRestoreInto),
	string(EdgeLabelScaledBy),
	string(EdgeLabelSource),
	string(EdgeLabelStorage),
	string(EdgeLabelView),
}

// EdgeLabelNames returns a list of possible string values of EdgeLabel.
func EdgeLabelNames() []string {
	tmp := make([]string, len(_EdgeLabelNames))
	copy(tmp, _EdgeLabelNames)
	return tmp
}

// EdgeLabelValues returns a list of the values for EdgeLabel
func EdgeLabelValues() []EdgeLabel {
	return []EdgeLabel{
		EdgeLabelAuthn,
		EdgeLabelAuthz,
		EdgeLabelAuthSecret,
		EdgeLabelBackupVia,
		EdgeLabelCatalog,
		EdgeLabelCertIssuer,
		EdgeLabelConfig,
		EdgeLabelConnectVia,
		EdgeLabelExposedBy,
		EdgeLabelEvent,
		EdgeLabelLocatedOn,
		EdgeLabelMonitoredBy,
		EdgeLabelOcmBind,
		EdgeLabelOffshoot,
		EdgeLabelOps,
		EdgeLabelPolicy,
		EdgeLabelRecommendedFor,
		EdgeLabelRestoreInto,
		EdgeLabelScaledBy,
		EdgeLabelSource,
		EdgeLabelStorage,
		EdgeLabelView,
	}
}

// String implements the Stringer interface.
func (x EdgeLabel) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x EdgeLabel) IsValid() bool {
	_, err := ParseEdgeLabel(string(x))
	return err == nil
}

var _EdgeLabelValue = map[string]EdgeLabel{
	"authn":           EdgeLabelAuthn,
	"authz":           EdgeLabelAuthz,
	"auth_secret":     EdgeLabelAuthSecret,
	"backup_via":      EdgeLabelBackupVia,
	"catalog":         EdgeLabelCatalog,
	"cert_issuer":     EdgeLabelCertIssuer,
	"config":          EdgeLabelConfig,
	"connect_via":     EdgeLabelConnectVia,
	"exposed_by":      EdgeLabelExposedBy,
	"event":           EdgeLabelEvent,
	"located_on":      EdgeLabelLocatedOn,
	"monitored_by":    EdgeLabelMonitoredBy,
	"ocm_bind":        EdgeLabelOcmBind,
	"offshoot":        EdgeLabelOffshoot,
	"ops":             EdgeLabelOps,
	"policy":          EdgeLabelPolicy,
	"recommended_for": EdgeLabelRecommendedFor,
	"restore_into":    EdgeLabelRestoreInto,
	"scaled_by":       EdgeLabelScaledBy,
	"source":          EdgeLabelSource,
	"storage":         EdgeLabelStorage,
	"view":            EdgeLabelView,
}

// ParseEdgeLabel attempts to convert a string to a EdgeLabel.
func ParseEdgeLabel(name string) (EdgeLabel, error) {
	if x, ok := _EdgeLabelValue[name]; ok {
		return x, nil
	}
	return EdgeLabel(""), fmt.Errorf("%s is %w", name, ErrInvalidEdgeLabel)
}

// MustParseEdgeLabel converts a string to a EdgeLabel, and panics if is not valid.
func MustParseEdgeLabel(name string) EdgeLabel {
	val, err := ParseEdgeLabel(name)
	if err != nil {
		panic(err)
	}
	return val
}
