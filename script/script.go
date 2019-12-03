package script

import (
	"github.com/yamakiller/magicLibs/script/js"
	"github.com/yamakiller/magicLibs/script/stack"
)

//NewJSStack desc
//@Method NewJSStack desc: Create a js virtual machine
//@Return (*stack.JSStack)
func NewJSStack() *stack.JSStack {
	stack := stack.MakeJSStack()
	js.Bundle(stack)
	return stack
}
