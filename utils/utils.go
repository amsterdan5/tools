package utils

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
)

var optMap = make(map[string]string, 0)
var Sep string

func OptStrDef(name string, d string) string {
	val, ok := optMap[name]
	if ok {
		return val
	}
	return d
}

func FindProjectRoot(mod string) string {
	abs, err := filepath.Abs(mod)
	if err != nil {
		log.Print("invalid path", mod)
	}

	for abs != "" {
		p := filepath.Join(abs, "proto")
		res, _ := IsDirectory(p)
		if res {
			return abs
		}

		prev := abs
		abs = filepath.Dir(abs)
		if prev == abs {
			break
		}
	}
	return ""
}

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func FileExist(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func FileRead(filename string) (content string, err error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("read file %s err %s", filename, err)
		return "", err
	}
	return string(buf), err
}

func FileExists(name string) bool {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func AdjPathSep(src string) string {
	var f, t string
	if runtime.GOOS == "windows" {
		f = `/`
		t = `\`
	} else {
		f = `\`
		t = `/`
	}
	return strings.Replace(src, f, t, -1)
}

func FileWrite(filename string, content string) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = out.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func GetConfigDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "innertools")
}

func UnderScore2Camel(name string) string {
	var buf []byte
	toggleUpper := true
	for i := 0; i < len(name); i++ {
		if name[i] == '_' {
			toggleUpper = true
		} else {
			c := name[i]
			if toggleUpper {
				toggleUpper = false
				if c >= 'a' && c <= 'z' {
					c = c - 'a' + 'A'
				}
			}
			if c >= '0' && c <= '9' {
				toggleUpper = true
			}
			buf = append(buf, c)
		}
	}
	return string(buf)
}

// 驼峰转下划线
func Camel2UnderScore(name string) string {
	var posList []int
	i := 1
	for i < len(name) {
		if name[i] >= 'A' && name[i] <= 'Z' {
			posList = append(posList, i)
			i++
			for i < len(name) && name[i] >= 'A' && name[i] <= 'Z' {
				i++
			}
		} else {
			i++
		}
	}
	lower := strings.ToLower(name)
	if len(posList) == 0 {
		return lower
	}
	b := strings.Builder{}
	left := 0
	for _, right := range posList {
		b.WriteString(lower[left:right])
		b.WriteByte('_')
		left = right
	}
	b.WriteString(lower[left:])
	return b.String()
}

func GetTomlFromFile(fpath string) (map[string]interface{}, error) {
	var config map[string]interface{}
	_, err := toml.DecodeFile(fpath, &config)
	if err != nil {
		log.Printf("decode toml fail %s, %s", fpath, err)
		return nil, err
	}
	return config, nil
}
func SetTomlToFile(fpath string, config map[string]interface{}) error {
	var configBuffer bytes.Buffer
	e := toml.NewEncoder(&configBuffer)
	err := e.Encode(config)
	if err != nil {
		log.Printf("toml encode err %+v", err)
	}
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("can not generate file %s,Error :%v", fpath, err)
		return err
	}
	if _, err := f.Write(configBuffer.Bytes()); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

const insertCost = 1
const deleteCost = 1
const editCost = 2

func MinEditDistance(target string, source string) (distance int) {
	targetR := []rune(target)
	sourceR := []rune(source)
	n := len(targetR)
	m := len(sourceR)
	// create distance matrix
	matrix := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		matrix[i] = make([]int, m+1)
	}
	for i := 0; i <= m; i++ {
		matrix[0][i] = i
	}
	for j := 0; j <= n; j++ {
		matrix[j][0] = j
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			var insertDistance int
			var substituteDistance int
			insertDistance = matrix[i-1][j] + insertCost
			if sourceR[j-1] == targetR[i-1] {
				substituteDistance = matrix[i-1][j-1]
			} else {
				substituteDistance = matrix[i-1][j-1] + editCost
			}
			deleteDistance := matrix[i][j-1] + deleteCost
			matrix[i][j] = Min(insertDistance, substituteDistance, deleteDistance).(int)
		}
	}
	distance = matrix[n][m]
	return
}

func OptHas(name string) bool {
	_, ok := optMap[name]
	return ok
}

func FilePathSplit(path string) (dirPath string, fileName string) {
	i := strings.LastIndex(path, Sep)
	return path[:i+1], path[i+1:]
}
