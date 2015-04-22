package javaclass

const (
	FrameSameLocals1StackItemExtended = 247
	FrameSameExtended                 = 251
)

type StackMapFrame interface {
	FrameType() uint8
}

type SameFrame struct {
	frameType uint8
}

func (s SameFrame) FrameType() uint8 {
	return s.frameType
}

type SameLocals1StackItemFrame struct {
	frameType uint8
	Stack     VerificationTypeInfo
}

func (s SameLocals1StackItemFrame) FrameType() uint8 {
	return s.frameType
}

type SameLocals1StackItemFrameExtended struct {
	OffsetDelta uint16
	Stack       VerificationTypeInfo
}

func (SameLocals1StackItemFrameExtended) FrameType() uint8 {
	return FrameSameLocals1StackItemExtended
}

type ChopFrame struct {
	frameType   uint8
	OffsetDelta uint16
}

func (c ChopFrame) FrameType() uint8 {
	return c.frameType
}

type SameFrameExtended struct {
	OffsetDelta uint16
}

func (SameFrameExtended) FrameType() uint8 {
	return FrameSameExtended
}

type AppendFrame struct {
	frameType   uint8
	OffsetDelta uint16
	Locals      []VerificationTypeInfo
}

func (a AppendFrame) FrameType() uint8 {
	return a.frameType
}
