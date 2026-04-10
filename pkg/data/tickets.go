package data

import (
	"encoding/csv"
	"fmt"
	"go-ticket-to-ride/pkg/game"
	"io"
	"os"
	"strconv"
	"strings"
)

// TicketsFromReader parses tickets from any io.Reader (CSV with header).
func TicketsFromReader(r io.Reader) ([]game.Ticket, error) {
	rdr := csv.NewReader(r)
	// read header
	if _, err := rdr.Read(); err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, fmt.Errorf("read header: %w", err)
	}
	var tickets []game.Ticket
	for {
		row, err := rdr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read: %w", err)
		}
		if len(row) < 3 {
			return nil, fmt.Errorf("invalid row, expected >=3 columns: %v", row)
		}
		x := strings.TrimSpace(row[0])
		y := strings.TrimSpace(row[1])
		scoreStr := strings.TrimSpace(row[2])
		score, err := strconv.Atoi(scoreStr)
		if err != nil {
			return nil, fmt.Errorf("invalid score %q: %w", scoreStr, err)
		}
		tickets = append(tickets, game.Ticket{X: game.City(x), Y: game.City(y), Value: score})
	}
	return tickets, nil
}

// Tickets opens the default tickets CSV and parses it.
func Tickets() ([]game.Ticket, error) {
	f, err := os.Open("./pkg/data/USA/tickets.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return TicketsFromReader(f)
}
