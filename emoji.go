// Package emoji provides the shared emoji registry for fizz.
// It embeds emojis.json and exposes a stable name→ID mapping.
package emoji

import (
	_ "embed"
	"encoding/json"
)

//go:embed emojis.json
var emojisJSON []byte

var (
	idSet    map[uint16]bool
	idToName map[uint16]string
	nameToID map[string]uint16
)

func init() {
	if err := json.Unmarshal(emojisJSON, &nameToID); err != nil {
		panic("emoji: invalid emojis.json: " + err.Error())
	}

	idSet = make(map[uint16]bool, len(nameToID))
	idToName = make(map[uint16]string, len(nameToID))

	for name, id := range nameToID {
		idSet[id] = true
		idToName[id] = name
	}
}

// Count returns the total number of known emojis.
func Count() uint16 { return uint16(len(nameToID)) }

// Valid reports whether id corresponds to a known emoji.
func Valid(id uint16) bool { return idSet[id] }

// NameByID returns the emoji name for the given ID, or "" if unknown.
func NameByID(id uint16) string { return idToName[id] }

// ID returns the hash-based ID for the named emoji and true,
// or 0 and false if the name is not known.
func ID(name string) (uint16, bool) {
	id, ok := nameToID[name]
	return id, ok
}

// All returns a copy of the name→ID mapping for every known emoji.
func All() map[string]uint16 {
	out := make(map[string]uint16, len(nameToID))
	for k, v := range nameToID {
		out[k] = v
	}
	return out
}
