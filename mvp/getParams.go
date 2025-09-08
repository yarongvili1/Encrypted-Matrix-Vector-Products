package mvp

import (
	"math"
)

func findb(sec float64, k int) int {
	for b := 3; b < 10000; b++ {
		if float64(k) < float64(b-1)*sec/math.Log2(float64(b)) {
			return b - 1
		}
	}
	return 10000
}

func findk(sec float64, b int) int {
	for k := 2; k < int(sec); k++ {
		if float64(k)*math.Log2(float64(k)) >= sec*float64(b-1) {
			return k
		}
	}
	return int(sec)
}

func prms(sec, f float64, ll int) (uint32, uint32, uint32, uint32) {
	b := int(math.Floor(f)) + 1
	k := int(math.Ceil(sec * float64(b-1) / math.Log2(float64(b))))
	l2 := int(math.Floor(float64(k) / (f - 1)))
	if l2 < ll {
		k = int(math.Ceil(float64(ll) * (f - 1)))
		b = findb(sec, k)
	} else {
		ll = l2
	}
	n := k + ll
	n = int(math.Ceil(float64(n)/float64(b))) * b
	k = n - ll
	s := float64(n) / float64(b)

	return uint32(ll), uint32(k), uint32(s), uint32(b)

}

func prms2(sec, f float64, ll int) (uint32, uint32, uint32, uint32, uint32) {
	b := int(math.Floor(f)) + 1
	k := findk(sec, b)
	l2 := int(math.Floor(float64(k) / (f - 1)))
	if l2 < ll {
		k = int(math.Ceil(float64(ll) * (f - 1)))
		b = 1 + int(math.Floor(float64(k)*math.Log2(float64(k))/sec))
	} else {
		ll = l2
	}
	n := k + ll
	n = int(math.Ceil(float64(n)/float64(b))) * b
	k = n - ll
	s := float64(n) / float64(b)

	return uint32(n), uint32(ll), uint32(k), uint32(s), uint32(b)
}
