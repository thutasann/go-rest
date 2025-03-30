package fundamentals

import "fmt"

func PointerSampleOne() {
	fmt.Println("---> Pointer Sample One ")
	num := 10
	ptr := &num
	fmt.Println(num)
	println("ptr --> ", ptr)
}

func ModifyPointerFunctionSample() {
	println("----> Modify Pointer Function")
	num := 10
	println("Before : ", num)

	modifyFn(&num)
	println("After : ", num)
}

func modifyFn(x *int) {
	*x = 100
}

type User struct {
	Name string
	Age  int
}

// Pointer receiver modifies the struct
func (u *User) birthday() {
	u.Age++
}

// Pointers to Structs & Methods
func PointerStructSampleOne() {
	user := User{Name: "Alice", Age: 25}
	user.birthday()
	fmt.Println(user)
}

func double_pointer_modify(ptr **int) {
	**ptr = 50
}

// Modify the original value with Double Pointer
func DoublePointer() {
	println("------> Double Pointer")
	num := 10
	ptr := &num
	ptr2 := &ptr // pointer to a pointer

	double_pointer_modify(ptr2)
	fmt.Println(num)
}

func arr_slice_modify(arr *[3]int) {
	(*arr)[0] = 99
}

func ArraySliceModify() {
	nums := [3]int{1, 2, 3}
	arr_slice_modify(&nums)
	fmt.Println(nums)
}
