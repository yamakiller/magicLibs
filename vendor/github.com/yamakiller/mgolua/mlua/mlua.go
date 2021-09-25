package mlua

/*
#cgo CFLAGS: -I ${SRCDIR}/lua
#cgo llua LDFLAGS: -llua
#cgo luaa LDFLAGS: -llua -lm -ldl
#cgo linux,!llua,!luaa LDFLAGS: -llua
#cgo darwin,!llua,!luaa LDFLAGS: -llua
#cgo freebsd,!luaa LDFLAGS: -llua
#cgo windows,!llua LDFLAGS: -L${SRCDIR} -llua -lmingwex -lmingw32

#include <lua.h>
#include <stdlib.h>

#include "mgolua.h"
*/
import "C"
import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"unsafe"
)

//LuaStackEntry :
type LuaStackEntry struct {
	name        string
	source      string
	shortSource string
	currentLine int
}

func newState(L *C.lua_State) *State {
	newstate := &State{L, nil}
	C.mlua_setgostate(L, C.uintptr_t(uintptr(unsafe.Pointer(newstate))))
	return newstate
}

// AbsIndex : lua_absindex
func (L *State) AbsIndex(index int) int {
	return int(C.lua_absindex(L._s, C.int(index)))
}

// Copy : lua_copy
func (L *State) Copy(fromindex int, toindex int) {
	C.lua_copy(L._s, C.int(fromindex), C.int(toindex))
}

// Type : lua_type
func (L *State) Type(index int) int {
	return int(C.lua_type(L._s, C.int(index)))
}

// TypeName : lua_typename
func (L *State) TypeName(tp int) string {
	return C.GoString(C.lua_typename(L._s, C.int(tp)))
}

// GetTop : lua_gettop
func (L *State) GetTop() int {
	return int(C.lua_gettop(L._s))
}

// SetTop : lua_settop
func (L *State) SetTop(index int) {
	C.lua_settop(L._s, C.int(index))
}

// Pop : lua_pop
func (L *State) Pop(n int) {
	C.lua_settop(L._s, C.int(-n-1))
}

// Insert : lua_insert
func (L *State) Insert(index int) {
	C.lua_rotate(L._s, C.int(index), C.int(1))
}

// Remove : lua_remove
func (L *State) Remove(index int) {
	C.lua_rotate(L._s, C.int(index), C.int(-1))
	L.Pop(1)
}

// Replace ： lua_replace
func (L *State) Replace(index int) {
	C.mlua_replace(L._s, C.int(index))
}

// PushBoolean : lua_pushboolean
func (L *State) PushBoolean(b bool) {
	var bint int
	if b {
		bint = 1
	} else {
		bint = 0
	}
	C.lua_pushboolean(L._s, C.int(bint))
}

// PushString : lua_pushstring
func (L *State) PushString(str string) {
	Cstr := C.CString(str)
	defer C.free(unsafe.Pointer(Cstr))
	C.lua_pushlstring(L._s, Cstr, C.size_t(len(str)))
}

// PushBytes :
func (L *State) PushBytes(b []byte) {
	C.lua_pushlstring(L._s, (*C.char)(unsafe.Pointer(&b[0])), C.size_t(len(b)))
}

// PushInteger : lua_pushinteger
func (L *State) PushInteger(n int64) {
	C.lua_pushinteger(L._s, C.lua_Integer(n))
}

// PushNil : lua_pushnil
func (L *State) PushNil() {
	C.lua_pushnil(L._s)
}

// PushNumber : lua_pushnumber
func (L *State) PushNumber(n float64) {
	C.lua_pushnumber(L._s, C.lua_Number(n))
}

// PushThread : lua_pushthread
func (L *State) PushThread() (isMain bool) {
	return C.lua_pushthread(L._s) != 0
}

// PushValue lua_pushvalue
func (L *State) PushValue(index int) {
	C.lua_pushvalue(L._s, C.int(index))
}

// PushGoFunction "" lua_pushcfunction -> PushGoFunction
func (L *State) PushGoFunction(f LuaGoFunction) {
	pf := unsafe.Pointer(&f)
	C.mlua_push_go_wrapper(L._s, pf, C.int(0))
}

// PushGoClosure : lua_pushcclosure -> PushGoClosure
func (L *State) PushGoClosure(f LuaGoFunction, args ...interface{}) {
	argsNum := 1
	pf := unsafe.Pointer(&f)
	C.lua_pushlightuserdata(L._s, pf)
	for _, val := range args {
		argsNum++
		switch reflect.TypeOf(val).Kind() {
		case reflect.Uint64:
			L.PushInteger(int64(val.(uint64)))
		case reflect.Uint32:
			L.PushInteger(int64(val.(uint32)))
		case reflect.Uint:
			L.PushInteger(int64(val.(uint)))
		case reflect.Int64:
			L.PushInteger(val.(int64))
		case reflect.Int32:
			L.PushInteger(int64(val.(int32)))
		case reflect.Int:
			L.PushInteger(int64(val.(int)))
			break
		case reflect.Float64:
			L.PushNumber(val.(float64))
		case reflect.Float32:
			L.PushNumber(float64(val.(float32)))
			break
		case reflect.String:
			L.PushString(reflect.ValueOf(val).String())
			break
		case reflect.Struct:
			L.PushUserGoStruct(val)
			break
		case reflect.Uintptr:
			L.PushLightGoStruct(unsafe.Pointer(val.(uintptr)))
		case reflect.UnsafePointer:
			L.PushLightGoStruct(unsafe.Pointer(reflect.ValueOf(val).Pointer()))
			break
		case reflect.Bool:
			L.PushBoolean(reflect.ValueOf(val).Bool())
			break
		default:
			panic(fmt.Sprintf("mlua go Closure %s Type not supported", reflect.TypeOf(val).Name()))
		}
	}
	C.mlua_push_fun_wrapper(L._s, C.int(argsNum))
}

// PushLiteral : mlua_pushliteral
func (L *State) PushLiteral(s string) {
	Cs := C.CString(s)
	defer C.free(unsafe.Pointer(Cs))
	C.mlua_pushliteral(L._s, Cs)
}

// UpvalueIndex : 闭包参数从第二个位置开始防蚊
func (L *State) UpvalueIndex(n int) int {
	return int(C.mlua_upvalueindex(C.int(n)))
}

// PushUserGoStruct : mlua_pushgostruct => lua_newuserdata
// 内存管理权交由 lua虚拟机管理
// 内存消耗略大
func (L *State) PushUserGoStruct(d interface{}) {
	var dby bytes.Buffer
	enc := gob.NewEncoder(&dby)
	err := enc.Encode(d)
	if err != nil {
		panic(err)
	}

	C.mlua_pushugostruct(L._s, (*C.char)(unsafe.Pointer(&dby.Bytes()[0])), C.size_t(len(dby.Bytes())))
}

// PushLightGoStruct : lua_pushlightuserdata =>nmlua_pushlgostruct
func (L *State) PushLightGoStruct(d unsafe.Pointer) {
	C.mlua_pushlgostruct(L._s, C.uintptr_t(uintptr(d)))
}

// SetGlobal : lua_setglobal
func (L *State) SetGlobal(name string) {
	Cname := C.CString(name)
	defer C.free(unsafe.Pointer(Cname))
	C.lua_setglobal(L._s, Cname)
}

// GetGlobal : lua_getglobal
func (L *State) GetGlobal(name string) {
	Cname := C.CString(name)
	defer C.free(unsafe.Pointer(Cname))
	C.lua_getglobal(L._s, Cname)
}

// ToString : lua_tostring
func (L *State) ToString(index int) string {
	var size C.size_t
	r := C.lua_tolstring(L._s, C.int(index), &size)
	return C.GoStringN(r, C.int(size))
}

// ToBytes : luaL_tolstring
func (L *State) ToBytes(index int) []byte {
	var size C.size_t
	b := C.lua_tolstring(L._s, C.int(index), &size)
	return C.GoBytes(unsafe.Pointer(b), C.int(size))
}

// ToInteger : ua_tointeger
func (L *State) ToInteger(index int) int64 {
	return int64(C.mlua_tointeger(L._s, C.int(index)))
}

// ToNumber : lua_tonumber
func (L *State) ToNumber(index int) float64 {
	return float64(C.mlua_tonumber(L._s, C.int(index)))
}

// ToBoolean : lua_toboolean
func (L *State) ToBoolean(index int) bool {
	return int(C.lua_toboolean(L._s, C.int(index))) != 0
}

//ToCheckString : luaL_checklstring
func (L *State) ToCheckString(index int) string {
	var size C.size_t
	r := C.luaL_checklstring(L._s, C.int(index), &size)
	return C.GoStringN(r, C.int(size))
}

// ToCheckInteger : luaL_checkinteger
func (L *State) ToCheckInteger(index int) int64 {
	return int64(C.luaL_checkinteger(L._s, C.int(index)))
}

// ToUserGoStruct lua_tougostruct => lua_touserdata
// 获取索引中的Go Struct 结构
// TODO： 思考感觉性能消耗不小!  没有没改进方案呢？
func (L *State) ToUserGoStruct(index int, s interface{}) {
	r := (*C.struct_GoStruct)(C.mlua_tougostruct(L._s, C.int(index)))
	n := int(r._sz)
	d := bytes.NewBuffer(C.GoBytes(unsafe.Pointer(&r._data[0]), C.int(n)))
	dec := gob.NewDecoder(d)
	err := dec.Decode(s)
	if err != nil {
		panic(err)
	}
}

// ToLightGoStruct :
func (L *State) ToLightGoStruct(index int) unsafe.Pointer {
	return unsafe.Pointer(C.mlua_tolgostruct(L._s, C.int(index)))
}

// RawLen : lua_rawlen
func (L *State) RawLen(index int) uint {
	return uint(C.lua_rawlen(L._s, C.int(index)))
}

// ToPointer : lua_topointer
func (L *State) ToPointer(index int) unsafe.Pointer {
	return unsafe.Pointer(C.lua_topointer(L._s, C.int(index)))
}

// RawEqual : lua_rawequal
func (L *State) RawEqual(index1 int, index2 int) int {
	return int(C.lua_rawequal(L._s, C.int(index1), C.int(index2)))
}

// GetTable : lua_gettable
func (L *State) GetTable(index int) int {
	return int(C.lua_gettable(L._s, C.int(index)))
}

// GetField : lua_getfield
func (L *State) GetField(index int, k string) int {
	Ck := C.CString(k)
	defer C.free(unsafe.Pointer(Ck))
	return int(C.lua_getfield(L._s, C.int(index), Ck))
}

// GetI : lua_geti
func (L *State) GetI(index int, n int64) int {
	return int(C.lua_geti(L._s, C.int(index), C.lua_Integer(n)))
}

// RawGet : lua_rawget
func (L *State) RawGet(index int) int {
	return int(C.lua_rawget(L._s, C.int(index)))
}

// RawGetI : lua_rawgeti
func (L *State) RawGetI(index int, n int64) int {
	return int(C.lua_rawgeti(L._s, C.int(index), C.lua_Integer(n)))
}

// RawGetP : lua_rawgetp
func (L *State) RawGetP(index int, p unsafe.Pointer) int {
	return int(C.lua_rawgetp(L._s, C.int(index), p))
}

// GetMetaTable : lua_getmetatable
func (L *State) GetMetaTable(objindex int) int {
	return int(C.lua_getmetatable(L._s, C.int(objindex)))
}

// GetUserValue lua_getuservalue
func (L *State) GetUserValue(index int) int {
	return int(C.lua_getuservalue(L._s, C.int(index)))
}

// SetTable lua_settable
func (L *State) SetTable(index int) {
	C.lua_settable(L._s, C.int(index))
}

// SetField : lua_setfield
func (L *State) SetField(index int, k string) {
	Ck := C.CString(k)
	defer C.free(unsafe.Pointer(Ck))
	C.lua_setfield(L._s, C.int(index), Ck)
}

// SetI : lua_seti
func (L *State) SetI(index int, n int64) {
	C.lua_seti(L._s, C.int(index), C.lua_Integer(n))
}

// RawSet : lua_rawset
func (L *State) RawSet(index int) {
	C.lua_rawset(L._s, C.int(index))
}

// RawSetI : lua_rawseti
func (L *State) RawSetI(index int, n int64) {
	C.lua_rawseti(L._s, C.int(index), C.lua_Integer(n))
}

// RawSetP :lua_rawsetp
func (L *State) RawSetP(index int, p unsafe.Pointer) {
	C.lua_rawsetp(L._s, C.int(index), p)
}

// SetMetaTable : lua_setmetatable
func (L *State) SetMetaTable(objindex int) {
	C.lua_setmetatable(L._s, C.int(objindex))
}

// SetUserValue : lua_setuservalue
func (L *State) SetUserValue(index int) {
	C.lua_setuservalue(L._s, C.int(index))
}

// Gc : lua_gc
func (L *State) Gc(what int, data int) int {
	return int(C.lua_gc(L._s, C.int(what), C.int(data)))
}

// Error : luaL_error
func (L *State) Error(sfmt string, v ...interface{}) int {
	Cerror := C.CString(fmt.Sprintf(sfmt, v...))
	defer C.free(unsafe.Pointer(Cerror))
	return int(C.mlua_error(L._s, Cerror))
}

// Concat : lua_concat
func (L *State) Concat(n int) {
	C.lua_concat(L._s, C.int(n))
}

// SetAllocF : lua_setallocf
func (L *State) SetAllocF(f LuaGoAllocFunction) {
	fp := unsafe.Pointer(&f)
	C.mlua_setallocf(L._s, fp)
}

// PushGlobalTable : lua_pushglobaltable
func (L *State) PushGlobalTable() {
	C.mlua_pushglobaltable(L._s)
}

// lua_pcall
func (L *State) pcall(nargs, nresults, errfunc int) int {
	return int(C.mlua_pcall(L._s, C.int(nargs), C.int(nresults), C.int(errfunc)))
}

func (L *State) callEx(nargs, nresults int, catch bool) (err error) {
	if catch {
		defer func() {
			if err2 := recover(); err2 != nil {
				if _, ok := err2.(error); ok {
					err = err2.(error)
				}
				return
			}
		}()
	}

	L.GetGlobal(C.GOLUA_PANIC_MSG_WARAPPER)
	erridx := L.GetTop() - nargs - 1
	L.Insert(erridx)
	r := L.pcall(nargs, nresults, erridx)
	L.Remove(erridx)
	if r != 0 {
		err = &LuaError{r, L.ToString(-1), L.StackTrace()}
		if !catch {
			return err
		}
	}
	return nil
}

// Call : lua_call
func (L *State) Call(nargs, nresults int) (err error) {
	return L.callEx(nargs, nresults, true)
}

// PCall : lua_pcall Extern
func (L *State) PCall(nargs, nresults, errfunc int) int {
	return L.pcall(nargs, nresults, errfunc)
}

// Register : Registers a Go function as a global variable
func (L *State) Register(name string, f LuaGoFunction) {
	L.PushGoFunction(f)
	L.SetGlobal(name)
}

// SetHook : lua_sethook
func (L *State) SetHook(f LuaGoHookFunction, mask int, count int) {
	L._h = &f
	C.mlua_sethook(L._s, C.int(mask), C.int(count))
}

// GetHookMask : lua_gethookmask
func (L *State) GetHookMask() int {
	return int(C.lua_gethookmask(L._s))
}

// GetHookCount : lua_gethookcount
func (L *State) GetHookCount() int {
	return int(C.lua_gethookcount(L._s))
}

// Gcc : lua_gc
func (L *State) Gcc(what int, data int) int {
	return int(C.lua_gc(L._s, C.int(what), C.int(data)))
}

// NewStateAlloc : exprot NewStateEnv
func NewStateAlloc(f LuaGoAllocFunction) *State {
	fp := unsafe.Pointer(&f)
	ls := C.mlua_newstate(fp)
	return newState(ls)
}

// Close : lua_close
func (L *State) Close() {
	C.lua_close(L._s)
}

// StackTrace :
func (L *State) StackTrace() []LuaStackEntry {
	r := []LuaStackEntry{}
	var d C.lua_Debug
	Sln := C.CString("Sln")
	defer C.free(unsafe.Pointer(Sln))

	for depth := 0; C.lua_getstack(L._s, C.int(depth), &d) > 0; depth++ {
		C.lua_getinfo(L._s, Sln, &d)
		ssb := make([]byte, C.LUA_IDSIZE)
		for i := 0; i < C.LUA_IDSIZE; i++ {
			ssb[i] = byte(d.short_src[i])
			if ssb[i] == 0 {
				ssb = ssb[:i]
				break
			}
		}
		ss := string(ssb)
		r = append(r, LuaStackEntry{C.GoString(d.name), C.GoString(d.source), ss, int(d.currentline)})
	}
	return r
}

// NewError :
func (L *State) NewError(msg string) *LuaError {
	return &LuaError{0, msg, L.StackTrace()}
}
