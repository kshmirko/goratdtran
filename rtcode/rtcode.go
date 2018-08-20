package rtcode

/*
#cgo CFLAGS: -O3
#cgo LDFLAGS: -L. -lcalc -lmiev0 -lrt3 -L /usr/local/Cellar/gcc/8.2.0/lib/gcc/8/ -lgfortran

void do_calc1(double r0, double r1, int npts, double wl, double mre, double mim,
                double gamma, double dens,
                double hpbl, double taua,
                int numazim, double galbedo,
                int nmu,  char* out_file);
*/
import (
	"C"
)

import (
	//	"unsafe"
	"bufio"
	"fmt"
	"io"
	"os"
)

// Параметры расчета переноса радиации
type RT3Params struct {
	r0, r1, wl, gamma, dens, hpbl, taua, galbedo float64
	npts, numazim, nmu                           int
	midx                                         complex64
	out_file                                     string
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
		nmu:      32,
		out_file: "rt3.0ut",
	}
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
		C.int(v.nmu),
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
