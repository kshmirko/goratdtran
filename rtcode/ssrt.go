package rtcode

import (
	"fmt"
	"math"
	"os"

	"github.com/kshmirko/radtran/libmath"
)

type SSRTData struct {
	Sza, Omega, Taue, Wl float64
	P1, P2               *[]float64
	NTheta               int
}

const deg2rad = 0.017453292519943
const rad2deg = 1.0 / deg2rad

// Cosine of zenith viewing angles
var Mu = []float64{
	0.024350292663424432509,
	0.0729931217877990394495,
	0.1214628192961205544704,
	0.1696444204239928180373,
	0.2174236437400070841497,
	0.264687162208767416374,
	0.311322871990210956158,
	0.3572201583376681159504,
	0.4022701579639916036958,
	0.446366017253464087985,
	0.489403145707052957479,
	0.531279464019894545658,
	0.5718956462026340342839,
	0.611155355172393250249,
	0.6489654712546573398578,
	0.6852363130542332425636,
	0.7198818501716108268489,
	0.7528199072605318966119,
	0.7839723589433414076102,
	0.8132653151227975597419,
	0.8406292962525803627517,
	0.8659993981540928197608,
	0.8893154459951141058534,
	0.9105221370785028057564,
	0.9295691721319395758215,
	0.9464113748584028160625,
	0.9610087996520537189186,
	0.9733268277899109637419,
	0.9833362538846259569313,
	0.9910133714767443207394,
	0.9963401167719552793469,
	0.9993050417357721394569,
}

//Конструктор объекта SSRT
func NewSSRT(fname string, sza, wl float64) *SSRTData {
	f, err := os.Open(fname)
	if err != nil {
		panic("Error opening scattering file")
	}

	defer f.Close()

	var ext, sca, omega float64
	var NL int

	// читаем коэффициенты лежандра из файла
	fmt.Fscanf(f, "%f\n", &ext)
	fmt.Fscanf(f, "%f\n", &sca)
	fmt.Fscanf(f, "%f\n", &omega)
	fmt.Fscanf(f, "%d\n", &NL)

	var P1, P2 []float64
	P1 = make([]float64, NL+1)
	P2 = make([]float64, NL+1)

	var j int
	var P3, P4, P5, P6 float64
	for i := 0; i <= NL; i++ {

		fmt.Fscanf(f, "%d %f %f %f %f %f %f\n", &j, &P1[i], &P2[i], &P3, &P4, &P5, &P6)
	}

	return &SSRTData{
		Sza:    math.Cos(sza * deg2rad),
		Omega:  omega,
		Taue:   ext,
		Wl:     wl,
		NTheta: 32,
		P1:     &P1,
		P2:     &P2,
	}
}

func (v *SSRTData) L0(tau, utheta, uphi float64) float64 {
	//Прямое солнечное излучение
	delta := math.Abs(utheta-v.Sza) + math.Abs(uphi-1.0)
	//fmt.Println(delta)
	if delta > 0.0001 {
		return 0.0
	}
	u0 := v.Sza
	u := utheta
	ua := u*u0 + math.Sqrt(1.0-u*u)*math.Sqrt(1.0-u0*u0)*uphi
	F1_i := libmath.Coef2phase(v.P1, ua)
	return tau * F1_i * math.Exp(-tau/utheta) / utheta
}

func (v *SSRTData) L1(tau, utheta, uphi float64) float64 {
	delta := math.Abs(utheta - v.Sza)
	if delta <= 0.0001 {
		return 0.0
	}

	//theta = theta * deg2rad
	//theta0 := v.Sza * deg2rad
	//phi = phi * deg2rad

	u0 := v.Sza
	u := utheta

	ua := u*u0 + math.Sqrt(1.0-u*u)*math.Sqrt(1.0-u0*u0)*uphi
	F1_i := libmath.Coef2phase(v.P1, ua)

	Ld := F1_i * u0 / (u0 - u) * (math.Exp(-tau/u0) - math.Exp(-tau/u))

	return Ld
}

func (v *SSRTData) L(tau, utheta, uphi float64) float64 {
	return v.L0(tau, utheta, uphi) + v.Omega*v.L1(tau, utheta, uphi)
}

func (v *SSRTData) Dump() {

	for i := 0; i < len(Mu); i++ {
		I := v.L(v.Taue, Mu[i], 1.0)
		fmt.Printf("%12.4f%12.3e\n", math.Acos(Mu[i])*rad2deg, I)
	}
	for i := len(Mu) - 1; i >= 0; i-- {
		I := v.L(v.Taue, Mu[i], -1.0)
		fmt.Printf("%12.4f%12.3e\n", -math.Acos(Mu[i])*rad2deg, I)
	}

}
