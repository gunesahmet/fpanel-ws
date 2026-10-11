package Helpers

import (
	"Logger"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var Logging Logger.Interface

func Trim(str string) string {
	s1 := strings.TrimSpace(str)
	s2 := strings.Trim(s1, "\t")
	s3 := strings.Trim(s2, "\n")
	return s3
}

func Timer(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func RandFloats(min float64, max float64, n int) []float64 {
	rand.Seed(time.Now().UnixNano())
	res := make([]float64, n)
	var floatn = float64(n)
	max = max / floatn
	for i := range res {
		res[i] = min + rand.Float64()*max
		//time.Sleep(time.Duration(100*time.Millisecond))
	}
	return res
}

func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func PrettyPrint(data []byte) {
	var res map[string]interface{}

	err := json.Unmarshal(data, &res)

	if err != nil {
		return
	}
	result, _ := json.MarshalIndent(res, "", "    ")
	fmt.Println(string(result))
}

func CurrentTimestamp() int64 {
	msec := time.Now().UnixNano() / 1000
	return msec
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func CreateDirectory(path string, recursive bool) bool {
	if recursive {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return false
		}
	} else {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return false
		}
	}
	return true
}

func ReadFile(filename string) []string {
	readFile, err := os.Open(filename)

	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileTextLines []string

	for fileScanner.Scan() {
		if len(fileScanner.Text()) > 0 {
			fileTextLines = append(fileTextLines, fileScanner.Text())
		}
	}

	_ = readFile.Close()

	return fileTextLines
}

func TrimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

func WriteFile(content string, filePath string) error {
	if FileExists(filePath) {
		err := os.Remove(filePath)
		if err != nil {
			return err
		}
	}
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	} else {
		_, err := f.WriteString(content)
		if err != nil {
			_ = f.Close()
			return err
		}
		err = f.Close()
		if err != nil {
			//fmt.Println(err)
			return err
		}
	}
	return nil
}

func Exec(command string, args ...string) (error, string) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	_ = cmd.Start()
	err := cmd.Wait()

	if err == nil {
		return nil, Trim(out.String())
	}
	return err, ""
}

func GetParamWithNeedle(param string, defaultVal string, needle string, needleVal string) string {
	if len(param) == 0 {
		param = defaultVal
	} else {
		if param == needle {
			param = needleVal
		} else {
			param = defaultVal
		}
	}
	return param
}

func GetParam(param string, defaultVal string) string {
	if len(param) == 0 {
		param = defaultVal
	} else {
	}
	return param
}

func ExecWithoutOutPut(command string, args ...string) {
	err := exec.Command(command, args...).Run()
	if err != nil {
		_, fn, line, _ := runtime.Caller(1)
		fmt.Println(fmt.Sprintf("Fn: %s Line: %d", fn, line))
		fmt.Println(command)
		fmt.Println(err)
	}
}

func ExecWithoutWait(command string, args ...string) {
	err := exec.Command(command, args...).Start()
	if err != nil {
		_, fn, line, _ := runtime.Caller(1)
		fmt.Println(fmt.Sprintf("Fn: %s Line: %d", fn, line))
		fmt.Println(command)
		fmt.Println(err)
	}
}

func RemoveAll(str string, remove string) string {
	str = strings.ReplaceAll(str, remove, "")
	if strings.Contains(str, remove) {
		return RemoveAll(str, remove)
	}
	return str
}

func ReplaceAllOld(str string, old string, new string) string {
	str = strings.ReplaceAll(str, old, new)
	if strings.Contains(str, old) {
		return ReplaceAllOld(str, old, new)
	}
	return str
}

func ReplaceAll(str string, old string, new string) string {
	re := regexp.MustCompile(old)
	return re.ReplaceAllString(str, new)
	//str = strings.ReplaceAll(str, old, new)
	//if strings.Contains(str, old) {
	//	return ReplaceAll(str, old, new)
	//}
	//return str
}

func ReplaceAllTest(subj *string, old string, new string) string {
	re := regexp.MustCompile(old)
	return re.ReplaceAllString(*subj, new)
	//str = strings.ReplaceAll(str, old, new)
	//if strings.Contains(str, old) {
	//	return ReplaceAll(str, old, new)
	//}
	//return str
}

func IsFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		return false
	case mode.IsRegular():
		return true
	}
	return false
}

func IsDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		return true
	case mode.IsRegular():
		return false
	}
	return false
}

func ParseVersion(s string, width int) int64 {
	strList := strings.Split(s, ".")
	format := fmt.Sprintf("%%s%%0%ds", width)
	v := ""
	for _, value := range strList {
		v = fmt.Sprintf(format, v, value)
	}
	var result int64
	var err error
	if result, err = strconv.ParseInt(v, 10, 64); err != nil {
		fmt.Printf("ugh: parseVersion(%s): error=%s", s, err)
		return 0
	}
	//fmt.Printf("parseVersion: [%s] => [%s] => [%d]\n", s, v, result);
	return result
}

func RemoveAfter(s string, d string) string {
	parts := strings.Split(s, d)
	return parts[0]
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func CopyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer func() {
		_ = sourcefile.Close()
	}()
	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer func() {
		_ = destfile.Close()
	}()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}
	}

	return
}

func CopyDir(source string, dest string) (err error) {

	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)
	objects, err := directory.Readdir(-1)

	for _, obj := range objects {
		sourcefilepointer := source + "/" + obj.Name()
		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// create sub-directories - recursively
			err = CopyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// perform copy
			err = CopyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return
}

func FindFile(directory string, fileName string, path *string) {
	directory = ReplaceAll(directory, "//", "/")
	files, err := ioutil.ReadDir(directory)

	if err == nil {
		for _, file := range files {
			if file.Mode().IsDir() {
				FindFile(directory+file.Name()+"/", fileName, path)
			} else {
				if strings.Compare(file.Name(), fileName) == 0 {
					*path = directory
				}
			}
		}
	}
}

func HandleError(err error, code int, msg string) {
	_, fn, line, _ := runtime.Caller(1)
	fn = ReplaceAll(fn, ".go", "")
	fn = ReplaceAll(fn, "/usr/local/src/sources/fpanel/fws/src/", "")
	fn = ReplaceAll(fn, "/usr/local/src/sources/fpanel/fws/", "")
	fn = ReplaceAll(fn, "H:/Common/fws/src/", "")
	fn = ReplaceAll(fn, "H:/Common/fws/", "")
	if err == nil {
		err = errors.New("nullptr")
	}
	Logging.Error(code, "[%v] %v : %v %v", msg, fn, line, err)
}

func AppendFile(str string, filePath string, perm os.FileMode) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	if _, err := f.WriteString(str + "\n"); err != nil {
		return err
	}
	return nil
}

func SearchStringInFile(needle string, filePath string) bool {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.Index(line, needle) > -1 {
			return true
		}
	}
	return true
}

func SearchStringInFileRegexp(rx string, filePath string) bool {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(b), "\n") {
		result, err := regexp.MatchString(rx, line)
		if err != nil {
			return false
		}
		if result {
			return true
		}

	}
	return false
}

func ClearEmptyLinesInFile(filePath string) {
	ExecWithoutOutPut("/bin/sh", "-c", "sed -i '/^$/d' "+filePath)
}

func Split(s string, d string) []string {
	a := strings.Split(s, d)
	return a
}

func TrimLeftChars(s string, n int) string {
	m := 0
	for i := range s {
		if m >= n {
			return s[i:]
		}
		m++
	}
	return s[:0]
}

func Sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func Find(slice []string, key int) bool {
	for i, _ := range slice {
		if i == key {
			return true
		}
	}
	return false
}
