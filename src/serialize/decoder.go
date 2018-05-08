package serialize

import "fmt"

type HessianDecoder struct {
	BytecodeReader
}

/**
 * read as
 * {
 *  type: com.xx.xx
 *  value: {}
 * }
 */
func (this *HessianDecoder) readObject() (interface{}, error) {
	tag := this.readUInt8()
	switch tag {
	case 'N':
		return nil, nil
	case 'T':
		return true, nil
	case 'F':
		return false, nil

	// one byte compact integer
	case 0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
		0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
		0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xaa, 0xab, 0xac, 0xad, 0xae, 0xaf,
		0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xba, 0xbb, 0xbc, 0xbd, 0xbe, 0xbf:
		return int8(tag - BC_INT_ZERO), nil

	// two byte compact int
	case 0xc0, 0xc1, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0xca, 0xcb, 0xcc, 0xcd, 0xce, 0xcf:
		lower := int16(this.readUInt8())
		return (int16(tag-BC_INT_BYTE_ZERO) << 8) | lower, nil

	// three byte compact int
	case 0xd0, 0xd1, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7:
		middle := int32(this.readUInt8())
		lower := int32(this.readUInt8())
		return (int32(tag-BC_INT_SHORT_ZERO) << 16) | (middle << 8) | lower, nil

	// int32
	case 'I':
		return parseInt32(this)

	// one byte direct long
	case 0xd8, 0xd9, 0xda, 0xdb, 0xdc, 0xdd, 0xde, 0xdf,
		0xe0, 0xe1, 0xe2, 0xe3, 0xe4, 0xe5, 0xe6, 0xe7,
		0xe8, 0xe9, 0xea, 0xeb, 0xec, 0xed, 0xee, 0xef:
		return int64(tag - BC_LONG_ZERO), nil

	// two byte compact long
	case 0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7, 0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff:
		lower := int64(this.readUInt8())
		return int64(tag-BC_LONG_BYTE_ZERO)<<8 | lower, nil

	case 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f:
		b16 := int64(this.readUInt8())
		b8 := int64(this.readUInt8())
		return int64(tag-BC_LONG_SHORT_ZERO)<<16 | b16<<8 | b8, nil

	// long as 32-bit int
	case BC_LONG_INT:
		ret, err := parseInt32(this)
		if err != nil {
			return nil, err
		}
		return int64(ret), nil

	case 'L':
		// TODO

	}

	result := make(map[string]interface{})
	return result, nil
}

func parseInt32(this *HessianDecoder) (int32, error) {
	b32 := int32(this.readUInt8())
	b24 := int32(this.readUInt8())
	b16 := int32(this.readUInt8())
	b8 := int32(this.readUInt8())
	return (b32 << 24) | (b24 << 16) | (b16 << 8) | b8, nil
}

// start message first three bytes must be 700200 in hessian 2 protocol
func (this *HessianDecoder) startMessage() uint16 {
	tag := this.readUInt8()
	major := uint16(this.readUInt8())
	minor := uint16(this.readUInt8())

	// check tag
	if !(tag == 'p' || tag == 'P') {
		panic("expected Hessian message ('p') at " + fmt.Sprintf("%x", tag))
	}

	// higher 8 bit is major version, lower 8 bit is minor version
	// as we use hession 2 protocol, the major version is 0x02
	return (major << 8) | minor
}

/**
 * Completes reading the message
 *
 * <p>A successful completion will have a single value:
 *
 * <pre>
 * z
 * </pre>
 */
func (this *HessianDecoder) completeMessage() {
	tag := this.readUInt8()

	if !(tag == 'Z' || tag == 'z') {
		panic("expected end of message at" + fmt.Sprintf("%x", tag))
	}
}
