package testdata

import (
	"github.com/metatexx/avrox"
	"time"
)

import _ "embed"

//go:generate avscgen -ns "testing" . TestStruct
//go:embed test_struct.avsc
var TestStructAVSC string

type TestPtrSubStruct struct {
	FieldSubString  string
	FieldPtrInt8    *int8
	FieldPtrFloat64 *float64
	FieldPtrFloat32 *float32
}

type TestSubStruct struct {
	FieldSubString string
}

type TestStruct struct {
	Magic                  avrox.Magic
	FieldString            string
	FieldSubStruct         TestSubStruct
	FieldPtrSubStruct      *TestPtrSubStruct
	FieldPtrSubStringNil   *TestPtrSubStruct
	FieldSlicePtrSubStruct []*TestPtrSubStruct
	FieldTime              time.Time
	FieldDate              time.Time `avsc:"type:int,logicalType:date"`
	FieldInt               int
}
