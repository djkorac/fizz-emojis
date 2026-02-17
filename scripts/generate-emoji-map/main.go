// Command generate-emoji-map reads newline-delimited emoji names from stdin,
// assigns stable uint16 IDs using an append-only registry, and writes
// a JSON object {"name": id, ...} to the path given by -o.
//
// On each run the existing output file is loaded first. Names already present
// keep their IDs unchanged (preserving wire compatibility with older clients).
// New names are sorted and assigned the lowest available uint16 IDs not
// already in use.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
)

func main() {
	outPath := flag.String("o", "emojis.json", "output JSON file path")
	flag.Parse()

	var names []string
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line := sc.Text()
		if line != "" {
			names = append(names, line)
		}
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading stdin: %v\n", err)
		os.Exit(1)
	}

	if len(names) == 0 {
		fmt.Fprintln(os.Stderr, "no emoji names provided")
		os.Exit(1)
	}

	// Load existing assignments to preserve wire-compatible IDs.
	result := make(map[string]uint16, len(names))
	if data, err := os.ReadFile(*outPath); err == nil {
		if err := json.Unmarshal(data, &result); err != nil {
			fmt.Fprintf(os.Stderr, "error parsing existing %s: %v\n", *outPath, err)
			os.Exit(1)
		}
	}

	// Build the set of already-used IDs.
	usedIDs := make(map[uint16]bool, len(result))
	for _, id := range result {
		usedIDs[id] = true
	}

	// Sort new names so ID assignment is deterministic regardless of input order.
	sort.Strings(names)

	// Assign the lowest available IDs (starting at 1) to names not yet in the
	// registry.  ID 0 is reserved as a sentinel for "no emoji".
	nextID := uint16(1)
	added := 0
	for _, name := range names {
		if _, ok := result[name]; ok {
			continue
		}
		for usedIDs[nextID] {
			nextID++
			if nextID == 0 {
				fmt.Fprintln(os.Stderr, "error: uint16 ID space exhausted")
				os.Exit(1)
			}
		}
		result[name] = nextID
		usedIDs[nextID] = true
		nextID++
		added++
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	data = append(data, '\n')

	if err := os.WriteFile(*outPath, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing %s: %v\n", *outPath, err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "wrote %d emojis (%d new) to %s\n", len(result), added, *outPath)
}
