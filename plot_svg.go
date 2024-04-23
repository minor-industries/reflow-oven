package main

import (
	"github.com/pkg/errors"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func plot_svg(
	profile Schedule,
	tcs []*Thermocouple,
	data [][]Point,
) error {
	p := plot.New()

	p.Title.Text = "Reflow Oven Temperature"
	p.X.Label.Text = "t"
	p.Y.Label.Text = "Temp"

	var vs []interface{}
	{
		vs = append(vs, "profile")
		var pts plotter.XYs
		for _, d := range profile {
			pts = append(pts, plotter.XY{
				d.T().Seconds(),
				d.Val,
			})
		}
		vs = append(vs, pts)
	}

	for i, tc := range tcs {
		vs = append(vs, tc.Description)
		pts := plotter.XYs{}
		for _, d := range data[i] {
			pts = append(pts, plotter.XY{
				d.T().Seconds(),
				d.Val,
			})
		}
		vs = append(vs, pts)
	}

	err := plotutil.AddLines(p, vs...)
	if err != nil {
		return errors.Wrap(err, "add lines")
	}

	err = p.Save(16*vg.Inch, 4*vg.Inch, "cook.svg")
	if err != nil {
		return errors.Wrap(err, "save svg")
	}

	return nil
}
