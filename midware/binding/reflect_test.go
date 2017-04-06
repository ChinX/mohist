package binding

import (
	"log"
	"reflect"
	"testing"
)

type MyStruct struct {
	name string
}

func (m *MyStruct) GetName() string {
	return m.name
}

type IStruct interface {
	GetName()string
}

func TestReflectTypeAndValue(t *testing.T) {
	s := "this is string"
	log.Println(reflect.TypeOf(s)) 			// output: "string"
	log.Println(reflect.ValueOf(s)) 		// output: "this is string"

	x := float64(3.4)
	log.Println(reflect.TypeOf(x))  		// output: "float64"
	log.Println(reflect.ValueOf(x))  		// output: 3.4

	a := MyStruct{name:"this is my name"}
	typ := reflect.TypeOf(a)
	// reflect.Type Base struct
	log.Println(typ) 				// output: "binding.MyStruct"
	log.Println(typ.NumMethod())                   	// output: 0
	//log.Println(typ.Method(0))                    // panic: reflect: Method index out of range
	log.Println(typ.Name())                        	// output: MyStruct
	log.Println(typ.PkgPath())                     	// output: "github.com/chinx/mohist/midware/binding"
	log.Println(typ.Size())                        	// output: 8
	log.Println(typ.String())                      	// output: "binding.MyStruct"
	//log.Println(typ.Elem().String())              // panic: reflect: Elem of invalid type
	//log.Println(typ.Elem().FieldByIndex([]int{0}))// panic: reflect: Elem of invalid type
	//log.Println(typ.Elem().FieldByName("name"))   // panic: reflect: Elem of invalid type
	//log.Println(typ.Elem().NumField()) 		// panic: reflect: Elem of invalid type

	log.Println(typ.Kind() == reflect.Ptr)                              	// output: false
	//log.Println(typ.Elem().Kind() == reflect.Struct)                  	// panic: reflect: Elem of invalid type
	//log.Println(typ.Implements(reflect.TypeOf((*IStruct)(nil)).Elem())) 	// panic: reflect: Elem of invalid type

	pTyp := reflect.TypeOf(&a)
	log.Println(pTyp) 				// output: "binding.MyStruct"
	log.Println(pTyp.NumMethod())                   // output: 1
	log.Println(pTyp.Method(0))                     // output: {GetName  func(*main.MyStruct) string <func(*main.MyStruct) string Value> 0}
	log.Println(pTyp.Name())                        // output: ""
	log.Println(pTyp.PkgPath())                     // output: ""
	log.Println(pTyp.Size())                        // output: 8
	log.Println(pTyp.String())                      // output: *main.MyStruct
	log.Println(pTyp.Elem().String())               // output: main.MyStruct
	log.Println(pTyp.Elem().FieldByIndex([]int{0})) // output: {name github.com/chinx/mohist/midware/binding string  0 [0] false}
	log.Println(pTyp.Elem().FieldByName("name"))    // output: {name github.com/chinx/mohist/midware/binding string  0 [0] false} true

	log.Println(pTyp.Kind() == reflect.Ptr)                              	// output: true
	log.Println(pTyp.Elem().Kind() == reflect.Struct)                    	// output: true
	log.Println(pTyp.Implements(reflect.TypeOf((*IStruct)(nil)).Elem())) 	// output: true
	log.Println(pTyp.Elem().NumField()) 		// output: 1

	log.Println(reflect.ValueOf(a)) 		// output: {this is my name}
	log.Println(reflect.ValueOf(&a)) 		// output: &{this is my name}

	// MethodByName, Call
	//b := reflect.ValueOf(a).MethodByName("GetName").Call([]reflect.Value{})// panic: reflect: call of reflect.Value.Call on zero Value
	b := reflect.ValueOf(&a).MethodByName("GetName").Call([]reflect.Value{})
	log.Println(b[0])  				// output: this is my name


	cha := make(chan int)
	log.Println(reflect.TypeOf(cha).ChanDir()) 	// output: chan

	var fun func(x int, y ...float64) string
	var fun2 func(x int, y float64) string
	log.Println(reflect.TypeOf(fun).IsVariadic())  	// output: true
	log.Println(reflect.TypeOf(fun2).IsVariadic()) 	// output: false
	log.Println(reflect.TypeOf(fun).In(0))         	// output: int
	log.Println(reflect.TypeOf(fun).In(1))         	// output: []float64
	log.Println(reflect.TypeOf(fun).NumIn())       	// output: 2
	log.Println(reflect.TypeOf(fun).NumOut())      	// output: 1
	log.Println(reflect.TypeOf(fun).Out(0))        	// output: string

	mp := make(map[string]int)
	mp["test1"] = 1
	log.Println(reflect.TypeOf(mp).Key()) 		// output: string

	arr := [1]string{"test"}
	log.Println(reflect.TypeOf(arr).Len()) 		// output: 1


}

func TestReflectCanSet(t *testing.T) {
	var a MyStruct
	b := &MyStruct{}
	log.Println(reflect.ValueOf(a))			// output: {}
	log.Println(reflect.ValueOf(b))			// output: &{}

	a.name = "this is my name"
	b.name = "this is my name"
	val := reflect.ValueOf(a).FieldByName("name")
	//val := reflect.ValueOf(b).FieldByName("name")	//panic: reflect: call of reflect.Value.FieldByName on ptr Value
	////指针的ValueOf返回的是指针的Type，它是没有Field的，所以也就不能使用FieldByName
	log.Println(val)				// output: "this is my name"

	log.Println(reflect.ValueOf(a).FieldByName("name").CanSet())	// output: false
	log.Println(reflect.ValueOf(&(a.name)).Elem().CanSet())		// output: true
	////CanSet当Value是可寻址的时候，返回true，否则返回false

	c := "this is my other name"
	p := reflect.ValueOf(&c)
	log.Println(p.CanSet())				// output: false
	log.Println(p.Elem().CanSet())			// output: true
	////CanSet是一个指针的时候（p）它是不可寻址的，但是当是p.Elem()(实际上就是*p)，它就是可以寻址的
	p.Elem().SetString("newName")
	log.Println(c)					// output: "newName"

	p2 := reflect.ValueOf(c)
	log.Println(p2.CanSet())			// output: false
	////CanSet是一个指针的时候（p）它是不可寻址的，但是当是p.Elem()(实际上就是*p)，它就是可以寻址的
	//p2.SetString("otherName")			// panic: reflect: reflect.Value.SetString using unaddressable value
	log.Println(c)					// output: "newName"
}

func TestReflectSimply(t *testing.T)  {
	x := 3.14
	v := reflect.ValueOf(x)
	log.Println("type of v:", v.Type(), "settability of v:", v.CanSet()) // output: type of v: float64 settability of v: false

	p := reflect.ValueOf(&x)
	log.Println("type of p:", p.Type(), "settability of p:", p.CanSet()) // output: type of p: *float64 settability of p: false

	e := p.Elem()
	log.Println("type of e:", e.Type(), "settability of e:", e.CanSet()) // output: type of e: float64 settability of e: true
	e.SetFloat(2.75)

	log.Println(x)
}

func TestReflectStruct(t *testing.T)  {
	type T struct {
		A int
		B string
	}

	at := T{23,"hello world"}
	s := reflect.ValueOf(&at).Elem()
	typeOfT := s.Type()
	for i:=0; i<s.NumField(); i++ {
		f := s.Field(i)
		log.Printf("%d: %s %s = %v\n", i,typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
	s.Field(0).SetInt(22)
	s.Field(1).SetString("XXOO")
	log.Println(at)
}
