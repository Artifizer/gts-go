/*
Copyright © 2025 Global Type System
Released under Apache License 2.0
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/GlobalTypeSystem/gts-go/gts"
	"github.com/GlobalTypeSystem/gts-go/server"
)

const (
	usageText = `GTS helpers CLI

Usage:
  gts [flags] <command> [command-flags]

Global flags:
  -v, --verbose int     Verbosity level (0=silent, 1=info, 2=debug) (default 0)
  --config string       Path to optional GTS config JSON to override defaults
  --path string         Path to json and schema files or directories (global default)

Commands:
  validate-id           Validate a GTS ID format
  parse-id              Parse a GTS ID into its components
  match-id-pattern      Match a GTS ID against a pattern
  uuid                  Generate UUID from a GTS ID
  validate-instance     Validate an instance against its schema
  resolve-relationships Resolve relationships for an entity
  compatibility         Check compatibility between two schemas
  cast                  Cast an instance or schema to a target schema
  query                 Query entities using an expression
  attr                  Get attribute value from a GTS entity
  list                  List all entities
  server                Start the GTS HTTP server
  openapi-spec          Generate OpenAPI specification

Run 'gts <command> -h' for more information on a command.
`
)

type config struct {
	verbose int
	config  string
	path    string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usageText)
		os.Exit(1)
	}

	// Parse global flags
	globalFlags := flag.NewFlagSet("gts", flag.ContinueOnError)
	cfg := &config{}
	globalFlags.IntVar(&cfg.verbose, "v", 0, "Verbosity level")
	globalFlags.IntVar(&cfg.verbose, "verbose", 0, "Verbosity level")
	globalFlags.StringVar(&cfg.config, "config", "", "Path to GTS config JSON")
	globalFlags.StringVar(&cfg.path, "path", "", "Path to json and schema files")
	globalFlags.Usage = func() {
		fmt.Fprint(os.Stderr, usageText)
	}

	// Find the command (first non-flag argument)
	cmdIdx := 1
	for cmdIdx < len(os.Args) && strings.HasPrefix(os.Args[cmdIdx], "-") {
		cmdIdx++
		// Skip flag value if present
		if cmdIdx < len(os.Args) && !strings.HasPrefix(os.Args[cmdIdx], "-") {
			cmdIdx++
		}
	}

	if cmdIdx >= len(os.Args) {
		fmt.Fprint(os.Stderr, usageText)
		os.Exit(1)
	}

	// Parse global flags before the command
	if err := globalFlags.Parse(os.Args[1:cmdIdx]); err != nil {
		os.Exit(1)
	}

	command := os.Args[cmdIdx]
	cmdArgs := os.Args[cmdIdx+1:]

	// Setup logging
	if cfg.verbose == 0 {
		log.SetOutput(os.Stderr)
		log.SetFlags(0)
	} else {
		log.SetOutput(os.Stderr)
		log.SetFlags(log.LstdFlags)
	}

	// Execute command
	switch command {
	case "validate-id":
		runValidateID(cmdArgs)
	case "parse-id":
		runParseID(cmdArgs)
	case "match-id-pattern":
		runMatchIDPattern(cmdArgs)
	case "uuid":
		runUUID(cmdArgs)
	case "validate-instance":
		runValidateInstance(cmdArgs, cfg)
	case "resolve-relationships":
		runResolveRelationships(cmdArgs, cfg)
	case "compatibility":
		runCompatibility(cmdArgs, cfg)
	case "cast":
		runCast(cmdArgs, cfg)
	case "query":
		runQuery(cmdArgs, cfg)
	case "attr":
		runAttr(cmdArgs, cfg)
	case "list":
		runList(cmdArgs, cfg)
	case "server":
		runServer(cmdArgs, cfg)
	case "openapi-spec":
		runOpenAPISpec(cmdArgs, cfg)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		fmt.Fprint(os.Stderr, usageText)
		os.Exit(1)
	}
}

func runValidateID(args []string) {
	fs := flag.NewFlagSet("validate-id", flag.ExitOnError)
	gtsID := fs.String("gts-id", "", "GTS ID to validate (required)")
	fs.Parse(args)

	if *gtsID == "" {
		fmt.Fprintln(os.Stderr, "Error: --gts-id is required")
		os.Exit(1)
	}

	result := gts.ValidateGtsID(*gtsID)
	outputJSON(result)
}

func runParseID(args []string) {
	fs := flag.NewFlagSet("parse-id", flag.ExitOnError)
	gtsID := fs.String("gts-id", "", "GTS ID to parse (required)")
	fs.Parse(args)

	if *gtsID == "" {
		fmt.Fprintln(os.Stderr, "Error: --gts-id is required")
		os.Exit(1)
	}

	result := gts.ParseGtsID(*gtsID)
	outputJSON(result)
}

func runMatchIDPattern(args []string) {
	fs := flag.NewFlagSet("match-id-pattern", flag.ExitOnError)
	pattern := fs.String("pattern", "", "Pattern to match against (required)")
	candidate := fs.String("candidate", "", "Candidate GTS ID (required)")
	fs.Parse(args)

	if *pattern == "" || *candidate == "" {
		fmt.Fprintln(os.Stderr, "Error: --pattern and --candidate are required")
		os.Exit(1)
	}

	result := gts.MatchIDPattern(*candidate, *pattern)
	outputJSON(result)
}

func runUUID(args []string) {
	fs := flag.NewFlagSet("uuid", flag.ExitOnError)
	gtsID := fs.String("gts-id", "", "GTS ID (required)")
	scope := fs.String("scope", "major", "UUID scope (major or full)")
	fs.Parse(args)

	if *gtsID == "" {
		fmt.Fprintln(os.Stderr, "Error: --gts-id is required")
		os.Exit(1)
	}

	result := gts.IDToUUID(*gtsID)

	// Note: scope parameter is currently not used in the Go implementation
	// It's included for API compatibility with Python version
	_ = scope

	outputJSON(result)
}

func runValidateInstance(args []string, cfg *config) {
	fs := flag.NewFlagSet("validate-instance", flag.ExitOnError)
	gtsID := fs.String("gts-id", "", "GTS ID of the object (required)")
	fs.Parse(args)

	if *gtsID == "" {
		fmt.Fprintln(os.Stderr, "Error: --gts-id is required")
		os.Exit(1)
	}

	store := createStore(cfg)
	result := store.ValidateInstance(*gtsID)
	outputJSON(result)
}

func runResolveRelationships(args []string, cfg *config) {
	fs := flag.NewFlagSet("resolve-relationships", flag.ExitOnError)
	gtsID := fs.String("gts-id", "", "GTS ID of the entity (required)")
	fs.Parse(args)

	if *gtsID == "" {
		fmt.Fprintln(os.Stderr, "Error: --gts-id is required")
		os.Exit(1)
	}

	store := createStore(cfg)
	result := store.BuildSchemaGraph(*gtsID)
	outputJSON(result)
}

func runCompatibility(args []string, cfg *config) {
	fs := flag.NewFlagSet("compatibility", flag.ExitOnError)
	oldSchemaID := fs.String("old-schema-id", "", "GTS ID of old schema (required)")
	newSchemaID := fs.String("new-schema-id", "", "GTS ID of new schema (required)")
	fs.Parse(args)

	if *oldSchemaID == "" || *newSchemaID == "" {
		fmt.Fprintln(os.Stderr, "Error: --old-schema-id and --new-schema-id are required")
		os.Exit(1)
	}

	store := createStore(cfg)
	result := store.CheckCompatibility(*oldSchemaID, *newSchemaID)
	outputJSON(result)
}

func runCast(args []string, cfg *config) {
	fs := flag.NewFlagSet("cast", flag.ExitOnError)
	fromID := fs.String("from-id", "", "GTS ID of instance or schema to be casted (required)")
	toSchemaID := fs.String("to-schema-id", "", "GTS ID of target schema (required)")
	fs.Parse(args)

	if *fromID == "" || *toSchemaID == "" {
		fmt.Fprintln(os.Stderr, "Error: --from-id and --to-schema-id are required")
		os.Exit(1)
	}

	store := createStore(cfg)
	result, err := store.Cast(*fromID, *toSchemaID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func runQuery(args []string, cfg *config) {
	fs := flag.NewFlagSet("query", flag.ExitOnError)
	expr := fs.String("expr", "", "Query expression (required)")
	limit := fs.Int("limit", 100, "Maximum number of entities to return")
	fs.Parse(args)

	if *expr == "" {
		fmt.Fprintln(os.Stderr, "Error: --expr is required")
		os.Exit(1)
	}

	store := createStore(cfg)
	result := store.Query(*expr, *limit)
	outputJSON(result)
}

func runAttr(args []string, cfg *config) {
	fs := flag.NewFlagSet("attr", flag.ExitOnError)
	gtsWithPath := fs.String("gts-with-path", "", "GTS ID with attribute path (required)")
	fs.Parse(args)

	if *gtsWithPath == "" {
		fmt.Fprintln(os.Stderr, "Error: --gts-with-path is required")
		os.Exit(1)
	}

	store := createStore(cfg)
	result := store.GetAttribute(*gtsWithPath)
	outputJSON(result)
}

func runList(args []string, cfg *config) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	limit := fs.Int("limit", 100, "Maximum number of entities to return")
	fs.Parse(args)

	store := createStore(cfg)
	result := store.List(*limit)
	outputJSON(result)
}

func runServer(args []string, cfg *config) {
	fs := flag.NewFlagSet("server", flag.ExitOnError)
	host := fs.String("host", "127.0.0.1", "Host to bind to")
	port := fs.Int("port", 8000, "Port to listen on")
	fs.Parse(args)

	store := createStore(cfg)

	fmt.Printf("starting the server @ http://%s:%d\n", *host, *port)
	if cfg.verbose == 0 {
		fmt.Println("use --verbose to see server logs")
	}

	srv := server.NewServer(store, *host, *port, cfg.verbose)
	log.Fatal(srv.Start())
}

func runOpenAPISpec(args []string, cfg *config) {
	fs := flag.NewFlagSet("openapi-spec", flag.ExitOnError)
	out := fs.String("out", "", "Destination file path for OpenAPI spec JSON (required)")
	host := fs.String("host", "127.0.0.1", "Server host")
	port := fs.Int("port", 8000, "Server port")
	fs.Parse(args)

	if *out == "" {
		fmt.Fprintln(os.Stderr, "Error: --out is required")
		os.Exit(1)
	}

	store := createStore(cfg)
	srv := server.NewServer(store, *host, *port, cfg.verbose)

	// Get OpenAPI spec
	spec := srv.GetOpenAPISpec()

	// Write to file
	file, err := os.Create(*out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(spec); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing OpenAPI spec: %v\n", err)
		os.Exit(1)
	}

	// Output result to stdout
	result := map[string]any{
		"ok":  true,
		"out": *out,
	}
	outputJSON(result)
}

func createStore(cfg *config) *gts.GtsStore {
	var reader gts.GtsReader

	if cfg.path != "" {
		paths := strings.Split(cfg.path, ",")
		for i := range paths {
			paths[i] = strings.TrimSpace(paths[i])
		}

		var gtsConfig *gts.GtsConfig
		if cfg.config != "" {
			gtsConfig = loadConfig(cfg.config)
		}

		reader = gts.NewGtsFileReader(paths, gtsConfig)
	}

	return gts.NewGtsStore(reader)
}

func loadConfig(configPath string) *gts.GtsConfig {
	file, err := os.Open(configPath)
	if err != nil {
		log.Printf("Warning: Could not open config file: %v", err)
		return gts.DefaultGtsConfig()
	}
	defer file.Close()

	var configData struct {
		EntityIDFields []string `json:"entity_id_fields"`
		SchemaIDFields []string `json:"schema_id_fields"`
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&configData); err != nil {
		log.Printf("Warning: Could not parse config file: %v", err)
		return gts.DefaultGtsConfig()
	}

	return &gts.GtsConfig{
		EntityIDFields: configData.EntityIDFields,
		SchemaIDFields: configData.SchemaIDFields,
	}
}

func outputJSON(v any) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}
