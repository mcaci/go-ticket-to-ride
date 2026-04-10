package data

import (
	"strings"
	"testing"
)

func TestCitiesFromReader(t *testing.T) {
	csv := `X,Y,City
10.5,20.7,Alpha
0,0,Beta
`
	m, err := CitiesFromReader(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("CitiesFromReader error: %v", err)
	}
	if len(m) != 2 {
		t.Fatalf("expected 2 cities, got %d", len(m))
	}
	if c, ok := m["Alpha"]; !ok {
		t.Fatalf("missing Alpha")
	} else {
		if c.X != 10 || c.Y != 20 {
			t.Fatalf("unexpected Alpha coords: %+v", c)
		}
	}
	if c, ok := m["Beta"]; !ok {
		t.Fatalf("missing Beta")
	} else {
		if c.X != 0 || c.Y != 0 {
			t.Fatalf("unexpected Beta coords: %+v", c)
		}
	}
}
