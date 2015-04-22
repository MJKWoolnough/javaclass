package javaclass

const (
	EVByte           = 'B'
	EVChar           = 'C'
	EVDouble         = 'D'
	EVFloat          = 'F'
	EVInt            = 'I'
	EVLong           = 'J'
	EVShort          = 'S'
	EVBoolean        = 'Z'
	EVString         = 's'
	EVEnumConstant   = 'e'
	EVClass          = 'c'
	EVAnnotationType = '@'
	EVArray          = '['
)

type ElementValue interface {
	Tag() uint8
}

type ConstValueIndex struct {
	tag   uint8
	Index uint16
}

func (c ConstValueIndex) Tag() uint8 {
	return c.tag
}

type EnumConstValue struct {
	TypeNameIndex, ConstNameIndex uint16
}

func (EnumConstValue) Tag() uint8 {
	return EVEnumConstant
}

type ClassInfoIndex struct {
	Index uint16
}

func (ClassInfoIndex) Tag() uint8 {
	return EVClass
}

type AnnotationValue struct {
	Annotation Annotation
}

func (AnnotationValue) Tag() uint8 {
	return EVAnnotationType
}

type ArrayValue struct {
	ArrayValues []ElementValue
}

func (ArrayValue) Tag() uint8 {
	return EVArray
}
