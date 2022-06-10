package multipkg

import (
	"github.com/marco-sacchi/go2jsonc/testdata/multipkg/network"
	alias "github.com/marco-sacchi/go2jsonc/testdata/multipkg/stats"
)

//go:generate go2jsonc -type MultiPackage -out multi_package.jsonc
//go:generate go2jsonc -type MultiPackage -doc-types NotStructFields -out multi_package_not_struct.jsonc
//go:generate go2jsonc -type MultiPackage -doc-types NotArrayFields -out multi_package_not_array.jsonc
//go:generate go2jsonc -type MultiPackage -doc-types NotMapFields -out multi_package_not_map.jsonc

// MultiPackage tests the multi-package and import aliasing case.
type MultiPackage struct {
	NetStatus  network.Status // Network status.
	alias.Info                // Statistics info.
}

func MultiPackageDefaults() *MultiPackage {
	return &MultiPackage{
		NetStatus: network.Status{
			Connected: true,
			State:     network.StateDisconnected,
		},
		Info: alias.Info{
			PacketLoss:    32 * 2,
			RoundTripTime: 123,
		},
	}
}
