package javaclass

type CPInfo interface {
	Type() int
}

type FieldInfo struct {
	AccessFlags, NameIndex, DescriptorIndex uint16
	Attributes                              []AttributeInfo
}

type MethodInfo struct {
	AccessFlags, NameIndex, DescriptorIndex uint16
	Attributes                              []AttributeInfo
}

type Class struct {
	Magic                              uint32
	Minor, Major                       uint16
	ConstantPool                       []CPInfo
	AccessFlags, ThisClass, SuperClass uint16
	Interfaces                         []uint16
	Fields                             []FieldInfo
	Methods                            []MethodInfo
	Attributes                         []AttributeInfo
}
