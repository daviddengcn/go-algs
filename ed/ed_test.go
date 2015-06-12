package ed

import (
	"fmt"
	"testing"

	"github.com/golangplus/testing/assert"
)

type unitStr struct {
	Base
	a, b []rune
}

func (us *unitStr) CostOfChange(iA, iB int) int {
	if us.a[iA] == us.b[iB] {
		return 0
	}

	return 1
}

func TestString(t *testing.T) {
	test := func(a, b string, d int) {
		actD := String(a, b)
		assert.Equal(t, fmt.Sprintf("Edit-distance between %s and %s", a, b), actD, d)

		actD = EditDistance(&unitStr{Base{len(a), len(b), 1}, []rune(a), []rune(b)})
		assert.Equal(t, fmt.Sprintf("Edit-distance between %s and %s", a, b), actD, d)
	}

	test("abcd", "bcde", 2)
	test("abcde", "", 5)
	test("", "abcde", 5)
	test("", "", 0)
	test("abcde", "abcde", 0)
	test("abcde", "dabce", 2)
	test("abcde", "abfde", 1)
}

func TestStringFull(t *testing.T) {
	test := func(a, b string, d int, lcs string) {
		actD, actLcs := StringFull(a, b)
		assert.Equal(t, fmt.Sprintf("Edit-distance between %s and %s", a, b), actD, d)
		assert.Equal(t, fmt.Sprintf("Longest-common-string between %s and %s", a, b), actLcs, lcs)
	}

	test("abcd", "bcde", 2, "bcd")
	test("abcde", "", 5, "")
	test("", "abcde", 5, "")
	test("", "", 0, "")
	test("abcde", "abcde", 0, "abcde")
	test("abcde", "dabce", 2, "abce")
	test("abcde", "abfde", 1, "abde")
}

func ExampleString() {
	fmt.Println(String("abcde", "bfdeg"))

	fmt.Println(StringFull("abcde", "bfdeg"))
	// Output:
	// 3
	// 3 bde
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
	}

	return 100
}

func (in *stringInterface) CostOfDel(iA int) int {
	return 100 + iA
}

func (in *stringInterface) CostOfIns(iB int) int {
	return 110 + iB
}

func TestEditDistance(t *testing.T) {
	test := func(a, b string, d int) {
		actD := EditDistance(&stringInterface{[]rune(a), []rune(b)})
		assert.Equal(t, fmt.Sprintf("Edit-distance between %s and %s", a, b), actD, d)
	}

	test("abcd", "bcde", 213)
	test("abcde", "", 510)
	test("", "abcde", 560)
	test("", "", 0)
	test("abcde", "abcde", 0)
	test("abcde", "dabce", 213)
	test("abcde", "abfde", 100)
}

func TestEditDistanceFull(t *testing.T) {
	test := func(a, b string, d int, matA, matB []int) {
		actD, actMatA, actMatB := EditDistanceFull(&stringInterface{[]rune(a), []rune(b)})
		assert.Equal(t, fmt.Sprintf("Edit-distance between %s and %s", a, b), actD, d)
		assert.StringEqual(t, fmt.Sprintf("matA for matchting between %s and %s", a, b), actMatA, matA)
		assert.StringEqual(t, fmt.Sprintf("matB for matchting between %s and %s", a, b), actMatB, matB)
	}

	test("abcd", "bcde", 213, []int{-1, 0, 1, 2}, []int{1, 2, 3, -1})
	test("abcde", "", 510, []int{-1, -1, -1, -1, -1}, []int{})
	test("", "abcde", 560, []int{}, []int{-1, -1, -1, -1, -1})
	test("", "", 0, []int{}, []int{})
	test("abcde", "abcde", 0, []int{0, 1, 2, 3, 4}, []int{0, 1, 2, 3, 4})
	test("abcde", "dabce", 213, []int{1, 2, 3, -1, 4}, []int{-1, 0, 1, 2, 4})
	test("abcde", "abfde", 100, []int{0, 1, 2, 3, 4}, []int{0, 1, 2, 3, 4})
}

func TestEditDistanceF(t *testing.T) {
	test := func(a, b string, d int) {
		ra, rb := []rune(a), []rune(b)
		actD := EditDistanceF(len(a), len(b),
			func(iA, iB int) int {
				if ra[iA] == rb[iB] {
					return 0
				}

				return 100
			},
			func(iA int) int {
				return 100 + iA
			},
			func(iB int) int {
				return 110 + iB
			})
		assert.Equal(t, fmt.Sprintf("Edit-distance between %s and %s", a, b), actD, d)
	}

	test("abcd", "bcde", 213)
	test("abcde", "", 510)
	test("", "abcde", 560)
	test("", "", 0)
	test("abcde", "abcde", 0)
	test("abcde", "dabce", 213)
	test("abcde", "abfde", 100)
}

func TestEditDistanceFFull(t *testing.T) {
	test := func(a, b string, d int, matA []int, matB []int) {
		ra, rb := []rune(a), []rune(b)
		actD, actMatA, actMatB := EditDistanceFFull(len(ra), len(rb),
			func(iA, iB int) int {
				if ra[iA] == rb[iB] {
					return 0
				}

				return 100
			},
			func(iA int) int {
				return 100 + iA
			},
			func(iB int) int {
				return 110 + iB
			})
		assert.Equal(t, fmt.Sprintf("Edit-distance between %s and %s", a, b), actD, d)
		assert.StringEqual(t, fmt.Sprintf("matA for matchting between %s and %s", a, b), actMatA, matA)
		assert.StringEqual(t, fmt.Sprintf("matB for matchting between %s and %s", a, b), actMatB, matB)
	}

	test("abcd", "bcde", 213, []int{-1, 0, 1, 2}, []int{1, 2, 3, -1})
	test("abcde", "", 510, []int{-1, -1, -1, -1, -1}, []int{})
	test("", "abcde", 560, []int{}, []int{-1, -1, -1, -1, -1})
	test("", "", 0, []int{}, []int{})
	test("abcde", "abcde", 0, []int{0, 1, 2, 3, 4}, []int{0, 1, 2, 3, 4})
	test("abcde", "dabce", 213, []int{1, 2, 3, -1, 4}, []int{-1, 0, 1, 2, 4})
	test("abcde", "abfde", 100, []int{0, 1, 2, 3, 4}, []int{0, 1, 2, 3, 4})
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
