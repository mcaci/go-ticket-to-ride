package data

import (
	"go-ticket-to-ride/pkg/game"
	"strings"
	"testing"
)

func TestRoutesFromReader(t *testing.T) {
	csv := `X,Y,Distance,Color
A,B,3,B
C,D,5,X
`
	b, err := RoutesFromReader(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("RoutesFromReader error: %v", err)
	}
	edges := b.Edges()
	if len(edges) < 2 {
		t.Fatalf("expected at least 2 edge entries, got %d", len(edges))
	}
	// find edge A->B
	var foundAB, foundCD bool
	for _, e := range edges {
		if e.X.E == game.City("A") && e.Y.E == game.City("B") {
			prop := e.P.(*game.TrainLineProperty)
			if prop.Distance != 3 || prop.Color != game.Blue {
				t.Fatalf("unexpected A->B prop: %+v", prop)
			}
			foundAB = true
		}
		if e.X.E == game.City("C") && e.Y.E == game.City("D") {
			prop := e.P.(*game.TrainLineProperty)
			if prop.Distance != 5 || prop.Color != game.All {
				t.Fatalf("unexpected C->D prop: %+v", prop)
			}
			foundCD = true
		}
	}
	if !foundAB || !foundCD {
		t.Fatalf("expected both edges found")
	}
}
