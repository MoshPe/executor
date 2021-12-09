package cmd

import (
	"database/sql"
	"executor/pkg"
	"github.com/compose-spec/compose-go/types"
)

var Project *types.Project
var ConfigFile = pkg.File{
	FileName: "",
	Path:     "",
	IsLoaded: false,
}
var db *sql.DB
