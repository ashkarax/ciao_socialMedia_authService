package randnumgene_authSvc

import (
	"math/rand"

	interface_randnumgene_authSvc "github.com/ashkarax/ciao_socialMedia_authService/utils/random_number/interface"
)

type RandomNum struct{}

func NewRandomNumUtil() interface_randnumgene_authSvc.IRandGene {
	return &RandomNum{}
}

func (rn RandomNum) RandomNumber() int {
	randomInt := rand.Intn(9000) + 1000
	return randomInt
}
