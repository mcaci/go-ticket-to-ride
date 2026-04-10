package data

import (
	"go-ticket-to-ride/pkg/game"
	"strings"
	"testing"
)

func TestTicketsFromReader(t *testing.T) {
	csv := `X,Y,Value
A,B,5
C,D,10
`
	tickets, err := TicketsFromReader(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("TicketsFromReader error: %v", err)
	}
	if len(tickets) != 2 {
		t.Fatalf("expected 2 tickets, got %d", len(tickets))
	}
	if tickets[0].X != game.City("A") || tickets[0].Y != game.City("B") || tickets[0].Value != 5 {
		t.Fatalf("unexpected first ticket: %+v", tickets[0])
	}
	if tickets[1].X != game.City("C") || tickets[1].Y != game.City("D") || tickets[1].Value != 10 {
		t.Fatalf("unexpected second ticket: %+v", tickets[1])
	}
}
