> Status: initial draft v0.1, not for production use

# GTS Go Library

A minimal, idiomatic Go library for working with **GTS** ([Global Type System](https://github.com/gts-spec/gts-spec)) identifiers and JSON/JSON Schema artifacts.

## Roadmap

Featureset:

- [x] **OP#1 - ID Validation**: Verify identifier syntax using regex patterns
- [x] **OP#2 - ID Extraction**: Fetch identifiers from JSON objects or JSON Schema documents
- [x] **OP#3 - ID Parsing**: Decompose identifiers into constituent parts (vendor, package, namespace, type, version, etc.)
- [x] **OP#4 - ID Pattern Matching**: Match identifiers against patterns containing wildcards
- [x] **OP#5 - ID to UUID Mapping**: Generate deterministic UUIDs from GTS identifiers
- [x] **OP#6 - Schema Validation**: Validate object instances against their corresponding schemas
- [x] **OP#7 - Relationship Resolution**: Load all schemas and instances, resolve inter-dependencies, and detect broken references
- [x] **OP#8 - Compatibility Checking**: Verify that schemas with different MINOR versions are compatible
- [x] **OP#8.1 - Backward compatibility checking**
- [x] **OP#8.2 - Forward compatibility checking**
- [x] **OP#8.3 - Full compatibility checking**
- [x] **OP#9 - Version Casting**: Transform instances between compatible MINOR versions
- [x] **OP#10 - Query Execution**: Filter identifier collections using the GTS query language
- [x] **OP#11 - Attribute Access**: Retrieve property values and metadata using the attribute selector (`@`)

TODO - need a file with Go code snippets for all Ops above

Other features:

- [x] **Web server** - a non-production web-server with REST API for the operations processing and testing
- [ ] **CLI** - command-line interface for all GTS operations
- [ ] **UUID for instances** - to support UUID as ID in JSON instances
- [ ] **Yaml** - Add YAML files support
- [ ] **TypeSpec support** - Add [typespec.io](https://typespec.io/) files (*.tsp) support

Technical Backlog:

- [ ] **Code coverage** - target is 90%
- [ ] **Documentation** - add documentation for all the features
- [ ] **Interface** - export publicly available interface and keep cli and others private
- [ ] **Final code cleanup** - remove unused code, denormalize, add critical comments, etc.

## Installation

```bash
go get github.com/GlobalTypeSystem/gts-go
```

## Usage

### Library

Import the GTS package in your Go code:

```go
import "github.com/GlobalTypeSystem/gts-go/gts"
```

#### OP#1 - ID Validation

```go
// Validate a GTS ID
if gts.IsValidGtsID("gts.vendor.pkg.ns.type.v1~") {
    fmt.Println("Valid GTS ID")
}

// Get detailed validation result
result := gts.ValidateGtsID("gts.vendor.pkg.ns.type.v1~")
if result.Valid {
    fmt.Printf("Valid: %s\n", result.ID)
} else {
    fmt.Printf("Invalid: %s\n", result.Error)
}
```

#### OP#2 - ID Extraction

```go
// Extract GTS ID from JSON content
content := map[string]any{
    "gtsId": "gts.vendor.pkg.ns.type.v1.0",
    "name":  "My Entity",
}

result := gts.ExtractID(content, nil)
fmt.Printf("ID: %s\n", result.ID)
fmt.Printf("Schema ID: %s\n", result.SchemaID)
```

#### OP#3 - ID Parsing

```go
// Parse a GTS ID into segments
result := gts.ParseGtsID("gts.vendor.pkg.ns.type.v1~")
if result.OK {
    for _, seg := range result.Segments {
        fmt.Printf("Vendor: %s, Package: %s, Type: %s, Version: %d\n",
            seg.Vendor, seg.Package, seg.Type, seg.VerMajor)
    }
}
```

#### OP#4 - Pattern Matching

```go
// Match GTS ID against a pattern
result := gts.MatchIDPattern(
    "gts.vendor.pkg.ns.type.v1.0",
    "gts.vendor.pkg.*",
)
if result.Match {
    fmt.Println("Pattern matched!")
}
```

#### OP#5 - UUID Generation

```go
// Generate deterministic UUID from GTS ID
result := gts.IDToUUID("gts.vendor.pkg.ns.type.v1~")
fmt.Printf("UUID: %s\n", result.UUID)
```

#### Using the GTS Store

```go
// Create a new store
store := gts.NewGtsStore(nil)

// Register an entity
entity := gts.NewJsonEntity(map[string]any{
    "gtsId": "gts.vendor.pkg.ns.type.v1.0",
    "name":  "My Entity",
}, gts.DefaultGtsConfig())

err := store.Register(entity)
if err != nil {
    log.Fatal(err)
}

// Query entities
result := store.Query("gts.vendor.pkg.*", 100)
fmt.Printf("Found %d entities\n", result.Count)

// Validate an instance
validation := store.ValidateInstance("gts.vendor.pkg.ns.type.v1.0")
if validation.OK {
    fmt.Println("Instance is valid")
}

// Attribute access
attr := store.GetAttribute("gts.vendor.pkg.ns.type.v1.0@name")
if attr.Resolved {
    fmt.Printf("Attribute value: %v\n", attr.Value)
}
```

### CLI

```bash
# TODO
```

### Library

TODO - See ...

### Web server

The web server is a non-production web-server with REST API for the operations processing and testing. It implements reference API for gts-spec [tests](https://github.com/GlobalTypeSystem/gts-spec/tree/main/tests)


```bash
# start the web server, default location is http://127.0.0.1:8000
...

# start the web server on different port
...

# pre-populate server with the JSON instancens and schemas from the gts-spec tests
...

# Generate the OpenAPI schema
...

# See the schema
...
```

### Testing

You can test the gts-go library by utilizing the shared test suite from the [gts-spec](https://github.com/GlobalTypeSystem/gts-spec) specification and executing the tests against the web server.

Executing gts-spec Tests on the Server:

```bash
# getting the tests
git clone https://github.com/GlobalTypeSystem/gts-spec.git
cd gts-spec/tests

# run tests against the web server on port 8000 (default)
pytest

# override server URL using GTS_BASE_URL environment variable
GTS_BASE_URL=http://127.0.0.1:8001 pytest

# or set it persistently
export GTS_BASE_URL=http://127.0.0.1:8001
pytest
```

## License

Apache License 2.0
