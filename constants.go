package javaclass

const (
	ConstantNull               = 0
	ConstantUTF8               = 1
	ConstantInteger            = 3
	ConstantFloat              = 4
	ConstantLong               = 5
	ConstantDouble             = 6
	ConstantClass              = 7
	ConstantString             = 8
	ConstantFieldRef           = 9
	ConstantMethodRef          = 10
	ConstantInterfaceMethodRef = 11
	ConstantNameAndType        = 12
	ConstantMethodHandle       = 15
	ConstantMethodType         = 16
	ConstantInvokeDynamic      = 18
)

type ConstantNullInfo struct{}

func (ConstantNullInfo) Type() int {
	return ConstantNull
}

type ConstantUTF8Info struct {
	String string
}

func (ConstantUTF8Info) Type() int {
	return ConstantUTF8
}

type ConstantIntegerInfo struct {
	Integer uint32
}

func (ConstantIntegerInfo) Type() int {
	return ConstantInteger
}

type ConstantFloatInfo struct {
	Float float32
}

func (ConstantFloatInfo) Type() int {
	return ConstantFloat
}

type ConstantLongInfo struct {
	Long uint64
}

func (ConstantLongInfo) Type() int {
	return ConstantLong
}

type ConstantDoubleInfo struct {
	Double float64
}

func (ConstantDoubleInfo) Type() int {
	return ConstantDouble
}

type ConstantClassInfo struct {
	NameIndex uint16
}

func (ConstantClassInfo) Type() int {
	return ConstantClass
}

type ConstantStringInfo struct {
	StringIndex uint16
}

func (ConstantStringInfo) Type() int {
	return ConstantString
}

type ConstantFieldRefInfo struct {
	ClassIndex, NameAndTypeIndex uint16
}

func (ConstantFieldRefInfo) Type() int {
	return ConstantFieldRef
}

type ConstantMethodRefInfo struct {
	ClassIndex, NameAndTypeIndex uint16
}

func (ConstantMethodRefInfo) Type() int {
	return ConstantMethodRef
}

type ConstantInterfaceMethodRefInfo struct {
	ClassIndex, NameAndTypeIndex uint16
}

func (ConstantInterfaceMethodRefInfo) Type() int {
	return ConstantInterfaceMethodRef
}

type ConstantNameAndTypeInfo struct {
	NameIndex, DescriptorIndex uint16
}

func (ConstantNameAndTypeInfo) Type() int {
	return ConstantNameAndType
}

type ConstantMethodHandleInfo struct {
	ReferenceKind  uint8
	ReferenceIndex uint16
}

func (ConstantMethodHandleInfo) Type() int {
	return ConstantMethodHandle
}

type ConstantMethodTypeInfo struct {
	DescriptorIndex uint16
}

func (ConstantMethodTypeInfo) Type() int {
	return ConstantMethodType
}

type ConstantInvokeDynamicInfo struct {
	BootstrapMethodAttrIndex, NameAndTypeIndex uint16
}

func (ConstantInvokeDynamicInfo) Type() int {
	return ConstantInvokeDynamic
}
