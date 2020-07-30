package database

import "github.com/BytecodeAgency/import-boundary-checker/examples/go-invalid-4/data/interactors"

type Database struct {
	username string
}

func New(username string) interactors.DatabaseInteractor {
	return Database{
		username: username,
	}
}

func (d Database) Username() string {
	return d.username
}
