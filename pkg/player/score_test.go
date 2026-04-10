package player

import (
	"testing"

	"go-ticket-to-ride/pkg/game"

	"github.com/mcaci/graphgo/graph"
)

// fakeScorer implements Scorer for testing
type fakeScorer struct {
	lines []*graph.Edge[game.City]
	ticks []game.Ticket
}

func (f *fakeScorer) TrainLines() []*graph.Edge[game.City] { return f.lines }
func (f *fakeScorer) Tickets() []game.Ticket               { return f.ticks }

func TestScore_Empty(t *testing.T) {
	f := &fakeScorer{}
	if Score(f) != 0 {
		t.Fatalf("expected 0, got %d", Score(f))
	}
}

func TestScore_TrainsAndTickets(t *testing.T) {
	// create vertices
	x := &graph.Vertex[game.City]{E: "A"}
	y := &graph.Vertex[game.City]{E: "B"}
	// create an edge with TrainLineProperty Distance=3 (score 4)
	prop := &game.TrainLineProperty{Distance: 3}
	edge := &graph.Edge[game.City]{X: x, Y: y, P: prop}

	// create tickets: one done valued 10, one not done valued 5
	tickets := []game.Ticket{{Value: 10, Done: true}, {Value: 5, Done: false}}

	f := &fakeScorer{lines: []*graph.Edge[game.City]{edge}, ticks: tickets}
	expected := 4 + 10 // train line score for distance 3 is 4
	if got := Score(f); got != expected {
		t.Fatalf("expected %d, got %d", expected, got)
	}
}
