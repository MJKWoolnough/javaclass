package javaclass

import (
	"io"

	"github.com/MJKWoolnough/byteio"
)

const (
	FrameMaxSame                      = 63
	FrameMaxSameLocals                = 127
	FrameSameLocals1StackItemExtended = 247
	FrameMaxChop                      = 250
	FrameSameExtended                 = 251
	FrameMaxAppend                    = 254
	FrameFull                         = 255
)

type StackMapFrame interface {
	FrameType() uint8
}

func readStackMapFrame(r io.Reader) (StackMapFrame, error) {
	br := byteio.BigEndianReader{r}
	frameType, _, err := br.ReadUint8()
	if err != nil {
		return nil, err
	}
	var stackMapFrame StackMapFrame
	switch {
	case frameType <= FrameMaxSame:
		stackMapFrame, err = readSameFrame(r, frameType)
	case frameType <= FrameMaxSameLocals:
		stackMapFrame, err = readSameLocalsFrame(r, frameType)
	case frameType == FrameSameLocals1StackItemExtended:
		stackMapFrame, err = readSameLocals1StackItemExtendedFrame(r, frameType)
	case frameType <= FrameMaxChop:
		stackMapFrame, err = readChopFrame(r, frameType)
	case frameType == FrameSameExtended:
		stackMapFrame, err = readSameExtendedFrame(r, frameType)
	case frameType <= FrameMaxAppend:
		stackMapFrame, err = readAppendFrame(r, frameType)
	case frameType == FrameFull:
		stackMapFrame, err = readFullFrame(r, frameType)
	}
	if err != nil {
		return nil, err
	}
	return stackMapFrame, nil
}

type SameFrame struct {
	frameType uint8
}

func readSameFrame(r io.Reader, frameType uint8) (StackMapFrame, error) {
	return SameFrame{frameType}, nil
}

func (s SameFrame) FrameType() uint8 {
	return s.frameType
}

type SameLocals1StackItemFrame struct {
	frameType uint8
	Stack     VerificationTypeInfo
}

func readSameLocalsFrame(r io.Reader, frameType uint8) (StackMapFrame, error) {
	stack, err := readVerificationTypeInfo(r)
	if err != nil {
		return nil, err
	}
	return SameLocals1StackItemFrame{frameType, stack}, nil
}

func (s SameLocals1StackItemFrame) FrameType() uint8 {
	return s.frameType
}

type SameLocals1StackItemFrameExtended struct {
	OffsetDelta uint16
	Stack       VerificationTypeInfo
}

func readSameLocals1StackItemExtendedFrame(r io.Reader, _ uint8) (StackMapFrame, error) {
	br := byteio.BigEndianReader{r}
	offsetDelta, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	stack, err := readVerificationTypeInfo(r)
	if err != nil {
		return nil, err
	}
	return SameLocals1StackItemFrameExtended{
		OffsetDelta: offsetDelta,
		Stack:       stack,
	}, nil
}

func (SameLocals1StackItemFrameExtended) FrameType() uint8 {
	return FrameSameLocals1StackItemExtended
}

type ChopFrame struct {
	frameType   uint8
	OffsetDelta uint16
}

func readChopFrame(r io.Reader, frameType uint8) (StackMapFrame, error) {
	br := byteio.BigEndianReader{r}
	offsetDelta, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	return ChopFrame{frameType, offsetDelta}, nil
}

func (c ChopFrame) FrameType() uint8 {
	return c.frameType
}

type SameFrameExtended struct {
	OffsetDelta uint16
}

func readSameExtendedFrame(r io.Reader, _ uint8) (StackMapFrame, error) {
	br := byteio.BigEndianReader{r}
	offsetDelta, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	return SameFrameExtended{offsetDelta}, nil
}

func (SameFrameExtended) FrameType() uint8 {
	return FrameSameExtended
}

type AppendFrame struct {
	frameType   uint8
	OffsetDelta uint16
	Locals      []VerificationTypeInfo
}

func readAppendFrame(r io.Reader, frameType uint8) (StackMapFrame, error) {
	br := byteio.BigEndianReader{r}
	offsetDelta, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	locals := make([]VerificationTypeInfo, frameType-251)
	for i := uint8(0); i < frameType-251; i++ {
		locals[i], err = readVerificationTypeInfo(r)
		if err != nil {
			return nil, err
		}
	}
	return AppendFrame{frameType, offsetDelta, locals}, nil
}

func (a AppendFrame) FrameType() uint8 {
	return a.frameType
}

type FullFrame struct {
	OffsetDelta   uint16
	Locals, Stack []VerificationTypeInfo
}

func readFullFrame(r io.Reader, _ uint8) (StackMapFrame, error) {
	br := byteio.BigEndianReader{r}
	offsetDelta, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	numberOfLocals, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	locals := make([]VerificationTypeInfo, numberOfLocals)
	for i := uint16(0); i < numberOfLocals; i++ {
		locals[i], err = readVerificationTypeInfo(r)
		if err != nil {
			return nil, err
		}
	}
	numberOfStackItems, _, err := br.ReadUint16()
	if err != nil {
		return nil, err
	}
	stack := make([]VerificationTypeInfo, numberOfStackItems)
	for i := uint16(0); i < numberOfStackItems; i++ {
		stack[i], err = readVerificationTypeInfo(r)
		if err != nil {
			return nil, err
		}
	}
	return FullFrame{
		OffsetDelta: offsetDelta,
		Locals:      locals,
		Stack:       stack,
	}, nil
}

func (FullFrame) FrameType() uint8 {
	return FrameFull
}
