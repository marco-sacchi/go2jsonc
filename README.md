# go2jsonc

go2jsonc is a standalone program, a go generator and a library that creates
jsonc files starting from a structure in go, including documentation blocks,
comments and default values, facilitating, for example, the maintenance of
configuration file templates.

Any struct embedded or nested in the starting one will be included, even when
coming from different packages.

The typed constants are resolved and included in the comments of generated
code to make it easier to verify and modify the values.

go2jsonc is written in such a way that parsing and extraction of the necessary
data from AST and types are decoupled from the generation of the jsonc file, 
making it possible to write virtually any format.

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

//go:generate go2jsonc -type MultiPackage -out multi_package.jsonc

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
    StateDisconnected ConnState = iota
    // StateConnecting signals the connection-pending state.
    StateConnecting
    // StateConnected signals the Connected state.
    StateConnected
)

const (
    // StateFailed signals the Failed state.
    StateFailed ConnState = iota + 5
    // StateReconnecting signals the Reconnecting state.
    StateReconnecting
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
        // StateDisconnected = 0  StateDisconnected signals the Disconnected state.
        // StateConnecting   = 1  StateConnecting signals the connection-pending state.
        // StateConnected    = 2  StateConnected signals the Connected state.
        // StateFailed       = 5  StateFailed signals the Failed state.
        // StateReconnecting = 6  StateReconnecting signals the Reconnecting state.
        "State": 0
    },
    // int - PacketLoss documentation block.
    // Packet loss comment.
    "packet_loss": 64,
    // int - Round-trip time in milliseconds.
    "round_trip_time": 123
}
```

## Installation

To install the standalone program / generator, run the following:

```shell
go install github.com/marco-sacchi/go2jsonc/cmd/go2jsonc@latest
```

## Running as a standalone program

When run as a standalone program, the syntax is as follows:

```shell
go2jsonc -type <type-name> [-doc-types bits] [-out jsonc-filename] [package-dir]
```

- `-doc-types` - `string`: pipe-separated bits representing struct fields types
  for which do not render the type in JSONC comments; when omitted all types
  will be rendered for all fields
- `-out` - `string`: output JSONC filepath; when omitted the code is written
  to `stdout`
- `-type` - `string`: struct type name for which generate JSONC; mandatory
- `package-dir`: directory that contains the go file where specified type is 
  defined; when omitted, current working directory will be used

Allowed constants for `-doc-types` flag:

- `NotStructFields`: do not show type on struct fields
- `NotArrayFields`: do not show type on array or slice fields
- `NotMapFields`: do not show type on map fields
- `NotFields`: do not show type on all fields (override all previous bits).

## Running as a generator

To run go2jsonc as a generator just add the `go:generate` comments where
desired, using the same syntax used above:

```
//go:generate go2jsonc -type type-name [-doc-types bits] -out outfile.jsonc [package-dir]
```

`package-dir` can be safely omitted in this use case. The directory of the file
in which the comment is present will be used.

## Importing packages

go2jsonc contains two packages:

- go2jsonc
- go2jsonc/distiller

You can import the first one and use the function:

```go
func Generate(dir, typeName string, mode DocTypesMode) (string, error)
```

to generate jsonc code from the specified package and type.

Or you can import the latter to easily extract information from the AST and
render other formats.

## Known limitations

go2jsonc supports maps, but under the following limitations:

- the keys must be strings;
- the values can be of any type, except the `interface{}` type, this means
  that all the values of the map must be of the same type;
- maps cannot be nested.

These choices have been made considering that adding map types does not make
it possible to integrate the documentation of the required structure and types,
completely canceling the intent of jsonc or any other output formats that need
comments to be easily understandable.

## License

MIT License, Copyright (c) 2022 Marco Sacchi
