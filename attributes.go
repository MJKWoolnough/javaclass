package javaclass

import (
	"errors"
	"io"
	"io/ioutil"

	"github.com/MJKWoolnough/byteio"
)

const (
	AttrConstantValue                        = "ConstantValue"
	AttrCode                                 = "Code"
	AttrStackMapTable                        = "StackMapTable"
	AttrExceptions                           = "Exceptions"
	AttrInnerClasses                         = "InnerClasses"
	AttrEnclosingMethod                      = "EnclosingMethod"
	AttrSynthetic                            = "Synthetic"
	AttrSignature                            = "Signature"
	AttrSourceFile                           = "SourceFile"
	AttrSourceDebugExtension                 = "SourceDebugExtension"
	AttrLineNumberTable                      = "LineNumberTable"
	AttrLocalVariableTable                   = "LocalVariableTable"
	AttrLocalVariableTypeTable               = "LocalVariableTypeTable"
	AttrDeprecated                           = "Deprecated"
	AttrRuntimeVisibleAnnotations            = "RuntimeVisibleAnnotations"
	AttrRuntimeInvisibleAnnotations          = "RuntimeInvisibleAnnotations"
	AttrRuntimeVisibleParameterAnnotations   = "RuntimeVisibleParameterAnnotations"
	AttrRuntimeInvisibleParameterAnnotations = "RuntimeInvisibleParameterAnnotations"
	AttrAnnotationDefault                    = "AnnotationDefault"
	AttrBootstrapMethods                     = "BootstrapMethods"
)

type AttributeInfo interface {
	Name() string
}

func (c *Class) readAttributes(r io.Reader) ([]AttributeInfo, error) {
	br := byteio.BigEndianReader{r}
	attributesCount, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	attributes := make([]AttributeInfo, 0, attributesCount)
	for i := uint16(0); i < attributesCount; i++ {
		ani, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		if int(ani) >= len(c.ConstantPool) {
			return nil, ErrInvalidConstantPoolIndex
		}
		cpi := c.ConstantPool[ani]
		if cpi.Type() != ConstantUTF8 {
			return nil, ErrInvalidConstantPoolType
		}
		cpUTF, ok := cpi.(ConstantUTF8Info)
		if !ok {
			return nil, ErrInvalidConstantPoolType
		}
		attributeLength, _, err := br.ReadUint32()
		if err != nil {
			return nil, err
		}
		lr := io.LimitReader(r, int64(attributeLength))
		var attributeInfo AttributeInfo
		switch cpUTF.String {
		case AttrConstantValue:
			attributeInfo, err = readConstantValue(lr)
		case AttrCode:
			attributeInfo, err = c.readCode(lr)
		case AttrStackMapTable:
			attributeInfo, err = readStackMapTable(lr)
		case AttrExceptions:
			attributeInfo, err = readExceptions(lr)
		case AttrInnerClasses:
			attributeInfo, err = readInnerClasses(lr)
		case AttrEnclosingMethod:
			attributeInfo, err = readEnclosingMethod(lr)
		case AttrSynthetic:
			attributeInfo, err = readSynthetic(lr)
		case AttrSignature:
			attributeInfo, err = readSignature(lr)
		case AttrSourceFile:
			attributeInfo, err = readSourceFile(lr)
		case AttrSourceDebugExtension:
			attributeInfo, err = readSourceDebugExtension(lr)
		case AttrLineNumberTable:
			attributeInfo, err = readLineNumberTable(lr)
		case AttrLocalVariableTable:
			attributeInfo, err = readLocalVariableTable(lr)
		case AttrLocalVariableTypeTable:
			attributeInfo, err = readLocalVariableTypeTable(lr)
		case AttrDeprecated:
			attributeInfo, err = readDeprecated(lr)
		case AttrRuntimeVisibleAnnotations:
			attributeInfo, err = readRuntimeVisibleAnnotations(lr)
		case AttrRuntimeInvisibleAnnotations:
			attributeInfo, err = readRuntimeInvisibleAnnotations(lr)
		case AttrRuntimeVisibleParameterAnnotations:
			attributeInfo, err = readRuntimeVisibleParameterAnnotations(lr)
		case AttrRuntimeInvisibleParameterAnnotations:
			attributeInfo, err = readRuntimeInvisibleParameterAnnotations(lr)
		case AttrAnnotationDefault:
			attributeInfo, err = readAnnotationDefault(lr)
		case AttrBootstrapMethods:
			attributeInfo, err = readBootstrapMethods(lr)
		default:
			return nil, ErrInvalidAttributeName
		}
		if err != nil {
			return nil, err
		}
		attributes = append(attributes, attributeInfo)
	}
	return attributes, nil
}

type ConstantValueAttribute struct {
	ConstantValue uint16
}

func readConstantValue(r io.Reader) (AttributeInfo, error) {
	br := byteio.BigEndianReader{r}
	n, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	return ConstantValueAttribute{n}, nil
}

func (ConstantValueAttribute) Name() string {
	return AttrConstantValue
}

type Exception struct {
	StartPC, EndPC, HandlerPC, CatchType uint16
}

func readExceptionsTable(r io.Reader) ([]Exception, error) {
	br := byteio.BigEndianReader{r}
	exceptionTableLength, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	exceptionsTable := make([]Exception, exceptionTableLength)
	for i := uint16(0); i < exceptionTableLength; i++ {
		startPC, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		endPC, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		handlerPC, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		catchType, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		exceptionsTable[i] = Exception{
			StartPC:   startPC,
			EndPC:     endPC,
			HandlerPC: handlerPC,
			CatchType: catchType,
		}
	}
	return exceptionsTable, nil
}

type CodeAttribute struct {
	MaxStack, MaxLocals uint16
	Code                []byte
	ExceptionTable      []Exception
	Attributes          []AttributeInfo
}

func (c *Class) readCode(r io.Reader) (AttributeInfo, error) {
	br := byteio.BigEndianReader{r}
	maxStack, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	maxLocals, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	codeLength, _, err := br.ReadUint32()
	if err != nil {
		return nil, err
	}
	code := make([]byte, codeLength)
	_, err = io.ReadFull(r, code)
	if err != nil {
		return nil, err
	}
	exceptions, err := readExceptionsTable(r)
	if err != nil {
		return nil, err
	}
	attributes, err := c.readAttributes(r)
	if err != nil {
		return nil, err
	}
	return &CodeAttribute{
		MaxStack:       maxStack,
		MaxLocals:      maxLocals,
		Code:           code,
		ExceptionTable: exceptions,
		Attributes:     attributes,
	}, nil
}

func (CodeAttribute) Name() string {
	return AttrCode
}

type StackMapTableAttribute struct {
	Entries []StackMapFrame
}

func readStackMapTable(r io.Reader) (AttributeInfo, error) {
	br := byteio.BigEndianReader{r}
	numEntries, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	entries := make([]StackMapFrame, numEntries)
	for i := uint16(0); i < numEntries; i++ {
		entries[i], err = readStackMapFrame(r)
		if err != nil {
			return nil, err
		}
	}
	return StackMapTableAttribute{entries}, nil
}

func (StackMapTableAttribute) Name() string {
	return AttrStackMapTable
}

type ExceptionsAttribute struct {
	ExceptionIndexTable []uint16
}

func readExceptions(r io.Reader) (AttributeInfo, error) {
	br := byteio.BigEndianReader{r}
	numExceptions, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	exceptions := make([]uint16, numExceptions)
	for i := uint16(0); i < numExceptions; i++ {
		exceptions[i], _, err = br.ReadUint16()
		if err != nil {
			return nil, err
		}
	}
	return ExceptionsAttribute{exceptions}, nil
}

func (ExceptionsAttribute) Name() string {
	return AttrExceptions
}

type ClassInfo struct {
	InnerClassInfoIndex, OuterClassInfoIndex, InnerClassNameIndex, InnerClassAccessFlags uint16
}

type InnerClassesAttribute struct {
	Classes []ClassInfo
}

func readInnerClasses(r io.Reader) (AttributeInfo, error) {
	br := byteio.BigEndianReader{r}
	numClasses, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	classes := make([]ClassInfo, numClasses)
	for i := uint16(0); i < numClasses; i++ {
		innerClassInfo, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		outerClassInfo, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		innerClassName, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		innerClassFlags, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		classes[i] = ClassInfo{
			InnerClassInfoIndex:   innerClassInfo,
			OuterClassInfoIndex:   outerClassInfo,
			InnerClassNameIndex:   innerClassName,
			InnerClassAccessFlags: innerClassFlags,
		}
	}
	return InnerClassesAttribute{classes}, nil
}

func (InnerClassesAttribute) Name() string {
	return AttrInnerClasses
}

type EnclosingMethodAttribute struct {
	ClassIndex, MethodIndex uint16
}

func readEnclosingMethod(r io.Reader) (AttributeInfo, error) {
	br := byteio.BigEndianReader{r}
	classIndex, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	methodIndex, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	return EnclosingMethodAttribute{
		ClassIndex:  classIndex,
		MethodIndex: methodIndex,
	}, nil
}

func (EnclosingMethodAttribute) Name() string {
	return AttrEnclosingMethod
}

type SyntheticAttribute struct{}

func readSynthetic(r io.Reader) (AttributeInfo, error) {
	return SyntheticAttribute{}, nil
}

func (SyntheticAttribute) Name() string {
	return AttrSynthetic
}

type SignatureAttribute struct {
	SignatureIndex uint16
}

func readSignature(r io.Reader) (AttributeInfo, error) {
	br := byteio.BigEndianReader{r}
	n, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	return SignatureAttribute{n}, nil
}

func (SignatureAttribute) Name() string {
	return AttrSignature
}

type SourceFileAttribute struct {
	SourceFileIndex uint16
}

func readSourceFile(r io.Reader) (AttributeInfo, error) {
	br := byteio.BigEndianReader{r}
	n, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	return SourceFileAttribute{n}, nil
}

func (SourceFileAttribute) Name() string {
	return AttrSourceFile
}

type SourceDebugAttribute struct {
	DebugExtension string
}

func readSourceDebugExtension(r io.Reader) (AttributeInfo, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return SourceDebugAttribute{string(data)}, nil
}

func (SourceDebugAttribute) Name() string {
	return AttrSourceDebugExtension
}

type LineNumber struct {
	StartPC, LineNumber uint16
}

type LineNumberTableAttribute struct {
	LineNumberTable []LineNumber
}

func readLineNumberTable(r io.Reader) (AttributeInfo, error) {
	br := byteio.BigEndianReader{r}
	lineNumberTableLength, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	lineNumberTable := make([]LineNumber, lineNumberTableLength)
	for i := uint16(0); i < lineNumberTableLength; i++ {
		startPC, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		lineNumber, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		lineNumberTable[i] = LineNumber{
			StartPC:    startPC,
			LineNumber: lineNumber,
		}
	}
	return LineNumberTableAttribute{lineNumberTable}, nil
}

func (LineNumberTableAttribute) Name() string {
	return AttrLineNumberTable
}

type LocalVariable struct {
	StartPC, Length, NameIndex, DescriptorIndex, Index uint16
}

type LocalVariableTableAttribute struct {
	LocalVariableTable []LocalVariable
}

func readLocalVariableTable(r io.Reader) (AttributeInfo, error) {
	br := byteio.BigEndianReader{r}
	localVariableTableLength, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	localVariableTable := make([]LocalVariable, localVariableTableLength)
	for i := uint16(0); i < localVariableTableLength; i++ {
		startPC, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		length, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		nameIndex, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		descriptorIndex, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		index, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		localVariableTable[i] = LocalVariable{
			StartPC:         startPC,
			Length:          length,
			NameIndex:       nameIndex,
			DescriptorIndex: descriptorIndex,
			Index:           index,
		}
	}
	return LocalVariableTableAttribute{localVariableTable}, nil
}

func (LocalVariableTableAttribute) Name() string {
	return AttrLocalVariableTable
}

type LocalVariableType struct {
	StartPC, Length, NameIndex, SignatureIndex, Index uint16
}

type LocalVariableTypeTableAttribute struct {
	LocalVariableTypeTable []LocalVariableType
}

func readLocalVariableTypeTable(r io.Reader) (AttributeInfo, error) {
	br := byteio.BigEndianReader{r}
	localVariableTypeTableLength, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	localVariableTypeTable := make([]LocalVariableType, localVariableTypeTableLength)
	for i := uint16(0); i < localVariableTypeTableLength; i++ {
		startPC, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		length, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		nameIndex, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		signatureIndex, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		index, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		localVariableTypeTable[i] = LocalVariableType{
			StartPC:        startPC,
			Length:         length,
			NameIndex:      nameIndex,
			SignatureIndex: signatureIndex,
			Index:          index,
		}
	}
	return LocalVariableTypeTableAttribute{localVariableTypeTable}, nil
}

func (LocalVariableTypeTableAttribute) Name() string {
	return AttrLocalVariableTypeTable
}

type DeprecatedAttribute struct{}

func readDeprecated(r io.Reader) (AttributeInfo, error) {
	return nil, nil
}

func (DeprecatedAttribute) Name() string {
	return AttrDeprecated
}

type ElementValuePair struct {
	ElementNameIndex uint16
	Value            ElementValue
}

func readElementValuePairs(r io.Reader) ([]ElementValuePair, error) {
	br := byteio.BigEndianReader{r}
	numElementValuePairs, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	elementValuePairs := make([]ElementValuePair, numElementValuePairs)
	for i := uint16(0); i < numElementValuePairs; i++ {
		elementNameIndex, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		elementValue, err := readElementValue(r)
		if err != nil {
			return nil, err
		}
		elementValuePairs[i] = ElementValuePair{
			ElementNameIndex: elementNameIndex,
			Value:            elementValue,
		}
	}
	return elementValuePairs, nil
}

type Annotation struct {
	TypeIndex         uint16
	ElementValuePairs []ElementValuePair
}

func readAnnotation(r io.Reader) (Annotation, error) {
	br := byteio.BigEndianReader{r}
	typeIndex, _, err := br.ReadUint16()
	if err != nil {
		return Annotation{}, err
	}
	elementValuePairs, err := readElementValuePairs(r)
	if err != nil {
		return Annotation{}, err
	}
	return Annotation{
		TypeIndex:         typeIndex,
		ElementValuePairs: elementValuePairs,
	}, nil
}

func readAnnotations(r io.Reader) ([]Annotation, error) {
	br := byteio.BigEndianReader{r}
	numAnnotations, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	annotations := make([]Annotation, numAnnotations)
	for i := uint16(0); i < numAnnotations; i++ {
		annotations[i], err = readAnnotation(r)
		if err != nil {
			return nil, err
		}
	}
	return annotations, nil
}

type RuntimeVisibleAnnotationsAttribute struct {
	Annotations []Annotation
}

func readRuntimeVisibleAnnotations(r io.Reader) (AttributeInfo, error) {
	annotations, err := readAnnotations(r)
	if err != nil {
		return nil, err
	}
	return RuntimeVisibleAnnotationsAttribute{annotations}, nil
}

func (RuntimeVisibleAnnotationsAttribute) Name() string {
	return AttrRuntimeVisibleAnnotations
}

type RuntimeInvisibleAnnotationsAttribute struct {
	Annotations []Annotation
}

func readRuntimeInvisibleAnnotations(r io.Reader) (AttributeInfo, error) {
	annotations, err := readAnnotations(r)
	if err != nil {
		return nil, err
	}
	return RuntimeInvisibleAnnotationsAttribute{annotations}, nil
}

func (RuntimeInvisibleAnnotationsAttribute) Name() string {
	return AttrRuntimeInvisibleAnnotations
}

type ParameterAnnotation struct {
	Annotations []Annotation
}

func readParameterAnnotations(r io.Reader) ([]ParameterAnnotation, error) {
	br := byteio.BigEndianReader{r}
	numAnnotations, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	parameterAnnotations := make([]ParameterAnnotation, numAnnotations)
	for i := uint16(0); i < numAnnotations; i++ {
		annotations, err := readAnnotations(r)
		if err != nil {
			return nil, err
		}
		parameterAnnotations[i] = ParameterAnnotation{annotations}
	}
	return parameterAnnotations, nil
}

type RuntimeVisibleParameterAnnotationsAttribute struct {
	ParameterAnnotations []ParameterAnnotation
}

func readRuntimeVisibleParameterAnnotations(r io.Reader) (AttributeInfo, error) {
	parameterAnnotations, err := readParameterAnnotations(r)
	if err != nil {
		return nil, err
	}
	return RuntimeVisibleParameterAnnotationsAttribute{parameterAnnotations}, nil
}

func (RuntimeVisibleParameterAnnotationsAttribute) Name() string {
	return AttrRuntimeVisibleParameterAnnotations
}

type RuntimeInvisibleParameterAnnotationsAttribute struct {
	ParameterAnnotations []ParameterAnnotation
}

func readRuntimeInvisibleParameterAnnotations(r io.Reader) (AttributeInfo, error) {
	parameterAnnotations, err := readParameterAnnotations(r)
	if err != nil {
		return nil, err
	}
	return RuntimeInvisibleParameterAnnotationsAttribute{parameterAnnotations}, nil
}

func (RuntimeInvisibleParameterAnnotationsAttribute) Name() string {
	return AttrRuntimeInvisibleParameterAnnotations
}

type AnnotationDefaultAttribute struct {
	DefaultValue ElementValue
}

func readAnnotationDefault(r io.Reader) (AttributeInfo, error) {
	return nil, nil
}

func (AnnotationDefaultAttribute) Name() string {
	return AttrAnnotationDefault
}

type BootstrapMethod struct {
	BootstrapMethodRef uint16
	BootstrapArguments []uint16
}

type BootstrapMethodsAttribute struct {
	BootstrapMethods []BootstrapMethod
}

func readBootstrapMethods(r io.Reader) (AttributeInfo, error) {
	return nil, nil
}

func (BootstrapMethodsAttribute) Name() string {
	return AttrBootstrapMethods
}

//Errors

var (
	ErrInvalidConstantPoolIndex = errors.New("invalid constant pool index")
	ErrInvalidConstantPoolType  = errors.New("invalid constant pool type")
	ErrInvalidAttributeName     = errors.New("invalid attribute name")
)
