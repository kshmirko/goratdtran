package rtcode

/*
#cgo CFLAGS: -O3
#cgo LDFLAGS: -L. -lradtran -lmiev0 -lrt3 -L /usr/local/Cellar/gcc/8.2.0/lib/gcc/8/ -lgfortran

#include "rt3v1.h"
*/
import (
	"C"
)

import (
	//	"unsafe"
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
)

// Параметры расчета переноса радиации
type RT3Params struct {
	r0, r1, wl, gamma, dens, hpbl, taua, galbedo, sza float64
	npts, numazim, nmu, nlays                         int
	midx                                              complex64
	out_file                                          string
}

// Инициализирует структуру и возвращает ссылку на объект
func New() *RT3Params {
	return &RT3Params{
		r0:       0.1,
		r1:       1.0,
		npts:     101,
		wl:       0.750,
		midx:     1.4 - 0.00i,
		gamma:    -4.0,
		dens:     300.0,
		hpbl:     3000.0,
		taua:     0.1,
		numazim:  2,
		galbedo:  0.0,
		sza:      10.0,
		nmu:      32,
		out_file: "rt3.out",
	}
}

//Геттеры и Сеттеры для параметров

// R0
func (v *RT3Params) R0() float64 {
	return v.r0
}

func (v *RT3Params) SetR0(r0 float64) {
	v.r0 = r0
}

// R1
func (v *RT3Params) R1() float64 {
	return v.r1
}

func (v *RT3Params) SetR1(r1 float64) {
	v.r1 = r1
}

// npts
func (v *RT3Params) Npts() int {
	return v.npts
}

func (v *RT3Params) SetNpts(npts int) {
	v.npts = npts
}

func (v *RT3Params) SetRadiusRange(r0, r1 float64, npts int) {
	v.r0 = r0
	v.r1 = r1
	v.npts = npts
}

func (v *RT3Params) SetSizeDistrib(r0, r1, gamma float64, npts int) {
	v.r0 = r0
	v.r1 = r1
	v.gamma = gamma
	v.npts = npts
}

// Wl
func (v *RT3Params) Wl() float64 {
	return v.wl
}

func (v *RT3Params) SetWl(wl float64) {
	v.wl = wl
}

// midx
func (v *RT3Params) Midx() complex64 {
	return v.midx
}

func (v *RT3Params) SetMidx(midx complex64) {
	v.midx = midx
}

// gamma
func (v *RT3Params) Gamma() float64 {
	return v.gamma
}

func (v *RT3Params) SetGamma(gamma float64) {
	v.gamma = gamma
}

// Dens
func (v *RT3Params) Dens() float64 {
	return v.dens
}

func (v *RT3Params) SetDens(dens float64) {
	v.dens = dens
}

// Hpbl
func (v *RT3Params) Hpbl() float64 {
	return v.hpbl
}

func (v *RT3Params) SetHpbl(hpbl float64) {
	v.hpbl = hpbl
}

// taua
func (v *RT3Params) Taua() float64 {
	return v.taua
}

func (v *RT3Params) SetTaua(taua float64) {
	v.taua = taua
}

// numazim
func (v *RT3Params) Numazim() int {
	return v.numazim
}

func (v *RT3Params) SetNumazim(numazim int) {
	v.numazim = numazim
}

// galbedo
func (v *RT3Params) Galbedo() float64 {
	return v.galbedo
}

func (v *RT3Params) SetGalbedo(galbedo float64) {
	v.galbedo = galbedo
}

// Sza
func (v *RT3Params) Sza() float64 {
	return v.sza
}

func (v *RT3Params) SetSza(sza float64) {
	v.sza = sza
}

// nmu
func (v *RT3Params) Nmu() int {
	return v.nmu
}

func (v *RT3Params) SetNmu(nmu int) {
	v.nmu = nmu
}

func (v *RT3Params) Nlays() int {
	return v.nlays
}

func (v *RT3Params) SetNlays(nlays int) {
	v.nlays = nlays
}

// gamma
func (v *RT3Params) Outfile() string {
	return v.out_file
}

func (v *RT3Params) SetOutfile(outf string) {
	v.out_file = outf
}

// Выполняе расчет освещенности
func (v *RT3Params) DoCalc() {
	outf := C.CString(v.out_file)
	C.do_calc1(C.double(v.r0),
		C.double(v.r1),
		C.int(v.npts),
		C.double(v.wl),
		C.double(real(v.midx)),
		C.double(imag(v.midx)),
		C.double(v.gamma),
		C.double(v.dens),
		C.double(v.hpbl),
		C.double(v.taua),
		C.int(v.numazim),
		C.double(v.galbedo),
		C.double(v.sza),
		C.int(v.nmu),
		C.int(v.nlays),
		(*C.char)(outf),
	)
}

type ResultData struct {
	Z          float64
	Phi        *[]float64
	Mu         *[]float64
	Ival, Qval *[]float64
}

func (v *RT3Params) UnmarshalData() (*ResultData, error) {
	reader, err := os.Open(v.out_file)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	bufferedReader := bufio.NewReader(reader)
	nlines := v.nmu * v.numazim * 2

	eof := false
	var z float64
	phi := make([]float64, nlines)
	mu := make([]float64, nlines)
	Ival := make([]float64, nlines)
	Qval := make([]float64, nlines)

	for lineno := 0; lineno < 11 && !eof; lineno++ {
		_, err := bufferedReader.ReadString('\n')
		if err == io.EOF {
			err = nil
			eof = true
		} else if err != nil {
			return nil, err
		}
	}

	for lineno := 0; lineno < nlines && !eof; lineno++ {
		line, err := bufferedReader.ReadString('\n')
		if err == io.EOF {
			err = nil
			eof = true
		} else if err != nil {
			return nil, err
		}
		fmt.Sscanf(line, "%f %f %f %f %f", &z, &phi[lineno], &mu[lineno], &Ival[lineno], &Qval[lineno])
	}

	return &ResultData{
		Z:    z,
		Phi:  &phi,
		Mu:   &mu,
		Ival: &Ival,
		Qval: &Qval,
	}, nil
}

func (v *ResultData) DumpDownwardRadiation(display bool) (tmpMu, tmpI, tmpQ []float64) {
	nlen := len(*v.Ival)
	deg2rad := math.Pi / 180.0
	rad2deg := 1.0 / deg2rad
	nnlines := nlen / 2
	tmpMu = make([]float64, nnlines)
	tmpI = make([]float64, nnlines)
	tmpQ = make([]float64, nnlines)

	j := 0
	for i := 0; i < nlen; i++ {
		mu := (*v.Mu)[i]
		if mu >= 0 {
			mu = math.Acos(mu) * math.Cos((*v.Phi)[i]*deg2rad) * rad2deg
			//fmt.Printf("%12.4f%12.3e%12.3e\n", mu, (*v.Ival)[i], (*v.Qval)[i])
			tmpMu[j] = mu
			tmpI[j] = (*v.Ival)[i]
			tmpQ[j] = (*v.Qval)[i]
			j++
		}
	}

	// Reverse second half of the array
	for i, j := nnlines/2, nnlines-1; i < j; i, j = i+1, j-1 {
		tmpMu[i], tmpMu[j] = tmpMu[j], tmpMu[i]
		tmpI[i], tmpI[j] = tmpI[j], tmpI[i]
		tmpQ[i], tmpQ[j] = tmpQ[j], tmpQ[i]
	}
	//fmt.Println()
	if display {
		for i := 0; i < nnlines; i++ {
			fmt.Printf("%12.4f%12.3e%12.3e\n", tmpMu[i], tmpI[i], tmpQ[i])
		}
	}
	return
}

// Поиск и удаление всех промежуточных файлов с матрицами рассеяния
func CleanUp() {
	//Поиск и удаление всех файов с индикатриссой рассеяния
	files, err := filepath.Glob(".scat*")

	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}
