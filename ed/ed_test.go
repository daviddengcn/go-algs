package ed

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

func o_() string {
	pc, _, _, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()
	if p := strings.LastIndexAny(name, `./\`); p >= 0 {
		name = name[p+1:]
	} // if
	fmt.Println("== BEGIN", name, "===")
	return name
}

func __(name string) {
	fmt.Println("== END", name, "===")
}

type unitStr struct {
	Base
	a, b []rune
}

func (us *unitStr) CostOfChange(iA, iB int) int {
	if us.a[iA] == us.b[iB] {
		return 0
	} // if

	return 1
}

func TestString(t *testing.T) {
	defer __(o_())

	cases := [][]interface{}{
		{"abcd", "bcde", 2},
		{"abcde", "", 5},
		{"", "abcde", 5},
		{"", "", 0},
		{"abcde", "abcde", 0},
		{"abcde", "dabce", 2},
		{"abcde", "abfde", 1}}

	for _, cs := range cases {
		d := String(cs[0].(string), cs[1].(string))
		if d != cs[2].(int) {
			t.Errorf("Edit-distance between %s and %s is expected to be %d, but %d got!", cs[0], cs[1], cs[2], d)
		} // if
	} // for case

	for _, cs := range cases {
		a, b := []rune(cs[0].(string)), []rune(cs[1].(string))
		d := EditDistance(&unitStr{Base{len(a), len(b), 1}, a, b})
		if d != cs[2].(int) {
			t.Errorf("Edit-distance between %s and %s is expected to be %d, but %d got!", cs[0], cs[1], cs[2], d)
		} // if
	} // for case
}

func TestStringFull(t *testing.T) {
	defer __(o_())

	cases := [][]interface{}{
		{"abcd", "bcde", 2, "bcd"},
		{"abcde", "", 5, ""},
		{"", "abcde", 5, ""},
		{"", "", 0, ""},
		{"abcde", "abcde", 0, "abcde"},
		{"abcde", "dabce", 2, "abce"},
		{"abcde", "abfde", 1, "abde"}}

	for _, cs := range cases {
		d, lcs := StringFull(cs[0].(string), cs[1].(string))
		if d != cs[2].(int) {
			t.Errorf("Edit-distance between %s and %s is expected to be %d, but %d got!", cs[0], cs[1], cs[2], d)
		} // if
		if lcs != cs[3].(string) {
			t.Errorf("Longest-common-string between %s and %s is expected to be %s, but %s got!", cs[0], cs[1], cs[3], lcs)
		} // if
	} // for case
}

func ExampleString() {
	fmt.Println(String("abcde", "bfdeg"))

	fmt.Println(StringFull("abcde", "bfdeg"))
	/* Output:
	3
	3 bde
	*/
}

type stringInterface struct {
	a, b []rune
}

func (in *stringInterface) LenA() int {
	return len(in.a)
}

func (in *stringInterface) LenB() int {
	return len(in.b)
}

func (in *stringInterface) CostOfChange(iA, iB int) int {
	if in.a[iA] == in.b[iB] {
		return 0
	} // if

	return 100
}

func (in *stringInterface) CostOfDel(iA int) int {
	return 100 + iA
}

func (in *stringInterface) CostOfIns(iB int) int {
	return 110 + iB
}

func TestEditDistance(t *testing.T) {
	defer __(o_())

	cases := [][]interface{}{
		{"abcd", "bcde", 213},
		{"abcde", "", 510},
		{"", "abcde", 560},
		{"", "", 0},
		{"abcde", "abcde", 0},
		{"abcde", "dabce", 213},
		{"abcde", "abfde", 100}}

	for _, cs := range cases {
		d := EditDistance(&stringInterface{[]rune(cs[0].(string)), []rune(cs[1].(string))})
		if d != cs[2].(int) {
			t.Errorf("Edit-distance between %s and %s is expected to be %d, but %d got!", cs[0], cs[1], cs[2].(int), d)
		} // if
	} // for case
}

func TestEditDistanceFull(t *testing.T) {
	defer __(o_())

	cases := [][]interface{}{
		{"abcd", "bcde", 213, []int{-1, 0, 1, 2}, []int{1, 2, 3, -1}},
		{"abcde", "", 510, []int{-1, -1, -1, -1, -1}, []int{}},
		{"", "abcde", 560, []int{}, []int{-1, -1, -1, -1, -1}},
		{"", "", 0, []int{}, []int{}},
		{"abcde", "abcde", 0, []int{0, 1, 2, 3, 4}, []int{0, 1, 2, 3, 4}},
		{"abcde", "dabce", 213, []int{1, 2, 3, -1, 4}, []int{-1, 0, 1, 2, 4}},
		{"abcde", "abfde", 100, []int{0, 1, 2, 3, 4}, []int{0, 1, 2, 3, 4}}}

	for _, cs := range cases {
		d, matA, matB := EditDistanceFull(&stringInterface{[]rune(cs[0].(string)), []rune(cs[1].(string))})
		if d != cs[2].(int) {
			t.Errorf("Edit-distance between %s and %s is expected to be %d, but %d got!", cs[0], cs[1], cs[2].(int), d)
		} // if
		if fmt.Sprint(matA) != fmt.Sprint(cs[3]) {
			t.Errorf("matA for matchting between %s and %s is expected to be %v, but %v got!", cs[0], cs[1], cs[3], matA)
		} // if
		if fmt.Sprint(matB) != fmt.Sprint(cs[4]) {
			t.Errorf("matB for matchting between %s and %s is expected to be %v, but %v got!", cs[0], cs[1], cs[4], matB)
		} // if
	} // for case
}

func TestEditDistanceF(t *testing.T) {
	defer __(o_())

	cases := [][]interface{}{
		{"abcd", "bcde", 213},
		{"abcde", "", 510},
		{"", "abcde", 560},
		{"", "", 0},
		{"abcde", "abcde", 0},
		{"abcde", "dabce", 213},
		{"abcde", "abfde", 100}}

	for _, cs := range cases {
		a, b := []rune(cs[0].(string)), []rune(cs[1].(string))
		d := EditDistanceF(len(a), len(b),
			func(iA, iB int) int {
				if a[iA] == b[iB] {
					return 0
				} // if

				return 100
			},
			func(iA int) int {
				return 100 + iA
			},
			func(iB int) int {
				return 110 + iB
			})
		if d != cs[2].(int) {
			t.Errorf("Edit-distance between %s and %s is expected to be %d, but %d got!", cs[0], cs[1], cs[2].(int), d)
		} // if
	} // for case
}

func TestEditDistanceFFull(t *testing.T) {
	defer __(o_())

	cases := [][]interface{}{
		{"abcd", "bcde", 213, []int{-1, 0, 1, 2}, []int{1, 2, 3, -1}},
		{"abcde", "", 510, []int{-1, -1, -1, -1, -1}, []int{}},
		{"", "abcde", 560, []int{}, []int{-1, -1, -1, -1, -1}},
		{"", "", 0, []int{}, []int{}},
		{"abcde", "abcde", 0, []int{0, 1, 2, 3, 4}, []int{0, 1, 2, 3, 4}},
		{"abcde", "dabce", 213, []int{1, 2, 3, -1, 4}, []int{-1, 0, 1, 2, 4}},
		{"abcde", "abfde", 100, []int{0, 1, 2, 3, 4}, []int{0, 1, 2, 3, 4}}}

	for _, cs := range cases {
		a, b := []rune(cs[0].(string)), []rune(cs[1].(string))
		d, matA, matB := EditDistanceFFull(len(a), len(b),
			func(iA, iB int) int {
				if a[iA] == b[iB] {
					return 0
				} // if

				return 100
			},
			func(iA int) int {
				return 100 + iA
			},
			func(iB int) int {
				return 110 + iB
			})
		if d != cs[2].(int) {
			t.Errorf("Edit-distance between %s and %s is expected to be %d, but %d got!", cs[0], cs[1], cs[2].(int), d)
		} // if
		if fmt.Sprint(matA) != fmt.Sprint(cs[3]) {
			t.Errorf("matA for matchting between %s and %s is expected to be %v, but %v got!", cs[0], cs[1], cs[3], matA)
		} // if
		if fmt.Sprint(matB) != fmt.Sprint(cs[4]) {
			t.Errorf("matB for matchting between %s and %s is expected to be %v, but %v got!", cs[0], cs[1], cs[4], matB)
		} // if
	} // for case
}

func ExampleEditDistanceF() {
	a, b := "abcd", "bcde"
	d := EditDistanceF(len(a), len(b), func(iA, iB int) int {
		return Ternary(a[iA] == b[iB], 0, 1)
	}, ConstCost(1), ConstCost(1))
	fmt.Println(a, b, d)
	// Output:
	// abcd bcde 2
}
