package testparrot

// R defines a global test recorder
var R = NewRecorder()

// export Ptr util method to get pointer of a value using reflection
var Ptr = valToPtr

// export Must util method, so we can use it with some methods like json.Unmarshal
var Must = must

var Record = R.Record

var RecordNext = R.RecordNext
