package actions

import (
	"github.com/kshmirko/rt3v1/rtcode"
	"github.com/urfave/cli"
)

func DoSSRT(ctx *cli.Context) error {
	fname := ctx.String("fname")
	sza := ctx.Float64("sza")
	wl := ctx.Float64("wl")

	p := rtcode.NewSSRT(fname, sza, wl)

	p.Dump()
	return nil
}

func DoRT3(ctx *cli.Context) error {
	r0 := ctx.Float64("r0")
	r1 := ctx.Float64("r1")
	gamma := ctx.Float64("gamma")
	sza := ctx.Float64("sza")
	wl := ctx.Float64("wl")
	mre := ctx.Float64("mre")
	mim := ctx.Float64("mim")
	midx := complex64(complex(mre, mim))
	galbedo := ctx.Float64("galbedo")
	display := ctx.Bool("display")
	nlays := ctx.Int("nlays")
	aot := ctx.Float64("taua")

	v := rtcode.New()
	v.SetSizeDistrib(r0, r1, gamma, 101)
	v.SetSza(sza)
	v.SetWl(wl)
	v.SetMidx(midx)
	v.SetGalbedo(galbedo)
	v.SetNlays(nlays)
	v.SetTaua(aot)
	v.DoCalc()
	a, _ := v.UnmarshalData()
	_, _, _ = a.DumpDownwardRadiation(display)

	rtcode.CleanUp()
	return nil
}
