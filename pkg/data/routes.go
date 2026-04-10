package data

import (
	"encoding/csv"
	"fmt"
	"go-ticket-to-ride/pkg/game"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/mcaci/graphgo/graph"
)

var colorMap = map[string]game.Color{
	"X": game.All,
	"B": game.Blue,
	"R": game.Red,
	"G": game.Green,
	"Y": game.Yellow,
	"W": game.White,
	"P": game.Pink,
	"O": game.Orange,
	"K": game.Black,
}

// RoutesFromReader parses CSV from r and returns a game.Board. Expected columns: X,Y,Distance,Color
func RoutesFromReader(r io.Reader) (game.Board, error) {
	rdr := csv.NewReader(r)
	// consume header
	if _, err := rdr.Read(); err != nil {
		if err == io.EOF {
			return graph.New[game.City](graph.AdjacencyListType, false), nil
		}
		return nil, fmt.Errorf("read header: %w", err)
	}
	b := graph.New[game.City](graph.AdjacencyListType, false)
	// vertex reuse map to ensure same *Vertex used for identical city
	verts := make(map[game.City]*graph.Vertex[game.City])
	for {
		row, err := rdr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read: %w", err)
		}
		if len(row) < 4 {
			return nil, fmt.Errorf("invalid row, expected >=4 columns: %v", row)
		}
		xName := game.City(strings.TrimSpace(row[0]))
		yName := game.City(strings.TrimSpace(row[1]))
		dStr := strings.TrimSpace(row[2])
		cl := strings.TrimSpace(row[3])
		d, err := strconv.Atoi(dStr)
		if err != nil {
			return nil, fmt.Errorf("invalid distance %q: %w", dStr, err)
		}
		color, ok := colorMap[cl]
		if !ok {
			return nil, fmt.Errorf("unknown color %q", cl)
		}
		X, ok := verts[xName]
		if !ok {
			v := &graph.Vertex[game.City]{E: xName}
			verts[xName] = v
			X = v
			b.AddVertex(X)
		}
		Y, ok := verts[yName]
		if !ok {
			v := &graph.Vertex[game.City]{E: yName}
			verts[yName] = v
			Y = v
			b.AddVertex(Y)
		}
		e := game.TrainLine{X: X, Y: Y, P: &game.TrainLineProperty{Color: color, Distance: d}}
		b.AddEdge((*graph.Edge[game.City])(&e))
	}
	return b, nil
}

// Routes opens the default routes CSV and parses it.
func Routes() (game.Board, error) {
	f, err := os.Open("./pkg/data/USA/routes.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return RoutesFromReader(f)
}
