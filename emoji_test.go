package emoji

import "testing"

func TestAllIDsUnique(t *testing.T) {
	seen := make(map[uint16]string, len(idToName))
	for id, name := range idToName {
		if prev, ok := seen[id]; ok {
			t.Fatalf("duplicate ID %d: %q and %q", id, prev, name)
		}
		seen[id] = name
	}
	if len(seen) != int(Count()) {
		t.Fatalf("expected %d unique IDs, got %d", Count(), len(seen))
	}
}

func TestValidAcceptsAllLoaded(t *testing.T) {
	for name, id := range nameToID {
		if !Valid(id) {
			t.Errorf("Valid(%d) = false for loaded emoji %q", id, name)
		}
	}
}

func TestValidRejectsUnknown(t *testing.T) {
	// Find an ID that is not in the set.
	var unknown uint16
	for id := uint16(0); ; id++ {
		if !idSet[id] {
			unknown = id
			break
		}
	}
	if Valid(unknown) {
		t.Errorf("Valid(%d) = true for unknown ID", unknown)
	}
}

func TestNameByIDRoundTrips(t *testing.T) {
	for name, id := range nameToID {
		got := NameByID(id)
		if got != name {
			t.Errorf("NameByID(%d) = %q, want %q", id, got, name)
		}
		gotID, ok := ID(got)
		if !ok || gotID != id {
			t.Errorf("ID(%q) = (%d, %v), want (%d, true)", got, gotID, ok, id)
		}
	}
}

func TestIDUnknownName(t *testing.T) {
	_, ok := ID("not_a_real_emoji_name_xyz")
	if ok {
		t.Error("ID returned ok=true for unknown name")
	}
}

func TestAllReturnsCorrectCount(t *testing.T) {
	all := All()
	if len(all) != int(Count()) {
		t.Errorf("All() has %d entries, Count() = %d", len(all), Count())
	}
}
