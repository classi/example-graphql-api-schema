package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	os.Exit(doMain())
}

var (
	rootDir string
	outPath string
)

func doMain() int {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	return 0
}

func run() error {
	flag.Parse()
	schemaPaths, err := filepath.Glob(filepath.Join(rootDir, "*.gql"))
	if err != nil {
		return fmt.Errorf("filepath.Glob: %w", err)
	}
	log.Printf("%d files found", len(schemaPaths))
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(schemaPaths); err != nil {
		return fmt.Errorf("json.Encoder.Encode: %w", err)
	}
	if err := f.Sync(); err != nil {
		return fmt.Errorf("os.File.Sync: %w", err)
	}
	return nil
}

func init() {
	flag.StringVar(&rootDir, "root", "./schemata", "root directory to search schemata")
	flag.StringVar(&outPath, "out", "./schema-files.json", "output file path")
}
