package javaclass

import (
	"errors"
	"io"

	"github.com/MJKWoolnough/byteio"
)

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

func readElementValue(r io.Reader) (ElementValue, error) {
	br := byteio.BigEndianReader{r}
	tag, _, err := br.ReadUint8()
	if err != nil {
		return nil, err
	}
	var elementValue ElementValue
	switch tag {
	case EVByte, EVChar, EVDouble, EVFloat, EVInt, EVLong, EVShort, EVBoolean, EVString:
		var cvi ConstValueIndex
		cvi, err = readEVConstValueIndex(r)
		cvi.tag = tag
		elementValue = cvi
	case EVEnumConstant:
		elementValue, err = readEVEnumConstant(r)
	case EVClass:
		elementValue, err = readEVClass(r)
	case EVAnnotationType:
		elementValue, err = readEVAnnotationType(r)
	case EVArray:
		elementValue, err = readEVArray(r)
	default:
		err = ErrUnknownElementValueTag
	}
	if err != nil {
		return nil, err
	}
	return elementValue, nil
}

type ConstValueIndex struct {
	tag   uint8
	Index uint16
}

func readEVConstValueIndex(r io.Reader) (ConstValueIndex, error) {
	br := byteio.BigEndianReader{r}
	index, _, err := br.ReadUint16()
	return ConstValueIndex{
		Index: index,
	}, err
}

func (c ConstValueIndex) Tag() uint8 {
	return c.tag
}

type EnumConstValue struct {
	TypeNameIndex, ConstNameIndex uint16
}

func readEVEnumConstant(r io.Reader) (ElementValue, error) {
	br := byteio.BigEndianReader{r}
	typeNameIndex, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	constNameIndex, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	return EnumConstValue{
		TypeNameIndex:  typeNameIndex,
		ConstNameIndex: constNameIndex,
	}, nil
}

func (EnumConstValue) Tag() uint8 {
	return EVEnumConstant
}

type ClassInfoIndex struct {
	Index uint16
}

func readEVClass(r io.Reader) (ElementValue, error) {
	br := byteio.BigEndianReader{r}
	index, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	return ClassInfoIndex{index}, nil
}

func (ClassInfoIndex) Tag() uint8 {
	return EVClass
}

type AnnotationValue struct {
	Annotation Annotation
}

func readEVAnnotationType(r io.Reader) (ElementValue, error) {
	annotation, err := readAnnotation(r)
	if err != nil {
		return nil, err
	}
	return AnnotationValue{annotation}, nil
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

func readEVArray(r io.Reader) (ElementValue, error) {
	br := byteio.BigEndianReader{r}
	numValues, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	arrayValues := make([]ElementValue, numValues)
	for i := uint16(0); i < numValues; i++ {
		arrayValues[i], err = readElementValue(r)
		if err != nil {
			return nil, err
		}

	}
	return ArrayValue{arrayValues}, nil
}

//Errors

var ErrUnknownElementValueTag = errors.New("unknown element value tag")
