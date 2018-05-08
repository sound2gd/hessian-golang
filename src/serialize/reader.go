package serialize

// read raw bytes and convert to golang's objects
type BytecodeReader struct {
	rawData []byte // original byte data
	pc      uint   // current PC
}

func (this *BytecodeReader) readUInt8() uint8 {
	ret := this.rawData[this.pc]
	this.pc++
	return ret
}
