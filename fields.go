package javaclass

import (
	"io"

	"github.com/MJKWoolnough/byteio"
)

type FieldInfo struct {
	AccessFlags, NameIndex, DescriptorIndex uint16
	Attributes                              []AttributeInfo
}

func (c *Class) readFields(r io.Reader) ([]FieldInfo, error) {
	br := byteio.BigEndianReader{Reader: r}
	fieldsCount, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	fields := make([]FieldInfo, fieldsCount)
	for i := uint16(0); i < fieldsCount; i++ {
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
		fields[i] = FieldInfo{
			AccessFlags:     af,
			NameIndex:       ni,
			DescriptorIndex: di,
			Attributes:      attributes,
		}
	}
	return fields, nil
}
