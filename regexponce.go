package regexponce

import (
	"go/types"
	"strings"

	"github.com/gostaticanalysis/analysisutil"
	"github.com/gostaticanalysis/comment/passes/commentmap"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

const doc = `regexp.Compile and below functions should be called at once for performance.
- regexp.MustCompile
- regexp.CompilePOSIX
- regexp.MustCompilePOSIX

Allow call in init, and main(exept for in for loop) functions because each function is called only once.
`

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "regexponce",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
		commentmap.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	fs := targetFuncs(pass)
	if len(fs) == 0 {
		return nil, nil
	}

	pass.Report = analysisutil.ReportWithoutIgnore(pass)
	srcFuncs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs
	for _, sf := range srcFuncs {
		if strings.HasPrefix(sf.Name(), "init#") {
			continue
		}

		for _, b := range sf.Blocks {
			var skipped bool
			if strings.HasPrefix(sf.Name(), "main") {
				skipped = true
			}
			if skipped && inFor(b) {
				skipped = false
			}
			if skipped {
				continue
			}
			for _, instr := range b.Instrs {
				for _, f := range fs {
					if Func(instr, f) {
						pass.Reportf(instr.Pos(), "%s must be called only once at initialize", f.FullName())
						break
					}
				}
			}
		}
	}

	return nil, nil
}

func inFor(b *ssa.BasicBlock) bool {
	p := b
	for {
		if p.Comment == "for.body" {
			return true
		}
		p = p.Idom()
		if p == nil {
			break
		}
	}
	return false
}

// Func returns true when f is called in the instr.
// If recv is not nil, Called also checks the receiver.
func Func(instr ssa.Instruction, f *types.Func) bool {
	call, ok := instr.(ssa.CallInstruction)
	if !ok {
		return false
	}

	common := call.Common()
	if common == nil {
		return false
	}

	callee := common.StaticCallee()
	if callee == nil {
		return false
	}

	fn, ok := callee.Object().(*types.Func)
	if !ok {
		return false
	}

	return fn == f
}

func targetFuncs(pass *analysis.Pass) []*types.Func {
	fs := make([]*types.Func, 0, 4)
	path := "regexp"
	fns := []string{"MustCompile", "Compile", "MustCompilePOSIX", "CompilePOSIX"}

	imports := pass.Pkg.Imports()
	for i := range imports {
		if path == analysisutil.RemoveVendor(imports[i].Path()) {
			for _, fn := range fns {
				obj := imports[i].Scope().Lookup(fn)
				if obj == nil {
					continue
				}

				if f, ok := obj.(*types.Func); ok {
					fs = append(fs, f)
				}
			}
		}
	}

	return fs
}
