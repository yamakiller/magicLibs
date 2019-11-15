package script

import (
	"github.com/yamakiller/magicLibs/script/js"
	"github.com/yamakiller/magicLibs/script/stack"
)

//NewJSStack desc
//@method NewJSStack desc: Create a js virtual machine
//@return (*stack.JSStack)
func NewJSStack() *stack.JSStack {
	stack := stack.MakeJSStack()
	js.Bundle(stack)
	return stack
}
