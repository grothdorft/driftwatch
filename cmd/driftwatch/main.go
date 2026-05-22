package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/driftwatch/internal/loader"
)

func main() {
	dir := flag.String("dir", ".", "directory containing source-of-truth YAML manifests")
	flag.Parse()

	log.SetFlags(0)
	log.SetPrefix("driftwatch: ")

	manifests, err := loader.LoadFromDir(*dir)
	if err != nil {
		log.Fatalf("failed to load manifests from %s: %v", *dir, err)
	}

	if len(manifests) == 0 {
		log.Printf("no YAML manifests found in %s", *dir)
		os.Exit(0)
	}

	fmt.Printf("Loaded %d manifest(s) from %s:\n", len(manifests), *dir)
	for _, m := range manifests {
		key, err := loader.ManifestKey(m)
		if err != nil {
			log.Printf("warning: skipping manifest %s: %v", m.Source, err)
			continue
		}
		fmt.Printf("  [%s] -> %s\n", key, m.Source)
	}
}
