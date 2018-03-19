package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"strings"
	"time"
)

type Vertex struct {
	X int
	Y int
}

func mapTests() {
	var mapTest map[string]int = make(map[string]int) // map a string to an int
	mapTest["a"] = 2
	mapTest["b"] = 3

	mapLiteral := map[string]float64{
		"pi":  3.14159,
		"one": 1.0,
	}

	fmt.Println(mapTest["a"])
	fmt.Println(mapTest["A"]) // Not found (case sensitive, so it returns the default (0))
	fmt.Println(mapLiteral["pi"])

	delete(mapLiteral, "pi")

	element, ok := mapLiteral["pi"]
	fmt.Println("The value:", element, "Present?", ok)

	element, ok = mapLiteral["one"]
	fmt.Println("The value:", element, "Present?", ok)
}

func structTests() {
	vert := Vertex{22, 33}
	fmt.Println(Vertex{22, 33})
	fmt.Println(vert.X)

	vertPointer := &vert
	fmt.Println(vertPointer.X) // Just like rust, we don't need to do (*vertPointer).x or have a C++ -> operator.
}

func arrayTests() {
	var anArray [22]int
	fmt.Println(anArray[12]) // May be initialized to default upon creation here, but I'm not going to count on it.

	bArray := [4]int{1, 2, 3, 4} // If you use [] instead of [4], it returns a slice instead of the underlying array.
	fmt.Println(bArray[2])

	slice := bArray[0:3] // Elements 0-2, same as Rust. It's just a reference, so you can modify it to modify the underlying array.
	fmt.Printf("Slice Type: %T\n", slice)

	// You can omit slice bounds, same as Rust.

	// Slice extending is rather odd, so I just grabbed the example directoy to mess with (https://tour.golang.org/moretypes/11)
	s := []int{2, 3, 5, 7, 11, 13}
	printSlice(s)

	// Slice the slice to give it zero length.
	s = s[:0]
	printSlice(s)

	// Extend its length.
	s = s[:4]
	printSlice(s)

	// Drop its first two values.
	s = s[2:]
	printSlice(s)

	// Extend and drop again
	s = s[:4] // Can only extend up to 4. You cannot extend 'backwards' AFAIK, only forwards.
	printSlice(s)

	// Drop its first two values.
	s = s[2:]
	printSlice(s)

	// Dynamically-sized arrays are with slices,
	dynaSlice := make([]int, 5) // 5-element zeroed array with length 5
	printSlice(dynaSlice)

	// Extend with append.
	dynaSlice = append(dynaSlice, 22)
	printSlice(dynaSlice)

	// strings test
	fmt.Println(strings.Join([]string{"a", "b", "c"}, "-"))

	// using range
	for idx, val := range dynaSlice {
		fmt.Printf("%d = %d\n", idx, val)
	}

	// Definitely useful: https://tour.golang.org/moretypes/18
	// Demonstrates you can modify slice elements while iterating through them
}

// https://tour.golang.org/moretypes/11
func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func pointerTests() {
	i := 22

	k := &i
	m := &k

	fmt.Println(m)
	fmt.Println(k)
	fmt.Println(&m)
	fmt.Println(*m)
	fmt.Println(**m)
	fmt.Println(*k)

	// *m = 44 Doesn't work as we're converting a pointer to a number.
	**m = 55
	fmt.Println(*k)
}

func moreTests() {
	defer fmt.Println("Defer 1")
	defer fmt.Println("Defer 2")

	switch os := runtime.GOOS; os {
	case "windows":
		fmt.Println("Microsoft Windows")
	case "linux":
		fmt.Println("Linux.")
	default:
		fmt.Printf("%s\n", os)
	}

	// Switch also supports generic comparisons, speeding up if-then-else chains, only if we don't provide a value.
	switch rand.Intn(10) {
	case 3:
		fmt.Println("It was a 3")
	case 7:
		fmt.Println("It was a 7")
	default:
		fmt.Println("It was not a 3")
	}

	// defer is interesting, it will be called now, in reverse order.
	fmt.Println("Defers: ")
}

func addMultMod(x int, y int, z int, w int) int {
	return ((x + y) * z) % w
}

func addMultModSep(x int, y int) (add int, mul int, mod int) {
	add = x + y
	mul = x * y
	mod = x % y
	return
}

var testValue, otherValue int = 2, 4

const Pi = 3.141592653589 // This is untyped, it takes the type needed by its context

const Big = 1 << 100
const Small = Big >> 99

type AFloat float64 // We can also define methods on these, basically you can only define methods on objects in the same package.

// Vertex is read-only in this context. It also is COPIED, so we realistically want a pointer type.
// We also don't want to intermix them on the same type
func (v Vertex) LengthSqd() float64 {
	return float64(v.X*v.X + v.Y*v.Y)
}

func (v *Vertex) Double() int {
	if v == nil { // Equivalent of null
		fmt.Println("Hit a null...")
		return 2
	}

	v.X *= 2
	v.Y *= 2
	return v.X + v.Y
}

// interfaces are just like C#, a set of methods
type AnInterface interface {
	Double() int
	LengthSqd() float64
}

func (v *Vertex) String() string {
	return fmt.Sprintf("[%v, %v]", v.X, v.Y)
}

func doWork(x, y int, channel chan int) {
	z := x + y
	channel <- z
}

func threadTests() {
	// Run on separate goroutines
	go fmt.Println("hi a b c d e f g h i j k l m n o p")
	go fmt.Println("there")

	// Pass data through channels
	c := make(chan int)
	go doWork(2, 7, c)
	result := <-c
	fmt.Println(result)

	ch := make(chan int, 2)
	ch <- 1
	ch <- 2 // Won't block, as we're buffered
	// ch <- 3 // Will block (channel is full), deadlocking
	fmt.Println(<-ch)
	fmt.Println(<-ch)

	close(ch) // As per docs, only necessary when the receiver *must* be told there are no more values coming.
	next, ok := <-ch
	fmt.Println("The value:", next, "Channel Alive?", ok)

	// See https://tour.golang.org/concurrency/5 for a very cool use of select to wait concurrently on multiple communication ops.

	// Mutexes are fairly standard: https://tour.golang.org/concurrency/9 and demonstrate using defer usefully.
}

func objectTests() {

	// Just like rust, you define methods on types and not classes.
	v := Vertex{2, 4}
	fmt.Println(v.LengthSqd())

	v.Double()      // Auto-interpreted as (&v).Double() or Double(&v) for us.
	fmt.Println(&v) // Interestingly, this uses our custom .String() method, but if we pass in v directly, we use the default.

	// var myVar AnInterface
	// myVar can hold Vertex* only if we implement LengthSqd for Vertex*, and Vertex only if we implement it for the raw data.
	// This is a rather interesting way of defining an interface, but makes sense.
	// No explicit declarations needed for interfaces, which is rather nice as well.

	// Methods can be called with null objects, which is rather cool, because then the METHOD gets to handle the potential null pointer exception

	// var myVar2 interface{} // Holds ANY type whatsoever. Internally, tuple of (value, type)

	// Type assertions, copied from https://tour.golang.org/methods/15 for reference
	var i interface{} = "hello"

	s := i.(string)
	fmt.Println(s)

	s, ok := i.(string)
	fmt.Println(s, ok)

	f, ok := i.(float64)
	fmt.Println(f, ok)

	// Can also switch on types
	// Error handling is just another tuple -- return an object (with .Error() returning a string) on error or nil on non error
	// Readers operate just as normal -- read into a slice, returns the bytes read and/or error
}

func takeMult(givenFunction func(float64, float64) float64) float64 {
	return givenFunction(2.0, 3.0)
}

func doubler() func() int {
	loc := 1 // Each doubler owns its own internal variable.
	return func() int {
		loc += loc
		return loc
	}
}

func closureTests() {
	first, second := doubler(), doubler()
	for i := 0; i < 10; i++ {
		fmt.Println(first())
		fmt.Println(second())
	}
}

func main() {
	threadTests()
	time.Sleep(5 * time.Second)

	objectTests()
	time.Sleep(5 * time.Second)

	// Functions can be passed along, like values
	multFunc := func(ax, ay float64) float64 {
		return ax * ay
	}

	fmt.Println(takeMult(multFunc))

	closureTests()
	mapTests()
	time.Sleep(5 * time.Second)

	structTests()
	arrayTests()
	time.Sleep(5 * time.Second)

	pointerTests()
	moreTests()
	time.Sleep(5 * time.Second)

	fmt.Println("Println with stuff", rand.Intn(10))
	fmt.Printf("Have a note %g %g\n", math.Sqrt(7), math.Pi)
	fmt.Println(addMultMod(1, 2, 3, 4))

	a, b, c := addMultModSep(1, 2)
	fmt.Printf("Add: %d. Mult: %d. Mod: %d.\n", a, b, c)

	var x int = 7
	y := 7 // Same idea, but implicit type
	var intValue, boolValue = 3, false
	var aValue, bValue = 'c', "stringValue"

	// basic types: https://tour.golang.org/basics/11
	// Nothing too out of the ordinary, recommends using int instead of a sized type unless needed

	var z float64 = float64(x) // C++ style explicit conversions for everything
	fmt.Printf("x type: %T. y type: %T, ivType: %T, bvType: %T, aType: %T, bType: %T. Z: %f\n", x, y, intValue, boolValue, aValue, bValue, z)
	// fmt.Println(int(Big)) // This will not build as it overflows on the integer conversion
	fmt.Println(int(Small))

	// The only loop in go is the for
	for i := 0; i < 100; i++ { // := to define the variable i. Could also set i = 0 if i is defined elsewhere.
		if i%3 == 0 {
			fmt.Println(i % 10)
		}
	}

	// Interestingly, if statemets can have a pre-execute condition
	if v := math.Pow(2, 2); v == 4 {
		fmt.Println(v) // v is in scope, only till end if the if-else.
	} else {
		fmt.Println(v + 1) // Never going to be called.
	}
}
