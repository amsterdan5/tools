package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"

	"github.com/amsterdan/tools/logger"
)

func showUsage() {
	fmt.Printf("\n%s -f <tools>.go\n\n", os.Args[0])
	os.Exit(1)
}

type funcNode struct {
	Name  string
	Usage string
}

const usageTxt = "usage:"

func init() {
	logger.InitLogger(logger.LogModConsole, "")
}

func main() {
	fnPtr := flag.String("f", "", "tools source file")
	flag.Parse()
	if *fnPtr == "" {
		showUsage()
	}
	fn := *fnPtr
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, fn, nil, parser.ParseComments)
	if err != nil {
		logger.Errorf("parse file error, fn %s, err %s", fn, err)
		return
	}
	pkgName := f.Name
	if pkgName == nil || pkgName.Name != "main" {
		logger.Errorf("package name must be 'main'")
		return
	}
	// get the directory
	dir := filepath.Dir(fn)
	if dir != "." {
		err := os.Chdir(dir)
		if err != nil {
			fmt.Printf("chdir to %s failed, %s", dir, err)
			os.Exit(1)
		}
	}
	base := filepath.Base(fn)
	toolsName := base
	if strings.HasSuffix(base, ".go") {
		toolsName = base[:len(base)-3]
	}
	outPath := fmt.Sprintf("%s_autogen.go", toolsName)
	logger.Infof("out path %s", outPath)
	funcList := make([]*funcNode, 0)
	for _, decl := range f.Decls {
		switch t := decl.(type) {
		case *ast.FuncDecl:
			funcName := t.Name.Name
			if funcName == "" {
				continue
			}
			if !unicode.IsUpper(rune(funcName[0])) {
				logger.Infof("- skip %s, first char not upper", funcName)
				continue
			}
			if t.Recv != nil {
				logger.Infof("- skip %s, only support raw function", funcName)
				continue
			}
			if len(t.Type.Params.List) > 0 {
				logger.Infof("- skip %s, params not empty", funcName)
				continue
			}
			usage := ""
			if t.Doc != nil {
				for _, comment := range t.Doc.List {
					txt := comment.Text
					p := strings.Index(txt, usageTxt)
					if p >= 0 {
						usage = strings.TrimSpace(txt[p+len(usageTxt):])
					}
				}
			}
			logger.Infof("func name %s, usage %s", funcName, usage)
			funcList = append(funcList, &funcNode{Name: funcName, Usage: usage})
		}
	}
	if len(funcList) == 0 {
		logger.Info("func list empty, skip")
		return
	}
	buf := ""
	imp := `package main
import (
	tools_lib "github.com/amsterdan/tools/builder/lib"
)
`
	buf += imp
	wrapperTemplate := `
func wrapper%s() {
%s()
}
`
	mainTemplate := `
func main() {
%s
tools_lib.Run()
}
`
	for _, f := range funcList {
		buf += fmt.Sprintf(wrapperTemplate, f.Name, f.Name)
	}
	regList := make([]string, 0)
	for _, f := range funcList {
		regList = append(
			regList,
			fmt.Sprintf("\ttools_lib.Register(\"%s\", `%s`, wrapper%s)",
				f.Name, f.Usage, f.Name))
	}
	regBlock := strings.Join(regList, "\n")
	buf += fmt.Sprintf(mainTemplate, regBlock)
	os.WriteFile(outPath, []byte(buf), 0666)
	logger.Infof("gen %s success", outPath)
	binName := toolsName
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	// build
	cmd := exec.Command(
		"go", "build", "-o", binName,
		fmt.Sprintf("%s.go", toolsName),
		fmt.Sprintf("%s_autogen.go", toolsName))
	out, err := cmd.CombinedOutput()
	if err != nil {
		println("build error", err.Error(), string(out))
		os.Exit(1)
	}
	print(string(out))
	logger.Info("build success")
}
