package testparrot

// R defines a global test recorder
var R = NewRecorder()

// export ValToPtr util method
var ValToPtr = valToPtr

var Record = R.Record

var RecordNext = R.RecordNext
