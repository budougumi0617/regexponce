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

func doNotWarnWithVariable(input string) {
	regexp.MustCompile(input) // OK because function parameter is a variable

	returnWord := func(input string) string {
		return input
	}
	regexp.MustCompile(returnWord(input)) // OK because function parameter is a function call

	const constVal = ".*"
	regexp.MustCompile(input + constVal) // OK because function parameter contains a variable

	regexp.MustCompile(constVal)            // want `regexp.MustCompile must be called only once at initialize`
	regexp.MustCompile(constVal + constVal) // want `regexp.MustCompile must be called only once at initialize`
}

func main() {
	var validID = regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`) // OK because main function runs only once.
	fmt.Println(validID.MatchString("adam[23]"))

	x := 10
	for i := 0; i < 10; i++ {
		validID = regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`) // want `regexp.MustCompile must be called only once at initialize`
		if x < 10 {
			validID = regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`) // want `regexp.MustCompile must be called only once at initialize`
		}
	}

	if x < 10 {
		validID = regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`) // OK
		for i := 0; i < 10; i++ {
			validID = regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`) // want `regexp.MustCompile must be called only once at initialize`
		}
	}
}
