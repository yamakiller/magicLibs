package js

import (
	"github.com/robertkrimen/otto"
	"github.com/yamakiller/magicLibs/script/stack"
)

//Bundle desc
//@method Bundle desc: Basic binding js library, extended here
//@param (*stack.JSStack)
func Bundle(stack *stack.JSStack) {
	stack.SetFunc("Refer", refer)
}

func refer(js otto.FunctionCall) otto.Value {
	switch js.Argument(0).String() {
	case "runtime":
		return jsruntimeBundle(js)
	default:
		return otto.Value{}
	}
}
