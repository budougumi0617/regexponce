package a

import (
	"fmt"
	"regexp"
)

var okRegexp = regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`) // OK

func init() {
	var validID = regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`) // OK
	fmt.Println(validID.MatchString("adam[23]"))
}

func initfake() {
	validID, _ := regexp.Compile(`^[a-z]+\[[0-9]+\]$`) // want `regexp.Compile must be called only once at initialize`
	fmt.Println(validID.MatchString("adam[23]"))
}

func f() {
	var validID = regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`) // want `regexp.MustCompile must be called only once at initialize`
	fmt.Println(validID.MatchString("adam[23]"))

	// lint:ignore regexponce allowed
	validID = regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`) // OK
	hoge := regexp.MustCompile
	hoge(`^[a-z]+\[[0-9]+\]$`) // want `regexp.MustCompile must be called only once at initialize`
}
func main() {
	var validID = regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`) // OK because main function runs only once.
	fmt.Println(validID.MatchString("adam[23]"))

	for i := 0; i < 10; i++ {
		validID = regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`) // want `regexp.MustCompile must be called only once at initialize`
	}
	x := 10
	if x < 10 {
		validID = regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`) // OK
	}
}
