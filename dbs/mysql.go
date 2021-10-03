package dbs

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	//import mysql base library
	_ "github.com/go-sql-driver/mysql"
	"github.com/yamakiller/magicLibs/util"
)

//MySQLValue doc
//@Summary mysql result value
//@Struct MySQLValue
//@Member interface{} value
//@Member reflect.Type value type
type MySQLValue struct {
	_v interface{}
	_t reflect.Type
}

//Print doc
//@Summary Print value informat
//@Method Print
func (slf *MySQLValue) Print() {
	fmt.Printf("%+v, %+v\n", slf._v, slf._t)
}

//IsEmpty doc
//@Summary Return Is Empty
//@Method IsEmpty
//@Return (bool) emtpy: true  no empty:false
func (slf *MySQLValue) IsEmpty() bool {
	if slf._v == nil {
		return true
	}
	return false
}

//ToString doc
//@Summary Return string value
//@Method ToString
//@Return (string) a string
func (slf *MySQLValue) ToString() string {
	return slf.getString()
}

//ToUint doc
//@Summary Return uint value
//@Method ToUint
//@Return (uint) a value
func (slf *MySQLValue) ToUint() uint {
	if v, e := slf.getNumber(); e == nil {
		return uint(v)
	}

	v, e := strconv.Atoi(slf.getString())
	if e != nil {
		return 0
	}

	return uint(v)
}

//ToInt doc
//@Summary Return int value
//@Method ToInt
//@Return (int) a value
func (slf *MySQLValue) ToInt() int {
	if v, e := slf.getNumber(); e == nil {
		return int(v)
	}

	v, e := strconv.Atoi(slf.getString())
	if e != nil {
		return 0
	}

	return v
}

//ToUint32 doc
//@Summary Return uint32 value
//@Method ToUint32
//@Return (uint32) a value
func (slf *MySQLValue) ToUint32() uint32 {
	if v, e := slf.getNumber(); e == nil {
		return uint32(v)
	}

	v, e := strconv.Atoi(slf.getString())
	if e != nil {
		return 0
	}

	return uint32(v)
}

//ToInt32 doc
//@Summary Return int32 value
//@Method ToInt32
//@Return (int32) a value
func (slf *MySQLValue) ToInt32() int32 {
	if v, e := slf.getNumber(); e == nil {
		return int32(v)
	}
	v, e := strconv.Atoi(slf.getString())
	if e != nil {
		return 0
	}

	return int32(v)
}

//ToUint64 doc
//@Summary Return uint64 value
//@Method ToUint64
//@Return (uint64) a value
func (slf *MySQLValue) ToUint64() uint64 {
	if v, e := slf.getNumber(); e == nil {
		return uint64(v)
	}

	v, e := strconv.ParseInt(slf.getString(), 10, 64)
	if e != nil {
		return 0
	}
	return uint64(v)
}

//ToInt64 doc
//@Summary Return int64 value
//@Method ToInt64
//@Return (int64) a value
func (slf *MySQLValue) ToInt64() int64 {
	if v, e := slf.getNumber(); e == nil {
		return int64(v)
	}

	v, e := strconv.ParseInt(slf.getString(), 10, 64)
	if e != nil {
		return 0
	}
	return v
}

//ToFloat doc
//@Summary Return float32 value
//@Method ToFloat
//@Return (float32) a value
func (slf *MySQLValue) ToFloat() float32 {
	if v, e := slf.getFloat(); e == nil {
		return float32(v)
	}

	v, e := strconv.ParseFloat(slf.getString(), 32)
	if e != nil {
		return 0.0
	}
	return float32(v)
}

//ToDouble doc
//@Summary Return float64 value
//@Method ToDouble
//@Return (float64) a value
func (slf *MySQLValue) ToDouble() float64 {
	if v, e := slf.getFloat(); e == nil {
		return v
	}

	v, e := strconv.ParseFloat(slf.getString(), 64)
	if e != nil {
		return 0.0
	}
	return v
}

//ToByte doc
//@Summary Return []byte value
//@Method ToByte
//@Return ([]byte) a value
func (slf *MySQLValue) ToByte() []byte {
	return ([]byte)(slf._v.([]uint8))
}

//ToTimeStamp doc
//@Summary  Return  time int64 value
//@Method ToTimeStamp
//@Return (int64) a value
func (slf *MySQLValue) ToTimeStamp() int64 {
	v := slf.ToDateTime()
	if v == nil {
		return 0
	}

	return v.Unix()
}

//ToDate doc
//@Summary  Return  time date value
//@Method ToDate
//@Return (*time.Time) a value
func (slf *MySQLValue) ToDate() *time.Time {
	v, e := time.Parse("2006-01-02", slf.getString())
	if e != nil {
		return nil
	}

	return &v
}

//ToDateTime doc
//@Summary  Return  time date value
//@Method ToDateTime
//@Return (*time.Time) a value
func (slf *MySQLValue) ToDateTime() *time.Time {
	v, e := time.Parse("2006-01-02 15:04:05", slf.getString())
	if e != nil {
		return nil
	}

	return &v
}

func (slf *MySQLValue) getString() string {
	return string(slf._v.([]uint8))
}

func (slf *MySQLValue) getNumber() (int64, error) {

	switch slf._t.Kind() {
	case reflect.Int64:
		return slf._v.(int64), nil
	case reflect.Int32:
		return int64(slf._v.(int32)), nil
	case reflect.Int16:
		return int64(slf._v.(int16)), nil
	default:
		return 0, fmt.Errorf("error: not number type")
	}
}

func (slf *MySQLValue) getFloat() (float64, error) {
	switch slf._t.Kind() {
	case reflect.Float32:
		return float64(slf._v.(float32)), nil
	case reflect.Float64:
		return slf._v.(float64), nil
	default:
		return 0, errors.New("error: not float type")
	}
}

//MySQLReader doc
//@Summary mysql reader
//@Struct MySQLReader
//@Member (int) count row of number
//@Member (int) read current row in index
//@Member ([]string) columns name
//@Member ([]MySQLValue) a mysql value
type MySQLReader struct {
	_rows       int
	_currentRow int
	_columns    []string
	_data       []MySQLValue
}

//GetAsNameValue doc
//@Summary Return column name to value
//@Method GetAsNameValue
//@Return (*MySQLValue) mysql value
//@Return (error) error informat
func (slf *MySQLReader) GetAsNameValue(name string) (*MySQLValue, error) {
	idx := slf.getNamePos(name)
	if idx == -1 {
		return nil, fmt.Errorf("mysql column %s is does not exist", name)
	}

	return slf.GetValue(idx)
}

//GetValue doc
//@Summary Return column index to value
//@Method GetValue
//@Return (*MySQLValue) mysql value
//@Return (error) error informat
func (slf *MySQLReader) GetValue(idx int) (*MySQLValue, error) {
	rpos := (slf._currentRow * len(slf._columns)) + idx
	if rpos >= len(slf._data) {
		return nil, fmt.Errorf("mysql column %d overload", idx)
	}

	return &slf._data[rpos], nil
}

//GetTryValue doc
//@Summary Return column index to value
//@Method GetTryValue
//@Return (*MySQLValue) mysql value
//@Return (error) error informat
func (slf *MySQLReader) GetTryValue(idx int) *MySQLValue {
	r, e := slf.GetValue(idx)
	if e != nil {
		return nil
	}

	return r
}

//Next doc
//@Summary Return Next row is success
//@Method Next
//@Return (bool) Next Success: true Next Fail:false
func (slf *MySQLReader) Next() bool {
	if (slf._currentRow + 1) >= slf._rows {
		return false
	}
	slf._currentRow++
	return true
}

//GetColumn doc
//@Summary Return Column of number
//@Method GetColumn
//@Return (int) a number
func (slf *MySQLReader) GetColumn() int {
	return len(slf._columns)
}

//GetRow doc
//@Summary Return row of number
//@Method GetRow
//@Return (int) a number
func (slf *MySQLReader) GetRow() int {
	return int(slf._rows)
}

//Rest doc
//@Summary go to first row
//@Method Rest
func (slf *MySQLReader) Rest() {
	slf._currentRow = -1
}

//Close doc
//@Summary clear data
//@Method Close
func (slf *MySQLReader) Close() {
	slf._columns = nil
	slf._data = nil
}

func (slf *MySQLReader) getNamePos(name string) int {
	for i, v := range slf._columns {
		if v == name {
			return i
		}
	}
	return -1
}

//MySQLDB doc
//@Summary mysql operation object
//@Struct MySQLDB
//@Member (string) mysql connection dsn
//@Member (*sql.DB) mysql connection object
type MySQLDB struct {
	_dsn string
	_db  *sql.DB
}

//Initial doc
//@Summary initialization mysql DB
//@Method Initial
//@Param (string) mysql connection dsn
//@Param (int) mysql connection max of number
//@Param (int) mysql connection idle of number
//@Param (int) mysql connection life time[util/sec]
//@Return (error) fail:return error, success: return nil
func (slf *MySQLDB) Initial(dsn string, maxConn, maxIdleConn, lifeSec int) error {
	var err error
	slf._db, err = sql.Open("mysql", dsn)
	util.AssertEmpty(slf._db, fmt.Sprintf("mysql open fail:%+v", err))
	slf._db.SetMaxOpenConns(maxConn)
	slf._db.SetMaxIdleConns(maxIdleConn)
	slf._db.SetConnMaxLifetime(time.Duration(lifeSec) * time.Second)
	slf._dsn = dsn

	err = slf._db.Ping()
	if err != nil {
		return err
	}

	return nil
}

//Query doc
//@Summary execute sql query
//@Method Query
//@Param (string) query sql
//@Param (...interface{}) sql params
//@Return (map[string]interface{}) query result
//@Return (error) fail: return error, success: return nil
func (slf *MySQLDB) Query(strSQL string, args ...interface{}) (*MySQLReader, error) {
	if perr := slf._db.Ping(); perr != nil {
		return nil, perr
	}

	rows, err := slf._db.Query(strSQL, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := &MySQLReader{_currentRow: -1}
	record._columns = columns
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		d := make([]MySQLValue, len(columns))
		for i, col := range values {
			d[i]._v = col
			d[i]._t = reflect.TypeOf(col)
		}
		record._data = append(record._data, d...)
		record._rows++
	}

	return record, nil
}

// QueryPage doc
// @Summary Query page data
// @Method  Query
// @Param   (string) table files (xxx,xxx)
// @Param   (string) table names (xxx,xxx)
// @Param   (string) query condition
// @Param   (string) query order mode
// @Param   (int) page
// @Param   (int) pageSize
// @Param   (...interface{}) where args
// @Return  (int) pageCount
// @Return  (*dbs.MySQLReader) reader
// @Return  (error) ree
func (slf *MySQLDB) QueryPage(fileds, tables, where, order string, page, pageSize int, args ...interface{}) (pageCount int, record *MySQLReader, err error) {
	if perr := slf._db.Ping(); perr != nil {
		return 0, nil, perr
	}

	idxFiled := strings.Split(fileds, ",")
	sqlWhere := ""
	sqlOrder := ""

	if where != "" {
		sqlWhere = fmt.Sprintf(" WHERE %s", where)
	}

	if order != "" {
		sqlOrder = fmt.Sprintf(" Order By %s", order)
	}

	strSQL := fmt.Sprintf("SELECT count(%s) as totalNum FROM %s%s%s", strings.Trim(idxFiled[0], " "), tables, sqlWhere, sqlOrder)

	var rows *sql.Rows
	rows, err = slf._db.Query(strSQL, args...)
	if err != nil {
		return 0, nil, err
	}

	for rows.Next() {
		err = rows.Scan(&pageCount)
		if err != nil {
			rows.Close()
			return 0, nil, err
		}
	}
	rows.Close()

	strSQL = fmt.Sprintf("SELECT %s as totalNum FROM %s%s%s LIMIT %d,%d", strings.Trim(idxFiled[0], " "),
		tables, sqlWhere, sqlOrder, page, pageSize)
	rows, err = slf._db.Query(strSQL, args...)
	if err != nil {
		return 0, nil, err
	}

	defer rows.Close()
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	for j := range values {
		scanArgs[j] = &values[j]
	}

	record = &MySQLReader{_currentRow: -1}
	record._columns = columns
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		d := make([]MySQLValue, len(columns))
		for i, col := range values {
			d[i]._v = col
			d[i]._t = reflect.TypeOf(col)
		}
		record._data = append(record._data, d...)
		record._rows++
	}
	return pageCount, record, nil
}

//Insert doc
//@Summary execute sql Insert
//@Method Insert
//@Param (string) Insert sql
//@Param (...interface{}) sql params
//@Return (int54) insert of number
//@Return (error) fail: return error, success: return nil
func (slf *MySQLDB) Insert(strSQL string, args ...interface{}) (int64, error) {
	if perr := slf._db.Ping(); perr != nil {
		return 0, perr
	}

	r, err := slf._db.Exec(strSQL, args...)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return r.LastInsertId()
}

//Update doc
//@Summary execute sql Update
//@Method Update
//@Param (string) Update sql
//@Param (...interface{}) sql params
//@Return (int54) Update of number
func (slf *MySQLDB) Update(strSQL string, args ...interface{}) (int64, error) {
	if perr := slf._db.Ping(); perr != nil {
		return 0, perr
	}

	r, err := slf._db.Exec(strSQL, args...)
	if err != nil {
		return 0, err
	}

	return r.RowsAffected()
}

//Close doc
//@Summary close mysql connection
//@Method CLose
func (slf *MySQLDB) Close() {
	slf._db.Close()
}
