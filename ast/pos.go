package ast

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
)

// Position provides interface to store code locations.
type Position struct {
	Line   int
	Column int
}

// Pos interface provies two functions to get/set the position for expression or statement.
type Pos interface {
	Position() Position
	SetPosition(Position)
	gob.GobEncoder
	gob.GobDecoder
}

// PosImpl provies commonly implementations for Pos.
type PosImpl struct {
	pos Position
}

// Position return the position of the expression or statement.
func (x *PosImpl) Position() Position {
	return x.pos
}

// SetPosition is a function to specify position of the expression or statement.
func (x *PosImpl) SetPosition(pos Position) {
	x.pos = pos
}

// Implementation of GobEncoder interface
func (x *PosImpl) GobEncode() ([]byte, error) {
	var err error
	var buf bytes.Buffer
	if err = binary.Write(&buf, binary.LittleEndian, (int32)(x.pos.Line)); err != nil {
		return nil, err
	}
	if err = binary.Write(&buf, binary.LittleEndian, (int32)(x.pos.Column)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

//
func (x *PosImpl) GobDecode(data []byte) error {
	var err error
	var tmpInt int32
	buf := bytes.NewBuffer(data)
	if err = binary.Read(buf, binary.LittleEndian, &tmpInt); err != nil {
		return err
	}
	x.pos.Line = (int)(tmpInt)
	if err = binary.Read(buf, binary.LittleEndian, &tmpInt); err != nil {
		return err
	}
	x.pos.Column = (int)(tmpInt)
	return nil
}