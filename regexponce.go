package regexponce

import (
	"fmt"
	"go/importer"
	"go/types"
	"strings"

	"github.com/gostaticanalysis/analysisutil"
	"github.com/gostaticanalysis/comment/passes/commentmap"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
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
	fo := analysisutil.ObjectOf(pass, "regexp", "MustCompile")
	f := fo.(*types.Func)

	f2 := fromDefault()
	fmt.Println("Pkg", f.Pkg(), "-------", f2.Pkg())
	fmt.Println("Name", f.FullName(), "-------", f2.FullName())
	fmt.Println("Id", f.Id(), "-------", f2.Id())
	fmt.Println("Pos", f.Pos(), "-------", f2.Pos())
	fmt.Println("Scope", f.Scope(), "-------", f2.Scope())
	return nil, nil
}

func fromDefault() *types.Func {
	pkg, err := importer.Default().Import("regexp")
	if err != nil {
		panic(err) // !!!!!!!! can't find import: "regexp" on playgournd
	}
	obj := pkg.Scope().Lookup("MustCompile")
	if f, ok := obj.(*types.Func); ok {
		return f
	}
	panic("unreachable")
}

func run2(pass *analysis.Pass) (interface{}, error) {
	// 関数の呼び出し箇所を取得する
	// regexpの該当の関数だけを抽出する
	// どのスコープで使われているか判定する。
	// initかパッケージ変数の初期化の場合は許可する
	// コメントで許可されているところは無視する。
	fs, err := targetFuncs()
	if err != nil {
		panic(err)
	}
	fs2 := restrictedFuncs(pass, "regexp.MustCompile")
	fs3 := analysisutil.ObjectOf(pass, "regexp", "MustCompile")
	fmt.Printf("fs: %#v\n", fs[0])
	fmt.Printf("fs2: %#v\n", fs2[0])
	fmt.Printf("fs2: %#v\n", fs3)
	fmt.Printf("result: %t\n", fs2[0] == fs[0])

	pass.Report = analysisutil.ReportWithoutIgnore(pass)
	srcFuncs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs
	for _, sf := range srcFuncs {
		for _, b := range sf.Blocks {
			for _, instr := range b.Instrs {
				// 	fmt.Printf("instr %#v\n", instr)
				for _, f := range fs {
					// for _, f := range fs2 {
					// for _, f := range []*types.Func{fs3.(*types.Func)} {
					//fmt.Println("try!", f.FullName())
					if Func(instr, nil, f) {
						fmt.Println("found!!!")
						fmt.Printf("%d: %s must be called only once at initialize\n", instr.Pos(), f.FullName())
						pass.Reportf(instr.Pos(), "%s must be called only once at initialize", f.FullName())
						break
					}
				}
			}
		}
	}

	return nil, nil
}
func restrictedFuncs(pass *analysis.Pass, names string) []*types.Func {
	var fs []*types.Func
	for _, fn := range strings.Split(names, ",") {
		ss := strings.Split(strings.TrimSpace(fn), ".")

		// package function: pkgname.Func
		if len(ss) < 2 {
			continue
		}
		f, _ := analysisutil.ObjectOf(pass, ss[0], ss[1]).(*types.Func)
		if f != nil {
			fs = append(fs, f)
			continue
		}

		// method: (*pkgname.Type).Method
		if len(ss) < 3 {
			continue
		}
		pkgname := strings.TrimLeft(ss[0], "(")
		typename := strings.TrimRight(ss[1], ")")
		if pkgname != "" && pkgname[0] == '*' {
			pkgname = pkgname[1:]
			typename = "*" + typename
		}

		typ := analysisutil.TypeOf(pass, pkgname, typename)
		if typ == nil {
			continue
		}

		m := analysisutil.MethodOf(typ, ss[2])
		if m != nil {
			fs = append(fs, m)
		}
	}

	return fs
}

// Func returns true when f is called in the instr.
// If recv is not nil, Called also checks the receiver.
func Func(instr ssa.Instruction, recv ssa.Value, f *types.Func) bool {

	// fmt.Println("Func start!!")
	call, ok := instr.(ssa.CallInstruction)
	if !ok {
		return false
	}

	// fmt.Println("CalleInstruction")
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
	// fmt.Println("got fn!", fn.FullName())

	if recv != nil &&
		common.Signature().Recv() != nil &&
		(len(common.Args) == 0 && recv != nil || common.Args[0] != recv &&
			!referrer(recv, common.Args[0])) {
		return false
	}
	// fmt.Println("convert", fn.FullName(), "and", f.FullName())
	// fmt.Println("Scope", fn.Scope(), "Scope", f.Scope())
	// fmt.Println("Parent", fn.Parent(), "-------", f.Parent())
	fmt.Println("Pkg", fn.Pkg(), "-------", f.Pkg())
	fmt.Println("Pos", fn.Pos(), "-------", f.Pos())
	fmt.Println("==", fn.Pkg().Path() == f.Pkg().Path())
	// if d := cmp.Diff(fn, f); len(d) != 0 {
	// 	fmt.Printf("differs: (-got +want)\n%s", d)
	// }

	return fn == f
}

func referrer(a, b ssa.Value) bool {
	return isReferrerOf(a, b) || isReferrerOf(b, a)
}

func isReferrerOf(a, b ssa.Value) bool {
	if a == nil || b == nil {
		return false
	}
	if b.Referrers() != nil {
		brs := *b.Referrers()

		for _, br := range brs {
			brv, ok := br.(ssa.Value)
			if !ok {
				continue
			}
			if brv == a {
				return true
			}
		}
	}
	return false
}

func targetFuncs() ([]*types.Func, error) {
	fs := make([]*types.Func, 0, 4)
	// fns := []string{"MustCompile", "Compile", "MustCompilePOSIX", "CompilePOSIX"}
	fns := []string{"MustCompile"}
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
