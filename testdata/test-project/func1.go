package testproject

import "fmt"

func Func1(bool1 bool, bool2 bool) {
	fmt.Println("func1")

	if bool1 {
		fmt.Println("bool1", bool1)

		fmt.Println("end bool1", bool2)
	}

	if bool2 {
		fmt.Println("bool2", bool2)

		fmt.Println("end bool2", bool2)
	}

	fmt.Println("end func1")
}
