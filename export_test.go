package regexponce

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var TargetFuncs func(*analysis.Pass) []*types.Func = targetFuncs
