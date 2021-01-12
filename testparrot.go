package testparrot

// R defines a global test recorder
var R = NewRecorder()

// export ValToPtr util method
var VarToPtr = valToPtr

var Record = R.Record

var RecordNext = R.RecordNext
