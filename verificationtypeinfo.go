package javaclass

import (
	"errors"
	"io"

	"vimagination.zapto.org/byteio"
)

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

func readVerificationTypeInfo(r io.Reader) (VerificationTypeInfo, error) {
	br := byteio.BigEndianReader{Reader: r}
	tag, _, err := br.ReadUint8()
	if err != nil {
		return nil, err
	}
	switch tag {
	case InfoTopVariable:
		return TopVariableInfo{}, nil
	case InfoIntegerVariableInfo:
		return IntegerVariableInfo{}, nil
	case InfoFloatVariable:
		return FloatVariableInfo{}, nil
	case InfoDoubleVariable:
		return DoubleVariableInfo{}, nil
	case InfoLongVariable:
		return LongVariableInfo{}, nil
	case InfoNullVariable:
		return NullVariableInfo{}, nil
	case InfoUnitializedThisVariable:
		return UninitializedThisVariableInfo{}, nil
	case InfoObjectVariable:
		cpoolIndex, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		return ObjectVariableInfo{cpoolIndex}, nil
	case InfoUnitializedVariable:
		offset, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		return UninitializedVariableInfo{offset}, nil
	default:
		return nil, ErrUnknownVerificationTypeTag
	}
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

//Errros

var ErrUnknownVerificationTypeTag = errors.New("unknown verification type tag")
