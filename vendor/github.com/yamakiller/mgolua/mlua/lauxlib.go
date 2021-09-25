package mlua

//#include <lua.h>
//#include <lauxlib.h>
//#include <lualib.h>
//#include <stdlib.h>
//#include "mgolua.h"
import "C"
import (
	"fmt"
	"unsafe"
)

// LuaError : 错误信息
type LuaError struct {
	code       int
	message    string
	stackTrace []LuaStackEntry
}

// LuaReg : 注册函数表
type LuaReg struct {
	Name string
	Func LuaGoFunction
}

// LuaDebug : 调试栈
type LuaDebug struct {
	Event           int
	Name            string
	NameWhat        string
	What            string
	Source          string
	CurrentLine     int
	LineDefined     int
	LastLineDefined int
	Nups            uint8
	NParams         uint8
	IsVararg        byte
	IsTailCall      byte
	ShortSrc        []byte
}

// LuaBuffer :
type LuaBuffer = C.luaL_Buffer

// Error : 错误信息
func (err *LuaError) Error() string {
	return err.message
}

// GetCode : 错误码
func (err *LuaError) GetCode() int {
	return err.code
}

// GetStackTrace : 错误的栈信息
func (err *LuaError) GetStackTrace() []LuaStackEntry {
	return err.stackTrace
}

// LoadFile : luaL_loadfile
func (L *State) LoadFile(filename string) int {
	Cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(Cfilename))
	return int(C.mlua_loadfile(L._s, Cfilename))
}

// DoFile Executes file, returns nil for no errors or the lua error string on failure
func (L *State) DoFile(filename string) error {
	if r := L.LoadFile(filename); r != 0 {
		return &LuaError{r, L.ToString(-1), L.StackTrace()}
	}
	return L.Call(0, LUAMULTRET)
}

// LoadString : luaL_loadstring
func (L *State) LoadString(s string) int {
	Cs := C.CString(s)
	defer C.free(unsafe.Pointer(Cs))
	return int(C.luaL_loadstring(L._s, Cs))
}

// DoString : luaL_dostring
func (L *State) DoString(s string) error {
	if r := L.LoadString(s); r != 0 {
		return &LuaError{r, L.ToString(-1), L.StackTrace()}
	}
	return L.Call(0, LUAMULTRET)
}

// LoadBuffer : luaL_loadbuffer
func (L *State) LoadBuffer(data *byte, sz uint, name string) int {
	Cname := C.CString(name)
	defer C.free(unsafe.Pointer(Cname))
	return int(C.mlua_loadbuffer(L._s, (*C.char)((unsafe.Pointer)(data)), C.size_t(sz), Cname))
}

// Argcheck : luaL_argcheck
// WARNING: before b30b2c62c6712c6683a9d22ff0abfa54c8267863 the function ArgCheck had the opposite behaviour
func (L *State) Argcheck(cond bool, narg int, extramsg string) {
	if !cond {
		Cextramsg := C.CString(extramsg)
		defer C.free(unsafe.Pointer(Cextramsg))
		C.luaL_argerror(L._s, C.int(narg), Cextramsg)
	}
}

// ArgError : luaL_argerror
func (L *State) ArgError(narg int, extramsg string) int {
	Cextramsg := C.CString(extramsg)
	defer C.free(unsafe.Pointer(Cextramsg))
	return int(C.luaL_argerror(L._s, C.int(narg), Cextramsg))
}

// CallMeta : luaL_callmeta
func (L *State) CallMeta(obj int, e string) int {
	Ce := C.CString(e)
	defer C.free(unsafe.Pointer(Ce))
	return int(C.luaL_callmeta(L._s, C.int(obj), Ce))
}

// IsLightUserdata : Returns true if the value at index is light user data
func (L *State) IsLightUserdata(index int) bool {
	return LuaValType(C.lua_type(L._s, C.int(index))) == LUATLIGHTUSERDATA
}

// IsBoolean : lua_isboolean
func (L *State) IsBoolean(index int) bool {
	return LuaValType(C.lua_type(L._s, C.int(index))) == LUATBOOLEAN
}

// IsNil lua_isnil
func (L *State) IsNil(index int) bool { return LuaValType(C.lua_type(L._s, C.int(index))) == LUATNIL }

// IsNone lua_isnone
func (L *State) IsNone(index int) bool { return LuaValType(C.lua_type(L._s, C.int(index))) == LUATNONE }

// IsNoneOrNil : lua_isnoneornil
func (L *State) IsNoneOrNil(index int) bool { return int(C.lua_type(L._s, C.int(index))) <= 0 }

// IsNumber : lua_isnumber
func (L *State) IsNumber(index int) bool { return C.lua_isnumber(L._s, C.int(index)) == 1 }

// IsString : lua_isstring
func (L *State) IsString(index int) bool { return C.lua_isstring(L._s, C.int(index)) == 1 }

// IsGFunction : lua_iscfunction -> LuaGFunction
func (L *State) IsGFunction(index int) bool {
	return C.lua_iscfunction(L._s, C.int(index)) == 1
}

// IsFunction : lua_function -> IsFunction
func (L *State) IsFunction(index int) bool {
	return LuaValType(C.lua_type(L._s, C.int(index))) == LUATFUNCTION
}

// IsTable : lua_istable
func (L *State) IsTable(index int) bool {
	return LuaValType(C.lua_type(L._s, C.int(index))) == LUATTABLE
}

// IsThread : lua_isthread
func (L *State) IsThread(index int) bool {
	return LuaValType(C.lua_type(L._s, C.int(index))) == LUATTHREAD
}

// IsUserdata : lua_isuserdata
func (L *State) IsUserdata(index int) bool { return C.lua_isuserdata(L._s, C.int(index)) == 1 }

// IsGoStruct : mlua_isgostruct
func (L *State) IsGoStruct(index int) bool {
	id := uint(C.mlua_isgostruct(L._s, C.int(index)))
	if id == 0 {
		return false
	}

	return true
}

// NewTable : lua_newtable
func (L *State) NewTable() {
	C.lua_createtable(L._s, 0, 0)
}

// NewUserData : lua_newuserdata
func (L *State) NewUserData(sz uint) unsafe.Pointer {
	return unsafe.Pointer(C.lua_newuserdata(L._s, C.size_t(sz)))
}

// NewThread : lua_newthread
func (L *State) NewThread() *State { //TODO: should have same lists as parent
	//		but may complicate gc
	s := C.lua_newthread(L._s)
	return &State{s, nil}
}

// Next : lua_next
func (L *State) Next(index int) int {
	return int(C.lua_next(L._s, C.int(index)))
}

// CheckAny : luaL_checkany
func (L *State) CheckAny(narg int) {
	C.luaL_checkany(L._s, C.int(narg))
}

// CheckInteger : luaL_checkinteger
func (L *State) CheckInteger(narg int) int {
	return int(C.luaL_checkinteger(L._s, C.int(narg)))
}

// CheckNumber : luaL_checknumber
func (L *State) CheckNumber(narg int) float64 {
	return float64(C.luaL_checknumber(L._s, C.int(narg)))
}

// CheckString : luaL_checkstring
func (L *State) CheckString(narg int) string {
	var length C.size_t
	return C.GoString(C.luaL_checklstring(L._s, C.int(narg), &length))
}

// CheckType : luaL_checktype
func (L *State) CheckType(narg int, t LuaValType) {
	C.luaL_checktype(L._s, C.int(narg), C.int(t))
}

// CheckStack : luaL_checkstack
func (L *State) CheckStack(sz int, msg string) {
	Cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(Cmsg))
	C.luaL_checkstack(L._s, C.int(sz), Cmsg)
}

// TestUData : luaL_testudata
func (L *State) TestUData(narg int, tname string) unsafe.Pointer {
	Ctname := C.CString(tname)
	defer C.free(unsafe.Pointer(Ctname))
	return unsafe.Pointer(C.luaL_testudata(L._s, C.int(narg), Ctname))
}

// CheckUdata : luaL_checkudata
func (L *State) CheckUdata(narg int, tname string) unsafe.Pointer {
	Ctname := C.CString(tname)
	defer C.free(unsafe.Pointer(Ctname))
	return unsafe.Pointer(C.luaL_checkudata(L._s, C.int(narg), Ctname))
}

// Len : luaL_len
func (L *State) Len(index int) int {
	return int(C.luaL_len(L._s, C.int(index)))
}

// GSub : luaL_gsub
func (L *State) GSub(s string, p string, r string) string {
	Cs := C.CString(s)
	Cp := C.CString(p)
	Cr := C.CString(r)

	defer func() {
		C.free(unsafe.Pointer(Cs))
		C.free(unsafe.Pointer(Cp))
		C.free(unsafe.Pointer(Cr))
	}()

	return C.GoString(C.luaL_gsub(L._s, Cs, Cp, Cr))
}

// GetSubTable : luaL_getsubtable
func (L *State) GetSubTable(index int, fname string) int {
	Cfname := C.CString(fname)
	defer C.free(unsafe.Pointer(Cfname))
	return int(C.luaL_getsubtable(L._s, C.int(index), Cfname))
}

// TraceBack : luaL_traceback
func (L *State) TraceBack(L1 *State, msg string, level int) {
	Cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(Cmsg))
	C.luaL_traceback(L._s, L1._s, Cmsg, C.int(level))
}

// GetMetaField : luaL_getmetafield
func (L *State) GetMetaField(obj int, e string) bool {
	Ce := C.CString(e)
	defer C.free(unsafe.Pointer(Ce))
	return C.luaL_getmetafield(L._s, C.int(obj), Ce) != 0
}

// NewMetaTable : luaL_newmetatable
func (L *State) NewMetaTable(tname string) bool {
	Ctname := C.CString(tname)
	defer C.free(unsafe.Pointer(Ctname))
	return C.luaL_newmetatable(L._s, Ctname) != 0
}

// SetMetatable : luaL_setmetatable
func (L *State) SetMetatable(tname string) {
	Ctname := C.CString(tname)
	defer C.free(unsafe.Pointer(Ctname))
	C.luaL_setmetatable(L._s, Ctname)
}

// GetMetatable :
func (L *State) GetMetatable(tname string) int {
	Ctname := C.CString(tname)
	defer C.free(unsafe.Pointer(Ctname))
	return int(C.mlua_getmetatable(L._s, Ctname))
}

// OptInteger : luaL_optinteger
func (L *State) OptInteger(narg int, d int) int {
	return int(C.luaL_optinteger(L._s, C.int(narg), C.lua_Integer(d)))
}

// OptNumber : luaL_optnumber
func (L *State) OptNumber(narg int, d float64) float64 {
	return float64(C.luaL_optnumber(L._s, C.int(narg), C.lua_Number(d)))
}

// OptString : luaL_optstring
func (L *State) OptString(narg int, d string) string {
	var length C.size_t
	Cd := C.CString(d)
	defer C.free(unsafe.Pointer(Cd))
	return C.GoString(C.luaL_optlstring(L._s, C.int(narg), Cd, &length))
}

// SetFuncs : luaL_setfuncs
func (L *State) SetFuncs(regs []LuaReg, nup int) {
	L.CheckStack(nup+1, "too many upvalues")
	for _, r := range regs {
		tmpcp := unsafe.Pointer(&r.Func)
		C.lua_pushlightuserdata(L._s, tmpcp)
		for i := 0; i < nup; i++ {
			L.PushValue(-(nup + 1))
		}

		C.mlua_push_fun_wrapper(L._s, C.int(nup+1))
		fmt.Println(L.Type(-(nup + 2)))
		L.SetField(-(nup + 2), r.Name)
	}
	L.Pop(nup)
}

// Ref : luaL_ref
func (L *State) Ref(t int) int {
	return int(C.luaL_ref(L._s, C.int(t)))
}

// LTypename : luaL_typename
func (L *State) LTypename(index int) string {
	return C.GoString(C.lua_typename(L._s, C.lua_type(L._s, C.int(index))))
}

// Unref : luaL_unref
func (L *State) Unref(t int, ref int) {
	C.luaL_unref(L._s, C.int(t), C.int(ref))
}

// Where : luaL_where
func (L *State) Where(lvl int) {
	C.luaL_where(L._s, C.int(lvl))
}

// NewState : luaL_newstate
func NewState() *State {
	ls := (C.luaL_newstate())
	if ls == nil {
		return nil
	}
	L := newState(ls)
	return L
}

// OpenLibs : luaL_openlibs
func (L *State) OpenLibs() {
	C.luaL_openlibs(L._s)
}

// BuffInit : luaL_buffinit
func (L *State) BuffInit(b *LuaBuffer) {
	C.luaL_buffinit(L._s, b)
}

// AddLString : luaL_addlstring
func (L *State) AddLString(b *LuaBuffer, s unsafe.Pointer, sz uint) {
	C.luaL_addlstring(b, (*C.char)(s), C.size_t(sz))
}

// AddString : luaL_addstring
func (L *State) AddString(b *LuaBuffer, s string) {
	Cs := C.CString(s)
	defer C.free(unsafe.Pointer(Cs))

	C.luaL_addstring(b, Cs)
}

// AddValue : luaL_addvalue
func (L *State) AddValue(b *LuaBuffer) {
	C.luaL_addvalue(b)
}

// PushResult : luaL_pushresult
func (L *State) PushResult(b *LuaBuffer) {
	C.luaL_pushresult(b)
}

// PushResultSize : luaL_pushresultsize
func (L *State) PushResultSize(b *LuaBuffer, sz uint) {
	C.luaL_pushresultsize(b, C.size_t(sz))
}

// BuffInitSize : luaL_buffinitsize
func (L *State) BuffInitSize(b *LuaBuffer, sz uint) unsafe.Pointer {
	return unsafe.Pointer(C.luaL_buffinitsize(L._s, b, C.size_t(sz)))
}

// PrepBuffSize : luaL_prepbuffsize
func (L *State) PrepBuffSize(b *LuaBuffer, sz uint) unsafe.Pointer {
	return unsafe.Pointer(C.luaL_prepbuffsize(b, C.size_t(sz)))
}

// PrepBuffer : luaL_prepbuffer
func (L *State) PrepBuffer(b *LuaBuffer) unsafe.Pointer {
	return unsafe.Pointer(C.luaL_prepbuffsize(b, C.size_t(C.mlua_buffersize())))
}
