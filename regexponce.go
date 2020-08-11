package regexponce

import (
	"go/importer"
	"go/types"

	"github.com/gostaticanalysis/analysisutil"
	"github.com/gostaticanalysis/comment/passes/commentmap"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
)

const doc = "regexponce is ..."

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
	// 関数の呼び出し箇所を取得する
	// regexpの該当の関数だけを抽出する
	// どのスコープで使われているか判定する。
	// initかパッケージ変数の初期化の場合は許可する
	// コメントで許可されているところは無視する。
	fs, err := targetFuncs()
	if err != nil {
		panic(err)
	}

	pass.Report = analysisutil.ReportWithoutIgnore(pass)
	srcFuncs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs
	for _, sf := range srcFuncs {
		for _, b := range sf.Blocks {
			for _, instr := range b.Instrs {
				for _, f := range fs {
					if analysisutil.Called(instr, nil, f) {
						pass.Reportf(instr.Pos(), "%s must not be called", f.FullName())
						break
					}
				}
			}
		}
	}

	return nil, nil
}

func targetFuncs() ([]*types.Func, error) {
	fs := make([]*types.Func, 0, 4)
	fns := []string{"MustCompile", "Compile", "MustCompilePOSIX", "CompilePOSIX"}
	pkg, err := importer.Default().Import("regexp")
	if err != nil {
		return nil, err
	}
	scp := pkg.Scope()

	for _, fn := range fns {
		obj := scp.Lookup(fn)
		if obj == nil {
			continue
		}
		if f, ok := obj.(*types.Func); ok {
			fs = append(fs, f)
		}
	}

	return fs, nil
}
