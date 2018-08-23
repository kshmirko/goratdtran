package main

import (
	_ "fmt"
	"github.com/kshmirko/rt3v1/actions"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		cli.Command{
			Name:     "ssrt",
			Category: "Радиационный код",
			Usage:    "Расчет переноса солнечной радиации в приближении однократного расеяния",
			Action:   actions.DoSSRT,
			Flags: []cli.Flag{
				cli.Float64Flag{
					Name:  "sza",
					Usage: "Solar zenith angle",
					Value: 10.0,
				},
			},
		},
		cli.Command{
			Name:     "rt3",
			Category: "Радиационный код",
			Usage:    "Расчет переноса солнечной радиации по модели polradtran, расчет компонент I и Q",
			Action:   actions.DoRT3,
			Flags: []cli.Flag{
				cli.Float64Flag{
					Name:  "sza",
					Usage: "Solar zenith angle",
					Value: 10.0,
				},
				cli.Float64Flag{
					Name:  "r0",
					Usage: "Left boundary of radius range, um",
					Value: 0.1,
				},
				cli.Float64Flag{
					Name:  "r1",
					Usage: "Right boundary of radius range, um",
					Value: 1.0,
				},
				cli.Float64Flag{
					Name:  "gamma",
					Usage: "Exponent decay of particles concentration < 0",
					Value: -3.5,
				},
				cli.Float64Flag{
					Name:  "wl",
					Usage: "Wavelength, um",
					Value: 0.750,
				},
				cli.Float64Flag{
					Name:  "mre",
					Usage: "Real part of refractive index",
					Value: 1.5,
				},
				cli.Float64Flag{
					Name:  "mim",
					Usage: "Imaginary part of refractive index",
					Value: 0.0,
				},
				cli.Float64Flag{
					Name:  "galbedo",
					Usage: "Ground reflectance",
					Value: 0.0,
				},
				cli.IntFlag{
					Name:  "nlays",
					Usage: "Number of atmosphere layers (each layer 1 km width)",
					Value: 40,
				},
				cli.BoolFlag{
					Name:  "display",
					Usage: "Display intensities on the screen",
				},
			},
		},
	}
	app.Action = nil
	app.Version = "0.1"
	app.Usage = "Программа для расчета нисходящей солнечной радиации в различных атмосферных условиях"
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	return

	//fmt.Println(mu, I, Q)
}
