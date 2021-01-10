package testparrot

var R = NewRegistry()

// export ValToPtr util method
var VarToPtr = ValToPtr

// export global registry Load method
var Load = R.Load

// export global registry Register method
var Register = R.Register

// export global registry MustRegister method
var MustRegister = R.MustRegister

var EnableRecording = R.EnableRecording
