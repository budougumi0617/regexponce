package main

import (
	"github.com/budougumi0617/regexponce"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(regexponce.Analyzer) }

