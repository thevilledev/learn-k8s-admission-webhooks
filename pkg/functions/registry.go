package functions

import (
	"log"

	"github.com/thevilledev/learn-admission-controllers/pkg/types"
)

var Registry map[string]types.AdmitFunc = make(map[string]types.AdmitFunc, 0)

func registerFunction(n string, f types.AdmitFunc) {
	if _, f := Registry[n]; f {
		log.Fatalf("function with name '%s' already exists", n)
	}
	Registry[n] = f
}
