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

type MapCity map[game.City]Coord

type Coord struct {
	X int
	Y int
}

// CitiesFromReader parses city coordinates from an io.Reader (CSV with header). Coordinates may be floats and are truncated to int.
func CitiesFromReader(r io.Reader) (MapCity, error) {
	rdr := csv.NewReader(r)
	if _, err := rdr.Read(); err != nil {
		if err == io.EOF {
			return make(MapCity), nil
		}
		return nil, fmt.Errorf("read header: %w", err)
	}
	m := make(MapCity)
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
		xS := strings.TrimSpace(row[0])
		yS := strings.TrimSpace(row[1])
		city := strings.TrimSpace(row[2])
		xf, err := strconv.ParseFloat(xS, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid x coord %q: %w", xS, err)
		}
		yf, err := strconv.ParseFloat(yS, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid y coord %q: %w", yS, err)
		}
		m[game.City(city)] = Coord{X: int(xf), Y: int(yf)}
	}
	return m, nil
}

// Cities opens the default cities CSV and parses it.
func Cities() (MapCity, error) {
	f, err := os.Open("./pkg/data/USA/cities.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return CitiesFromReader(f)
}
