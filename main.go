package main

import (
    "flag"
    "fmt"
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
    "math/rand"
    "time"

    "github.com/StephaneBunel/bresenham"
)


func main() {
    if err := Run(os.Args[1:]); err != nil {
        slog.Error("fatal", "err", err)
        os.Exit(1)
    }
}

// Run executes the application logic using the provided args (typically os.Args[1:]).
// It returns an error instead of exiting, which makes it testable.
func Run(args []string) error {
    fs := flag.NewFlagSet("ttr", flag.ContinueOnError)
    show := fs.Bool("show", false, "show the game on the map")
    mapPath := fs.String("map", "./pkg/data/USA/USA_map.jpg", "path to base map image")
    outImage := fs.String("out-image", "./pkg/data/USA/USA_map_out.jpg", "output image path")
    outGif := fs.String("out-gif", "./pkg/data/USA/USA_map_out.gif", "output gif path")
    frameDelay := fs.Int("frame-delay", 45, "frame delay for gif in 100ths of a second")
    mode := fs.String("mode", "tickets", "player mode: random|tickets")
    routesPath := fs.String("routes", "./pkg/data/USA/routes.csv", "path to routes CSV")
    citiesPath := fs.String("cities", "./pkg/data/USA/cities.csv", "path to cities CSV")
    ticketsPath := fs.String("tickets", "./pkg/data/USA/tickets.csv", "path to tickets CSV")
    if err := fs.Parse(args); err != nil {
        return fmt.Errorf("parse flags: %w", err)
    }

    // load routes from provided path
    rf, err := os.Open(*routesPath)
    if err != nil {
        return fmt.Errorf("open routes: %w", err)
    }
    defer rf.Close()
    routes, err := data.RoutesFromReader(rf)
    if err != nil {
        return fmt.Errorf("load routes: %w", err)
    }

    // load cities from provided path
    cf, err := os.Open(*citiesPath)
    if err != nil {
        return fmt.Errorf("open cities: %w", err)
    }
    defer cf.Close()
    cities, err := data.CitiesFromReader(cf)
    if err != nil {
        return fmt.Errorf("load cities: %w", err)
    }

    // Create players according to selected mode
    type agent interface {
        Play() func(game.Board) (game.City, game.City)
        player.Scorer
    }
    var p1 agent
    var p2 agent
    switch *mode {
    case "random":
        // seed RNG for randomness
        rand.Seed(time.Now().UnixNano())
        p1 = player.NewRandom(1)
        p2 = player.NewRandom(2)
    case "tickets":
        // Load tickets from provided path and create ticket-aware players
        tf, err := os.Open(*ticketsPath)
        if err != nil {
            return fmt.Errorf("open tickets: %w", err)
        }
        defer tf.Close()
        tickets, err := data.TicketsFromReader(tf)
        if err != nil {
            return fmt.Errorf("load tickets: %w", err)
        }
        // Create players that receive 3 tickets each from the tickets pool
        p1 = player.NewWithTickets(1, game.GetTickets(3, &tickets))
        p2 = player.NewWithTickets(2, game.GetTickets(3, &tickets))
    default:
        return fmt.Errorf("unknown mode: %s", *mode)
    }

    var layer *image.NRGBA
    if *show {
        var err error
        layer, err = data.LoadMap(*mapPath)
        if err != nil {
            return fmt.Errorf("load map: %w", err)
        }
    }

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

    // Log scores
    slog.Info("end game", "Score P1", player.Score(p1), "Score P2", player.Score(p2))

    if !*show {
        return nil
    }
    if err := render.Map(layer, frames, *outImage, *outGif, *frameDelay); err != nil {
        return fmt.Errorf("render outputs: %w", err)
    }
    return nil
}
