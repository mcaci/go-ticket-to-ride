package main

import (
	"flag"
	"go-ticket-to-ride/pkg/data"
	"go-ticket-to-ride/pkg/game"
	"go-ticket-to-ride/pkg/player"
	"go-ticket-to-ride/pkg/render"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"log/slog"
	"os"

	"github.com/StephaneBunel/bresenham"
)

func main() {
	show := flag.Bool("show", false, "show the game on the map")
	mapPath := flag.String("map", "./pkg/data/USA/USA_map.jpg", "path to base map image")
	outImage := flag.String("out-image", "./pkg/data/USA/USA_map_out.jpg", "output image path")
	outGif := flag.String("out-gif", "./pkg/data/USA/USA_map_out.gif", "output gif path")
	frameDelay := flag.Int("frame-delay", 45, "frame delay for gif in 100ths of a second")
	mode := flag.String("mode", "tickets", "player mode: random|tickets")
	routesPath := flag.String("routes", "./pkg/data/USA/routes.csv", "path to routes CSV")
	citiesPath := flag.String("cities", "./pkg/data/USA/cities.csv", "path to cities CSV")
	ticketsPath := flag.String("tickets", "./pkg/data/USA/tickets.csv", "path to tickets CSV")
	flag.Parse()

	// load routes from provided path
	rf, err := os.Open(*routesPath)
	if err != nil {
		slog.Error("failed to open routes file", "err", err)
		os.Exit(1)
	}
	defer rf.Close()
	routes, err := data.RoutesFromReader(rf)
	if err != nil {
		slog.Error("failed to load routes", "err", err)
		os.Exit(1)
	}

	// load cities from provided path
	cf, err := os.Open(*citiesPath)
	if err != nil {
		slog.Error("failed to open cities file", "err", err)
		os.Exit(1)
	}
	defer cf.Close()
	cities, err := data.CitiesFromReader(cf)
	if err != nil {
		slog.Error("failed to load cities", "err", err)
		os.Exit(1)
	}
	_ = ticketsPath // currently unused but kept for future use

	var layer *image.NRGBA
	if *show {
		var err error
		layer, err = data.LoadMap(*mapPath)
		if err != nil {
			slog.Error("cannot load base map image", "err", err)
			os.Exit(1)
		}
	}

	// Create players according to selected mode
	type agent interface {
		Play() func(game.Board) (game.City, game.City)
		player.Scorer
	}
	var p1, p2 agent
	switch *mode {
	case "random":
		p1 = player.NewRandom(1)
		p2 = player.NewRandom(2)
	case "tickets":
		// Load tickets from provided path and create ticket-aware players
		tf, err := os.Open(*ticketsPath)
		if err != nil {
			slog.Error("failed to open tickets file", "err", err)
			os.Exit(1)
		}
		defer tf.Close()
		tickets, err := data.TicketsFromReader(tf)
		if err != nil {
			slog.Error("failed to load tickets", "err", err)
			os.Exit(1)
		}
		// Create players that receive 3 tickets each from the tickets pool
		p1 = player.NewWithTickets(1, game.GetTickets(3, &tickets))
		p2 = player.NewWithTickets(2, game.GetTickets(3, &tickets))
	default:
		slog.Error("unknown mode", "mode", *mode)
		os.Exit(1)
	}

	// For debugging purposes, we can use the fix the tickets
	// p1, p2 := player.NewTAPl(1, []game.Ticket{tickets[3], tickets[26], tickets[21]}), player.NewTAPl(2, []game.Ticket{tickets[12], tickets[2], tickets[22]})
	var coin bool
	var frames []*image.Paletted
	for game.FreeRoutesAvailable(routes) {
		var play func(game.Board) (game.City, game.City)
		coin = !coin
		switch coin {
		case true:
			play = p1.Play()
		case false:
			play = p2.Play()
		}
		a, b := play(routes)

		if !*show {
			continue
		}
		var c color.Color
		switch coin {
		case true:
			c = color.RGBA{R: 0, G: 0, B: 255, A: 255}
		case false:
			c = color.RGBA{R: 255, G: 0, B: 255, A: 255}
		}
		bresenham.DrawLine(layer, cities[a].X, cities[a].Y, cities[b].X, cities[b].Y, c)
		p := image.NewPaletted(layer.Bounds(), palette.Plan9)
		draw.Draw(p, p.Bounds(), layer, image.Point{}, draw.Over)
		frames = append(frames, p)
	}

	slog.Info("end game", "Score P1", player.Score(p1), "Score P2", player.Score(p2))

	if !*show {
		os.Exit(0)
	}
	if err := render.Map(layer, frames, *outImage, *outGif, *frameDelay); err != nil {
		slog.Error("failed to render/save outputs", "err", err)
		os.Exit(1)
	}
}
