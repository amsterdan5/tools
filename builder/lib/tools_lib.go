package tools_lib

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/amsterdan/tools/logger"
)

type Func func()
type FuncNode struct {
	Name  string
	Usage string
	Func  Func
}

var fList []*FuncNode
var mu sync.Mutex
var optMap = make(map[string]string, 0)
var UsageTail string
var DefaultMethod string

func Register(name, usage string, f Func) {
	mu.Lock()
	defer mu.Unlock()
	fList = append(fList, &FuncNode{Name: name, Usage: usage, Func: f})
}

func showUsage() {
	fmt.Println()
	for _, f := range fList {
		fmt.Printf("\t-f %s %s\n", f.Name, f.Usage)
	}
	if UsageTail != "" {
		fmt.Printf("\n%s\n", UsageTail)
	}
	if DefaultMethod != "" {
		fmt.Println()
		fmt.Printf("default action: %s\n", DefaultMethod)
	}
	fmt.Println()
	os.Exit(1)
}

func OptInt(name string) int {
	val := OptStr(name)
	intVal, err := strconv.Atoi(val)
	if err == nil {
		return intVal
	}
	fmt.Printf("opt %s required type int, input %s\n\n", name, val)
	os.Exit(1)
	return 0
}

func OptFloat32(name string) float32 {
	val := OptStr(name)
	intVal, err := strconv.ParseFloat(val, 32)
	if err == nil {
		return float32(intVal)
	}
	fmt.Printf("opt %s required type uint32, input %s\n\n", name, val)
	os.Exit(1)
	return 0
}

func OptFloat64(name string) float64 {
	val := OptStr(name)
	intVal, err := strconv.ParseFloat(val, 64)
	if err == nil {
		return intVal
	}
	fmt.Printf("opt %s required type uint32, input %s\n\n", name, val)
	os.Exit(1)
	return 0
}

func OptUint32(name string) uint32 {
	val := OptStr(name)
	intVal, err := strconv.ParseUint(val, 10, 32)
	if err == nil {
		return uint32(intVal)
	}
	fmt.Printf("opt %s required type uint32, input %s\n\n", name, val)
	os.Exit(1)
	return 0
}

func OptUint32Def(name string, def uint32) uint32 {
	val := OptStrDef(name, "")
	if val == "" {
		return def
	}
	intVal, err := strconv.ParseUint(val, 10, 32)
	if err == nil {
		return uint32(intVal)
	}
	fmt.Printf("opt %s required type uint32, input %s\n\n", name, val)
	os.Exit(1)
	return 0
}

func OptUint64(name string) uint64 {
	val := OptStr(name)
	intVal, err := strconv.ParseUint(val, 10, 64)
	if err == nil {
		return intVal
	}
	fmt.Printf("opt %s required type uint64, input %s\n\n", name, val)
	os.Exit(1)
	return 0
}

func OptUint64Def(name string, def uint64) uint64 {
	val := OptStrDef(name, "")
	if val == "" {
		return def
	}
	intVal, err := strconv.ParseUint(val, 10, 64)
	if err == nil {
		return intVal
	}
	fmt.Printf("opt %s required type uint64, input %s\n\n", name, val)
	os.Exit(1)
	return 0
}

func OptInt64(name string) int64 {
	val := OptStr(name)
	intVal, err := strconv.ParseInt(val, 10, 64)
	if err == nil {
		return intVal
	}
	fmt.Printf("opt %s required type int, input %s, err %s\n\n", name, val, err)
	os.Exit(1)
	return 0
}

func OptInt64Def(name string, d uint64) int64 {
	val := OptStrDef(name, fmt.Sprintf("%d", d))
	intVal, err := strconv.ParseInt(val, 10, 64)
	if err == nil {
		return intVal
	}
	fmt.Printf("opt %s required type int, input %s, err %s\n\n", name, val, err)
	os.Exit(1)
	return 0
}

func OptIntDef(name string, d int) int {
	val := OptStrDef(name, strconv.Itoa(d))
	intVal, err := strconv.Atoi(val)
	if err == nil {
		return intVal
	}
	fmt.Printf("opt %s required type int, input %s\n", name, val)
	os.Exit(1)
	return d
}

func OptStr(name string) string {
	val, ok := optMap[name]
	if ok {
		return val
	}
	fmt.Printf("missed parameter -%s\n\n", name)
	os.Exit(1)
	return ""
}

func OptStrDef(name string, d string) string {
	val, ok := optMap[name]
	if ok {
		return val
	}
	return d
}

func OptStrSlice(name string) []string {
	str := OptStrDef(name, "")
	slice := strings.Split(str, StringSeparator)
	return slice
}

func OptUint32Slice(name string) []uint32 {
	strSlice := OptStrSlice(name)
	intSlice := make([]uint32, 0)
	for _, str := range strSlice {
		intVal, err := strconv.ParseUint(str, 10, 32)
		if err == nil {
			intSlice = append(intSlice, uint32(intVal))
		} else {
			fmt.Printf("opt %s required type int, input %s\n", name, str)
			os.Exit(1)
		}
	}
	return intSlice
}

func OptUint64Slice(name string) []uint64 {
	strSlice := OptStrSlice(name)
	var intSlice []uint64
	for _, str := range strSlice {
		intVal, err := strconv.ParseUint(str, 10, 64)
		if err == nil {
			intSlice = append(intSlice, intVal)
		} else {
			fmt.Printf("opt %s required type int, input %s\n", name, str)
			os.Exit(1)
		}
	}
	return intSlice
}

func OptUint64SliceDef(name string, def []uint64) []uint64 {
	val := OptStrDef(name, "")
	if val == "" {
		return def
	}
	strSlice := strings.Split(val, StringSeparator)
	var intSlice []uint64
	for _, str := range strSlice {
		intVal, err := strconv.ParseUint(str, 10, 64)
		if err == nil {
			intSlice = append(intSlice, intVal)
		} else {
			fmt.Printf("opt %s required type int, input %s\n", name, str)
			os.Exit(1)
		}
	}
	return intSlice
}

func OptInt64Slice(name string) []int64 {
	strSlice := OptStrSlice(name)
	var intSlice []int64
	for _, str := range strSlice {
		intVal, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			intSlice = append(intSlice, intVal)
		} else {
			fmt.Printf("opt %s required type int, input %s\n", name, str)
			os.Exit(1)
		}
	}
	return intSlice
}

func OptHas(name string) bool {
	_, ok := optMap[name]
	return ok
}

func OptBool(name string) bool {
	val := OptStr(name)
	val = strings.ToLower(val)
	if val == "0" || val == "false" {
		return false
	}
	return true
}

func OptBoolDef(name string, b bool) bool {
	val := OptStrDef(name, "")
	if len(val) == 0 {
		return b
	}
	val = strings.ToLower(val)
	if val == "0" || val == "false" {
		return false
	}
	return true
}

func OptStrPrompt(tips string) (string, error) {
	fmt.Println()
	fmt.Println(tips)
	var optValue string
	_, err := fmt.Scanf("%s", &optValue)
	if err != nil {
		logger.Errorf("err:%v", err)
		return "", err
	}
	return optValue, nil
}

const (
	StringSeparator = ","
)

func Run() {
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "-h" || os.Args[i] == "--help" {
			showUsage()
		}
		if strings.HasPrefix(os.Args[i], "-") {
			opt := strings.Trim(os.Args[i], "-")
			if i+1 < len(os.Args) {
				if strings.HasPrefix(os.Args[i+1], "-") {
					optMap[opt] = ""
				} else {
					optMap[opt] = os.Args[i+1]
				}
			} else {
				optMap[opt] = ""
			}
		}
	}
	f, ok := optMap["f"]
	if !ok {
		if DefaultMethod != "" {
			f = DefaultMethod
		} else {
			fmt.Printf("\nmissed -f opt\n\n")
			showUsage()
		}
	}
	name := path.Base(strings.Replace(os.Args[0], "\\", "/", -1))
	if strings.HasPrefix(name, "./") {
		name = name[2:]
	}
	// get handler
	var fu Func
	for _, x := range fList {
		if strings.ToLower(x.Name) == strings.ToLower(f) {
			fu = x.Func
			break
		}
	}
	if fu == nil {
		fmt.Printf("\nnot found func %s\n\n", f)
		showUsage()
	}
	regExitSignals()
	fu()
}
