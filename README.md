# go2jsonc

go2jsonc is a standalone program, and a go generator that creates jsonc files
from a structure in go, including documentation blocks, comments and default
values, facilitating, for example, the maintenance of configuration file
templates.

Any struct embedded or nested in the starting one will be included, even when
coming from different packages.

## Default values

If a function that matches the signature:
```go
func StructTypeNameDefaults() *StructTypeName
```
exists in the same package for which you are generating the jsonc code, it
will be parsed to extract the default values for each field. Note that for
simplicity of parsing the body must have the syntax visible below in the
rendering example.

## Rendering example

Source code:
```go
package multipkg

import (
    "github.com/marco-sacchi/go2jsonc/testdata/multipkg/network"
    alias "github.com/marco-sacchi/go2jsonc/testdata/multipkg/stats"
)

//go:generate go2jsonc -type=MultiPackage -out multi_package.jsonc

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
```

```go
package network

type ConnState int

const (
    // StateDisconnected signals the Disconnected state.
    StateDisconnected ConnState = iota // StateDisconnected comment.
    StateConnecting                    // StateConnecting comment.
    // StateConnected signals the Connected state.
    StateConnected // StateConnected comment.
)

const (
    // StateFailed signals the Failed state.
    StateFailed ConnState = iota + 5 // StateFailed comment.
    // StateReconnecting signals the Reconnecting state.
    StateReconnecting // StateReconnecting comment.
)

// Status reports connection status.
type Status struct {
    Connected bool      // Connected flag comment.
    State     ConnState // Connection state comment.
}
```

```go
package stats

// Info reports statistical info.
type Info struct {
    // PacketLoss documentation block.
    PacketLoss    int `json:"packet_loss"`     // Packet loss comment.
    RoundTripTime int `json:"round_trip_time"` // Round-trip time in milliseconds.
}

```

Rendered output:
```json5
{
    // network.Status - Network status.
    "NetStatus": {
        // bool - Connected flag comment.
        "Connected": true,
        // network.ConnState - Connection state comment.
        // Allowed values:
        // StateDisconnected = 0
        // StateConnecting = 1
        // StateConnected = 2
        // StateFailed = 5
        // StateReconnecting = 6
        "State": 0,
    },
    // int - PacketLoss documentation block.
    // Packet loss comment.
    "packet_loss": 64,
    // int - Round-trip time in milliseconds.
    "round_trip_time": 123,
}

```

## Run as a standalone program

When run as a standalone program, the syntax is as follows:
```
go2jsonc -type=type-name -out=outfile.jsonc [package-dir]
```

When `-out` flag is omitted the code will be written to stdout.

The optional `package-dir` argument can be used when you want to run go2jsonc
from a path outside the package that contains the specified type.

## Running as a generator

To run go2jsonc as a generator just add the `go:generate` comments where
desired, using the same syntax used above:
```
//go:generate go2jsonc -type=type-name -out=outfile.jsonc [package-dir]
```

`package-dir` can be safely omitted in this use case. The directory of the file
in which the comment is present will be used.

## Minor notes

go2jsonc is written in such a way that parsing and extraction of the necessary
data from the AST and the types are decoupled from the generation of the jsonc
file, making it possible to write virtually any format.

## License

MIT License, Copyright (c) 2022 Marco Sacchi
