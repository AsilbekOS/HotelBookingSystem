package randcode

import (
	"math/rand"
	"strconv"
)

func RandomCode() string {
	min := 100000
	max := 999999
	intrand := rand.Intn(max-min+1) + min
	randomnumber := strconv.Itoa(intrand)

	return randomnumber
}
