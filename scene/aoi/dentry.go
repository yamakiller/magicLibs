package aoi

import "github.com/yamakiller/magicLibs/math"

//DEntry doc
type DEntry struct {
	_xPrev *DEntry
	_xNext *DEntry
	_yPrev *DEntry
	_yNext *DEntry
	_vKey  uint32
	_vPos  math.Vector2
}
