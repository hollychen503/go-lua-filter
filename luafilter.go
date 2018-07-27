package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

type Filter struct {
	state                    *lua.LState
	exceptionHandlerFunction *lua.LFunction
}

func NewFilter() *Filter {
	state := lua.NewState()
	luajson.Preload(state)

	filter := &Filter{
		state: state,
	}
	filter.exceptionHandlerFunction = state.NewFunction(
		filter.exceptionHandler)
	return filter
}

func (f *Filter) exceptionHandler(L *lua.LState) int {
	panic("exception in lua code")
	return 0
}

func (f *Filter) LoadScript(filename string) error {
	return f.state.DoFile(filename)
}

func (f *Filter) LoadScriptString(str string) error {
	return f.state.DoString(str)
}

func (f *Filter) ValidateScript() error {
	fn := f.state.GetGlobal("filter")
	if fn.Type() != lua.LTFunction {
		return errors.New("Function 'filter' not found")
	}
	return nil
}

func (f *Filter) ValidateEvent(event string) (bool, error) {
	fn := f.state.GetGlobal("filter")

	f.state.Push(fn.(*lua.LFunction))
	f.state.Push(lua.LString(event))

	// one argument and one return value
	err := f.state.PCall(1, 1, f.exceptionHandlerFunction)
	if err != nil {
		return false, err
	}

	top := f.state.GetTop()
	returnValue := f.state.Get(top)
	if returnValue.Type() != lua.LTBool {
		return false, errors.New("Invalid return value")
	}

	return lua.LVAsBool(returnValue), err
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func TestSimple() {
	const str = `
	local json = require("json")
	assert(type(json) == "table")
	assert(type(json.decode) == "function")
	assert(type(json.encode) == "function")
	assert(json.encode(true) == "true")
	assert(json.encode(1) == "1")
	assert(json.encode(-10) == "-10")
	assert(json.encode(nil) == "{}")
	local obj = {"a",1,"b",2,"c",3}
	local jsonStr = json.encode(obj)
	local jsonObj = json.decode(jsonStr)
	for i = 1, #obj do
		assert(obj[i] == jsonObj[i])
	end
	local obj = {name="Tim",number=12345}
	local jsonStr = json.encode(obj)
	local jsonObj = json.decode(jsonStr)
	assert(obj.name == jsonObj.name)
	assert(obj.number == jsonObj.number)
	local obj = {"a","b",what="c",[5]="asd"}
	local jsonStr = json.encode(obj)
	local jsonObj = json.decode(jsonStr)
	assert(obj[1] == jsonObj["1"])
	assert(obj[2] == jsonObj["2"])
	assert(obj.what == jsonObj["what"])
	assert(obj[5] == jsonObj["5"])
	assert(json.decode("null") == nil)
	assert(json.decode(json.encode({person={name = "tim",}})).person.name == "tim")
	local obj = {
		abc = 123,
		def = nil,
	}
	local obj2 = {
		obj = obj,
	}
	obj.obj2 = obj2
	assert(json.encode(obj) == nil)
	local a = {}
	for i=1, 5 do
		a[i] = i
	end
	assert(json.encode(a) == "[1,2,3,4,5]")
	`
	s := lua.NewState()
	luajson.Preload(s)
	if err := s.DoString(str); err != nil {
		panic(err)
	}
	log.Println("OK")
}

func main() {
	if len(os.Args) != 2 {
		println("provide filter script")
		return
	}

	bgn := time.Now().UnixNano()

	///Size of the callstack & registry is fixed for mainly performance. You can change the default size of the callstack & registry.
	lua.RegistrySize = 1024 * 20
	lua.CallStackSize = 10240

	/* L := lua.NewState(lua.Options{
		CallStackSize: 120,
		RegistrySize:  120*20,
	}) */

	event := "func (f *Filter) ValidateEvent(event string) (bool, error) { if"

	///

	dat, err := ioutil.ReadFile(os.Args[1])
	check(err)
	datStr := string(dat)

	imax := 1 //1000
	jmax := 1 //10000

	for i := 0; i < imax; i++ {
		filter := NewFilter()
		//err := filter.LoadScript(os.Args[1])
		err := filter.LoadScriptString(datStr)
		if err != nil {
			panic(err.Error())
		}

		err = filter.ValidateScript()
		if err != nil {
			panic(err.Error())
		}

		for j := 0; j < jmax; j++ {

			isValid, err := filter.ValidateEvent(event)
			if err != nil || !isValid {
				panic(err.Error())
			}
		}
		log.Println(i)
	}
	end := time.Now().UnixNano()
	diff := end - bgn

	fmt.Println("total req:", imax*jmax)
	fmt.Println("used ", diff, "nano seconds, ", (float64(diff) / (float64)(time.Second)), "seconds")
	fmt.Println(float64(imax*jmax)/(float64(diff)/(float64)(time.Second)), " req per second")

	TestSimple()
	/*
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			event := scanner.Text()
			isValid, err := filter.ValidateEvent(event)
			if err != nil {
				panic(err.Error())
			}

			if isValid {
				println(event)
			}
		}*/
}
