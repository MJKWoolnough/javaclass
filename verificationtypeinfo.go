package javaclass

const (
	InfoTopVariable             = 0
	InfoIntegerVariableInfo     = 1
	InfoFloatVariable           = 2
	InfoDoubleVariable          = 3
	InfoLongVariable            = 4
	InfoNullVariable            = 5
	InfoUnitializedThisVariable = 6
	InfoObjectVariable          = 7
	InfoUnitializedVariable     = 8
)

type VerificationTypeInfo interface {
	Tag() int
}

type TopVariableInfo struct{}

func (TopVariableInfo) Tag() int {
	return InfoTopVariable
}

type IntegerVariableInfo struct{}

func (IntegerVariableInfo) Tag() int {
	return InfoIntegerVariableInfo
}

type FloatVariableInfo struct{}

func (FloatVariableInfo) Tag() int {
	return InfoFloatVariable
}

type DoubleVariableInfo struct{}

func (DoubleVariableInfo) Tag() int {
	return InfoDoubleVariable
}

type LongVariableInfo struct{}

func (LongVariableInfo) Tag() int {
	return InfoLongVariable
}

type NullVariableInfo struct{}

func (NullVariableInfo) Tag() int {
	return InfoNullVariable
}

type UninitializedThisVariableInfo struct{}

func (UninitializedThisVariableInfo) Tag() int {
	return InfoUnitializedThisVariable
}

type ObjectVariableInfo struct {
	CPoolIndex uint16
}

func (ObjectVariableInfo) Tag() int {
	return InfoObjectVariable
}

type UninitializedVariableInfo struct {
	Offset uint16
}

func (UninitializedVariableInfo) Tag() int {
	return InfoUnitializedVariable
}
