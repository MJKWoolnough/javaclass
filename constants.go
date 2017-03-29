package javaclass

import (
	"io"

	"github.com/MJKWoolnough/byteio"
)

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

type CPInfo interface {
	Type() int
}

func readConstantPool(r io.Reader) ([]CPInfo, error) {
	br := byteio.BigEndianReader{Reader: r}
	constantPoolCount, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	constantPool := make([]CPInfo, 1, constantPoolCount)
	constantPool[0] = ConstantNullInfo{}
	for i := uint16(1); i < constantPoolCount; i++ {
		tag, _, err := br.ReadUint8()
		if err != nil {
			return nil, err
		}
		var (
			cpInfo  CPInfo
			addNull bool
		)
		switch tag {
		case ConstantUTF8:
			cpInfo, err = readConstantUTF8(r)
		case ConstantInteger:
			cpInfo, err = readConstantInteger(r)
		case ConstantFloat:
			cpInfo, err = readConstantFloat(r)
		case ConstantLong:
			cpInfo, err = readConstantLong(r)
			addNull = true
		case ConstantDouble:
			cpInfo, err = readConstantDouble(r)
			addNull = true
		case ConstantClass:
			cpInfo, err = readConstantClass(r)
		case ConstantString:
			cpInfo, err = readConstantString(r)
		case ConstantFieldRef:
			cpInfo, err = readConstantFieldRef(r)
		case ConstantMethodRef:
			cpInfo, err = readConstantMethodRef(r)
		case ConstantInterfaceMethodRef:
			cpInfo, err = readConstantInterfaceMethodRef(r)
		case ConstantNameAndType:
			cpInfo, err = readConstantNameAndType(r)
		case ConstantMethodHandle:
			cpInfo, err = readConstantMethodHandle(r)
		case ConstantMethodType:
			cpInfo, err = readConstantMethodType(r)
		case ConstantInvokeDynamic:
			cpInfo, err = readConstantInvokeDynamic(r)
		default:
			err = ErrUnknownConstantPoolTag{tag}
		}
		if err != nil {
			return nil, err
		}
		constantPool = append(constantPool, cpInfo)
		if addNull {
			constantPool = append(constantPool, ConstantNullInfo{})
			i++
		}
	}
	return constantPool, nil
}

type ConstantNullInfo struct{}

func (ConstantNullInfo) Type() int {
	return ConstantNull
}

type ConstantUTF8Info struct {
	String string
}

func readConstantUTF8(r io.Reader) (CPInfo, error) {
	br := byteio.BigEndianReader{Reader: r}
	length, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	toString := make([]byte, length)
	_, err = br.Read(toString)
	if err != nil {
		return nil, err
	}
	return ConstantUTF8Info{string(toString)}, nil
}

func (ConstantUTF8Info) Type() int {
	return ConstantUTF8
}

type ConstantIntegerInfo struct {
	Integer uint32
}

func readConstantInteger(r io.Reader) (CPInfo, error) {
	br := byteio.BigEndianReader{Reader: r}
	n, _, err := br.ReadUint32()
	if err != nil {
		return nil, err
	}
	return ConstantIntegerInfo{n}, nil
}

func (ConstantIntegerInfo) Type() int {
	return ConstantInteger
}

type ConstantFloatInfo struct {
	Float float32
}

func readConstantFloat(r io.Reader) (CPInfo, error) {
	br := byteio.BigEndianReader{Reader: r}
	n, _, err := br.ReadFloat32()
	if err != nil {
		return nil, err
	}
	return ConstantFloatInfo{n}, nil
}

func (ConstantFloatInfo) Type() int {
	return ConstantFloat
}

type ConstantLongInfo struct {
	Long uint64
}

func readConstantLong(r io.Reader) (CPInfo, error) {
	br := byteio.BigEndianReader{Reader: r}
	n, _, err := br.ReadUint64()
	if err != nil {
		return nil, err
	}
	return ConstantLongInfo{n}, nil
}

func (ConstantLongInfo) Type() int {
	return ConstantLong
}

type ConstantDoubleInfo struct {
	Double float64
}

func readConstantDouble(r io.Reader) (CPInfo, error) {
	br := byteio.BigEndianReader{Reader: r}
	n, _, err := br.ReadFloat64()
	if err != nil {
		return nil, err
	}
	return ConstantDoubleInfo{n}, nil
}

func (ConstantDoubleInfo) Type() int {
	return ConstantDouble
}

type ConstantClassInfo struct {
	NameIndex uint16
}

func readConstantClass(r io.Reader) (CPInfo, error) {
	br := byteio.BigEndianReader{Reader: r}
	n, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	return ConstantClassInfo{n}, nil
}

func (ConstantClassInfo) Type() int {
	return ConstantClass
}

type ConstantStringInfo struct {
	StringIndex uint16
}

func readConstantString(r io.Reader) (CPInfo, error) {
	br := byteio.BigEndianReader{Reader: r}
	n, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	return ConstantStringInfo{n}, nil
}

func (ConstantStringInfo) Type() int {
	return ConstantString
}

type ConstantFieldRefInfo struct {
	ClassIndex, NameAndTypeIndex uint16
}

func readConstantFieldRef(r io.Reader) (CPInfo, error) {
	br := byteio.BigEndianReader{Reader: r}
	i, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	n, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	return ConstantFieldRefInfo{i, n}, nil
}

func (ConstantFieldRefInfo) Type() int {
	return ConstantFieldRef
}

type ConstantMethodRefInfo struct {
	ClassIndex, NameAndTypeIndex uint16
}

func readConstantMethodRef(r io.Reader) (CPInfo, error) {
	br := byteio.BigEndianReader{Reader: r}
	i, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	n, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	return ConstantMethodRefInfo{i, n}, nil
}

func (ConstantMethodRefInfo) Type() int {
	return ConstantMethodRef
}

type ConstantInterfaceMethodRefInfo struct {
	ClassIndex, NameAndTypeIndex uint16
}

func readConstantInterfaceMethodRef(r io.Reader) (CPInfo, error) {
	br := byteio.BigEndianReader{Reader: r}
	i, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	n, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	return ConstantInterfaceMethodRefInfo{i, n}, nil
}

func (ConstantInterfaceMethodRefInfo) Type() int {
	return ConstantInterfaceMethodRef
}

type ConstantNameAndTypeInfo struct {
	NameIndex, DescriptorIndex uint16
}

func readConstantNameAndType(r io.Reader) (CPInfo, error) {
	br := byteio.BigEndianReader{Reader: r}
	n, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	d, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	return ConstantNameAndTypeInfo{n, d}, nil
}

func (ConstantNameAndTypeInfo) Type() int {
	return ConstantNameAndType
}

type ConstantMethodHandleInfo struct {
	ReferenceKind  uint8
	ReferenceIndex uint16
}

func readConstantMethodHandle(r io.Reader) (CPInfo, error) {
	br := byteio.BigEndianReader{Reader: r}
	k, _, err := br.ReadUint8()
	if err != nil {
		return nil, err
	}
	i, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	return ConstantMethodHandleInfo{k, i}, nil
}

func (ConstantMethodHandleInfo) Type() int {
	return ConstantMethodHandle
}

type ConstantMethodTypeInfo struct {
	DescriptorIndex uint16
}

func readConstantMethodType(r io.Reader) (CPInfo, error) {
	br := byteio.BigEndianReader{Reader: r}
	i, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	return ConstantMethodTypeInfo{i}, nil
}

func (ConstantMethodTypeInfo) Type() int {
	return ConstantMethodType
}

type ConstantInvokeDynamicInfo struct {
	BootstrapMethodAttrIndex, NameAndTypeIndex uint16
}

func readConstantInvokeDynamic(r io.Reader) (CPInfo, error) {
	br := byteio.BigEndianReader{Reader: r}
	b, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	i, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	return ConstantInvokeDynamicInfo{b, i}, nil
}

func (ConstantInvokeDynamicInfo) Type() int {
	return ConstantInvokeDynamic
}

// Error types

type ErrUnknownConstantPoolTag struct {
	Tag uint8
}

func (ErrUnknownConstantPoolTag) Error() string {
	return "unknown constant pool tag"
}
