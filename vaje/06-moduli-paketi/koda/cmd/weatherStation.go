package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/davors/weatherstation/weather"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "weatherStation",
		Usage: "Weather station displays Temperature, Humidity, and Pressure",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "timeout",
				Usage: "Timeout for reading sensor in seconds",
				Value: 5,
			},
			&cli.IntFlag{
				Name:  "pollTime",
				Usage: "Sensor refresh interval in ms",
				Value: 500,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			timeoutSec := cmd.Int("timeout")
			pollTime := cmd.Int("pollTime")
			return runTUI(pollTime, timeoutSec)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func runTUI(pollTime int, timeoutSec int) error {
	station := weather.NewStation(time.Duration(pollTime)*time.Millisecond, time.Duration(timeoutSec)*time.Second)

	app := tview.NewApplication()

	mainView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetChangedFunc(func() {
			app.Draw()
		})
	mainView.SetBorder(true).SetTitle(" Vremenska postaja ")

	footer := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[gray]Pritisnite [white]q[gray] za izhod").
		SetTextAlign(tview.AlignLeft)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainView, 0, 1, true).
		AddItem(footer, 3, 0, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q', 'Q':
			app.Stop()
			return nil
		}
		return event
	})

	go func() {
		defer station.Stop()
		meritve := map[string]float32{"Temperatura": 0.0, "Vlaga": 0.0, "Tlak": 0.0}
		for {
			if data, err := station.GetData(); err != nil {
				app.QueueUpdateDraw(func() {
					footer.SetText("[red]" + err.Error() + "\n[gray]Pritisnite [white]q[gray] za izhod")
				})
				return
			} else {
				switch data.MType {
				case "Temperature":
					meritve["Temperature"] = data.Value
				case "Humidity":
					meritve["Humidity"] = data.Value
				case "Pressure":
					meritve["Pressure"] = data.Value
				}
				text := fmt.Sprintf(
					"[yellow:]Temperatura: [white] %3.1f °C\n"+
						"[blue:]Vlažnost:    [white] %3.1f %%\n"+
						"[green:]Tlak:        [white] %3.1f mbar",
					meritve["Temperature"],
					meritve["Humidity"],
					meritve["Pressure"],
				)

				app.QueueUpdateDraw(func() {
					mainView.SetText(text)
				})
			}
		}
	}()

	return app.SetRoot(layout, true).Run()
}
