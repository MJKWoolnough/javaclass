package javaclass

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

type ConstantValueAttribute struct {
	ConstantValue uint16
}

func (ConstantValueAttribute) Name() string {
	return AttrConstantValue
}

type Exception struct {
	StartPC, EndPC, HandlerPC, CatchType uint16
}

type CodeAttribute struct {
	MaxStack, MaxLocals uint16
	Code                []byte
	ExceptionTable      []Exception
	Attributes          []AttributeInfo
}

func (CodeAttribute) Name() string {
	return AttrCode
}

type StackMapTableAttribute struct {
	Entries []StackMapFrame
}

func (StackMapTableAttribute) Name() string {
	return AttrStackMapTable
}

type ExceptionsAttribute struct {
	ExceptionIndexTable []uint16
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

func (InnerClassesAttribute) Name() string {
	return AttrInnerClasses
}

type EnclosingMethodAttribute struct {
	ClassIndex, MethodIndex uint16
}

func (EnclosingMethodAttribute) Name() string {
	return AttrEnclosingMethod
}

type SyntheticAttribute struct {
	SignatureIndex uint16
}

func (SyntheticAttribute) Name() string {
	return AttrSynthetic
}

type SourceFileAttribute struct {
	SourceFileIndex uint16
}

func (SourceFileAttribute) Name() string {
	return AttrSourceFile
}

type SourceFileDebugAttribute struct {
	DebugExtension []uint8
}

func (SourceFileDebugAttribute) Name() string {
	return AttrSourceDebugExtension
}

type LineNumber struct {
	StartPC, LineNumber uint16
}

type LineNumberTableAttribute struct {
	LineNumberTable []LineNumber
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

func (LocalVariableTableAttribute) Name() string {
	return AttrLocalVariableTable
}

type LocalVariableType struct {
	StartPC, Length, NameIndex, SignatureIndex, Index uint16
}

type LocalVariableTypeTableAttribute struct {
	LocalVariableTypeTable []LocalVariableType
}

func (LocalVariableTypeTableAttribute) Name() string {
	return AttrLocalVariableTypeTable
}

type DeprecatedAttribute struct{}

func (DeprecatedAttribute) Name() string {
	return AttrDeprecated
}

type ElementValuePair struct {
	ElementNameIndex uint16
	Value            ElementValue
}

type Annotation struct {
	ElementValuePairs []ElementValuePair
}

type RuntimeVisibleAnnotationsAttribute struct {
	Annotations []Annotation
}

func (RuntimeVisibleAnnotationsAttribute) Name() string {
	return AttrRuntimeVisibleAnnotations
}

type RuntimeInvisibleAnnotationsAttribute struct {
	Annotations []Annotation
}

func (RuntimeInvisibleAnnotationsAttribute) Name() string {
	return AttrRuntimeInvisibleAnnotations
}

type ParameterAnnotation struct {
	Annotations []Annotation
}

type RuntimeVisibleParameterAnnotationsAttribute struct {
	ParameterAnnotations []ParameterAnnotation
}

func (RuntimeVisibleParameterAnnotationsAttribute) Name() string {
	return AttrRuntimeVisibleParameterAnnotations
}

type RuntimeInvisibleParameterAnnotationsAttribute struct {
	ParameterAnnotations []ParameterAnnotation
}

func (RuntimeInvisibleParameterAnnotationsAttribute) Name() string {
	return AttrRuntimeInvisibleParameterAnnotations
}

type AnnotationDefaultAttribute struct {
	DefaultValue ElementValue
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

func (BootstrapMethodsAttribute) Name() string {
	return AttrBootstrapMethods
}
