package javaclass // import "vimagination.zapto.org/javaclass"

import (
	"errors"
	"io"

	"vimagination.zapto.org/byteio"
)

const Magic = 0xCAFEBABE

var (
	ErrInvalidMagic = errors.New("read invalid magic string")
)

type Class struct {
	Minor, Major                       uint16
	ConstantPool                       []CPInfo
	AccessFlags, ThisClass, SuperClass uint16
	Interfaces                         []uint16
	Fields                             []FieldInfo
	Methods                            []MethodInfo
	Attributes                         []AttributeInfo
}

func Read(r io.Reader) (*Class, error) {
	br := byteio.BigEndianReader{Reader: r}
	magic, _, err := br.ReadUint32()
	if err != nil {
		return nil, err
	}
	if magic != Magic {
		return nil, ErrInvalidMagic
	}
	minor, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	major, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	cp, err := readConstantPool(r)
	if err != nil {
		return nil, err
	}
	accessFlags, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	thisClass, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	superClass, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}

	interfacesCount, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	interfaces := make([]uint16, interfacesCount)
	for i := uint16(0); i < interfacesCount; i++ {
		interfaces[i], _, err = br.ReadUint16()
		if err != nil {
			return nil, err
		}
	}

	c := &Class{
		Minor:        minor,
		Major:        major,
		ConstantPool: cp,
		AccessFlags:  accessFlags,
		ThisClass:    thisClass,
		SuperClass:   superClass,
		Interfaces:   interfaces,
	}

	c.Fields, err = c.readFields(r)
	if err != nil {
		return nil, err
	}

	c.Methods, err = c.readMethods(r)
	if err != nil {
		return nil, err
	}

	c.Attributes, err = c.readAttributes(r)
	if err != nil {
		return nil, err
	}

	return c, err
}
