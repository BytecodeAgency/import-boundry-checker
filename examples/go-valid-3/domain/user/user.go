package user

import "github.com/BytecodeAgency/import-boundary-checker/examples/go-valid-3/data/interactions"

func Validate() error {
	return interactions.Validate()
}

func GetTheUser() string {
	return "admin"
}
