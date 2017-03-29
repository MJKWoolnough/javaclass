package javaclass

import (
	"io"

	"github.com/MJKWoolnough/byteio"
)

type MethodInfo struct {
	AccessFlags, NameIndex, DescriptorIndex uint16
	Attributes                              []AttributeInfo
}

func (c *Class) readMethods(r io.Reader) ([]MethodInfo, error) {
	br := byteio.BigEndianReader{Reader: r}
	methodsCount, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	methods := make([]MethodInfo, methodsCount)
	for i := uint16(0); i < methodsCount; i++ {
		af, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		ni, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		di, _, err := br.ReadUint16()
		if err != nil {
			return nil, err
		}
		attributes, err := c.readAttributes(r)
		if err != nil {
			return nil, err
		}
		methods[i] = MethodInfo{
			AccessFlags:     af,
			NameIndex:       ni,
			DescriptorIndex: di,
			Attributes:      attributes,
		}
	}
	return methods, nil
}
