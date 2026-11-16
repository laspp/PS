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
	mainView.SetBorder(true).SetTitle(" Weather station ")

	footer := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[gray]Press [white]q[gray] to quit").
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
		measurements := map[string]float32{"Temperature": 0.0, "Humidity": 0.0, "Pressure": 0.0}
		for {
			data, err := station.GetData()
			if err != nil {
				app.QueueUpdateDraw(func() {
					footer.SetText("[red]Error: " + err.Error() + "\n[gray]Press [white]q[gray] to quit")
				})
				return
			}
			measurements[data.MType] = data.Value
			text := fmt.Sprintf(
				"[yellow:]Temperature: [white] %3.1f °C\n"+
					"[blue:]Humidity:    [white] %3.1f %%\n"+
					"[green:]Pressure:    [white] %3.1f mbar",
				measurements["Temperature"],
				measurements["Humidity"],
				measurements["Pressure"],
			)

			app.QueueUpdateDraw(func() {
				mainView.SetText(text)
			})
		}
	}()

	return app.SetRoot(layout, true).Run()
}
