package testparrot

// type RecordField struct {
// 	name     string
// 	recorder KVRecorder
// }

// func (f *RecordField) init(name string, recorder KVRecorder) {
// 	f.name = name
// 	f.recorder = recorder
// }

// func (f *RecordField) Record(val interface{}) interface{} {
// 	return f.recorder.Record(f.name, val)
// }

// func (f *RecordField) Value() interface{} {
// 	return f.recorder.Value(f.name)
// }

type StructRecorder struct {
	KVRecorder

	// self holds pointer for a struct, which contain fields used for testing
	self Recorder
}

func (r *StructRecorder) Init(self Recorder) {
	r.KVRecorder.Init(self)

	r.self = self
	r.name = getStructName(r.self)
}

func (r *StructRecorder) RecordField(fieldPtr interface{}, value interface{}) interface{} {
	fieldName := getFieldName(r.self, fieldPtr)

	value = r.Record(fieldName, value)

	if err := setStructValue(r.self, fieldName, value); err != nil {
		panic(err)
	}

	return value
}
