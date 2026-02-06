package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/djsam/mcp-gateway-envoy/internal/config"
	"github.com/djsam/mcp-gateway-envoy/internal/controller"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		printUsage()
		return nil
	}

	switch args[0] {
	case "init":
		return runInit(args[1:])
	case "validate":
		return runValidate(args[1:])
	case "plan":
		return runPlan(args[1:])
	case "help", "-h", "--help":
		printUsage()
		return nil
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func runInit(args []string) error {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	output := fs.String("output", "gateway.yaml", "path to output config file")
	force := fs.Bool("force", false, "overwrite output if file already exists")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if !*force {
		if _, err := os.Stat(*output); err == nil {
			return fmt.Errorf("output file %q already exists (use --force to overwrite)", *output)
		}
	}

	dir := filepath.Dir(*output)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create output directory: %w", err)
		}
	}

	if err := os.WriteFile(*output, []byte(config.DefaultTemplateYAML), 0o644); err != nil {
		return fmt.Errorf("write template config: %w", err)
	}

	fmt.Printf("created %s\n", *output)
	fmt.Println("next: gateway validate --file", *output)
	return nil
}

func runValidate(args []string) error {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	file := fs.String("file", "gateway.yaml", "path to config file")
	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg, err := config.LoadFile(*file)
	if err != nil {
		return err
	}

	fmt.Printf("config is valid: %s\n", *file)
	fmt.Printf("gateway=%s servers=%d routes=%d secureDefault=%t\n",
		cfg.Gateway.Name, len(cfg.Servers), len(cfg.Routes), cfg.Auth.RequireAuth)
	return nil
}

func runPlan(args []string) error {
	fs := flag.NewFlagSet("plan", flag.ContinueOnError)
	file := fs.String("file", "gateway.yaml", "path to config file")
	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg, err := config.LoadFile(*file)
	if err != nil {
		return err
	}

	resources, err := controller.BuildResources(cfg)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(resources, "", "  ")
	if err != nil {
		return err
	}
	if len(resources) == 0 {
		return errors.New("no resources produced")
	}

	fmt.Println(string(b))
	return nil
}

func printUsage() {
	fmt.Print(`mcp-gateway-envoy

Usage:
  gateway init [--output gateway.yaml] [--force]
  gateway validate [--file gateway.yaml]
  gateway plan [--file gateway.yaml]
`)
}
