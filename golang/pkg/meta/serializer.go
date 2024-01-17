package meta

import (
	"bytes"
	"encoding/binary"
)

// internel inteface for req reqSerializer and resp deserializer
type reqSerializer interface {
	serialize(serializer) []byte
}

type respDeserializer interface {
	deserialize(deserializer) error
}

type serializer interface {
	serializeString(str string)
	serializeStringArray(arr []string)
	serializeBool(b bool)
	serializeUINT8(i uint8)
	serializeINT64(i int64)
	serializeUINT32(i uint32)
	serializeHeader(header *headerRequest)
	getBytes() []byte
	reset()
}

type deserializer interface {
	deserializeHeader() (*HeaderResponse, error)
	deserializeString() (string, error)
	deserializeStringArray() ([]string, error)
	deserializeINT8() (int8, error)
	deserializeBOOL() (bool, error)
	deserializeUINT32() (uint32, error)
	deserializeINT64() (int64, error)
	deserializeUINT64() (uint64, error)
	setBytes([]byte)
	reset()
}

type defaultSerializer struct {
	writer *bytes.Buffer
}

type defaultDeserializer struct {
	reader *bytes.Reader
}

func newDefaultSerializer() *defaultSerializer {
	return &defaultSerializer{
		writer: bytes.NewBuffer(nil),
	}
}

func (s *defaultSerializer) reset() {
	s.writer.Reset()
}

func (s *defaultDeserializer) reset() {
	s.reader.Reset(nil)
}

func (s *defaultSerializer) serializeString(str string) {
	length := uint32(len(str))
	binary.Write(s.writer, binary.LittleEndian, length)
	binary.Write(s.writer, binary.LittleEndian, []byte(str))
}

func (s *defaultSerializer) serializeStringArray(arr []string) {
	length := uint32(len(arr))
	binary.Write(s.writer, binary.LittleEndian, length)
	for _, str := range arr {
		s.serializeString(str)
	}
}

func (s *defaultSerializer) serializeBool(b bool) {
	if b {
		s.serializeUINT8(1)
	} else {
		s.serializeUINT8(0)
	}
}

func (s *defaultSerializer) serializeINT64(i int64) {
	binary.Write(s.writer, binary.LittleEndian, i)
}

func (s *defaultSerializer) serializeUINT32(i uint32) {
	binary.Write(s.writer, binary.LittleEndian, i)
}

func (s *defaultSerializer) serializeUINT8(i uint8) {
	binary.Write(s.writer, binary.LittleEndian, i)
}

func (s *defaultSerializer) serializeHeader(header *headerRequest) {
	s.serializeString(header.requestType)
	s.serializeINT64(header.clusterId)
}

func (s *defaultSerializer) getBytes() []byte {
	return s.writer.Bytes()
}

func newDefaultDeserializer() *defaultDeserializer {
	return &defaultDeserializer{
		reader: bytes.NewReader(nil),
	}
}

func (d *defaultDeserializer) setBytes(bs []byte) {
	d.reader.Reset(nil)
	d.reader = bytes.NewReader(bs)
}

func (d *defaultDeserializer) deserializeString() (string, error) {
	var length uint32
	if err := binary.Read(d.reader, binary.LittleEndian, &length); err != nil {
		return "", err
	}
	buf := make([]byte, length)
	if _, err := d.reader.Read(buf); err != nil {
		return "", err
	}
	return string(buf), nil
}

func (d *defaultDeserializer) deserializeStringArray() ([]string, error) {
	var length uint32
	if err := binary.Read(d.reader, binary.LittleEndian, &length); err != nil {
		return nil, err
	}
	arr := make([]string, length)
	for i := 0; i < int(length); i++ {
		str, err := d.deserializeString()
		if err != nil {
			return nil, err
		}
		arr[i] = str
	}
	return arr, nil
}

func (d *defaultDeserializer) deserializeINT8() (int8, error) {
	var i int8
	if err := binary.Read(d.reader, binary.LittleEndian, &i); err != nil {
		return 0, err
	}
	return i, nil
}

func (d *defaultDeserializer) deserializeBOOL() (bool, error) {
	var i int8
	if err := binary.Read(d.reader, binary.LittleEndian, &i); err != nil {
		return false, err
	}
	return i != 0, nil
}

func (d *defaultDeserializer) deserializeINT64() (int64, error) {
	var i int64
	if err := binary.Read(d.reader, binary.LittleEndian, &i); err != nil {
		return 0, err
	}
	return i, nil
}

func (d *defaultDeserializer) deserializeUINT64() (uint64, error) {
	var i uint64
	if err := binary.Read(d.reader, binary.LittleEndian, &i); err != nil {
		return 0, err
	}
	return i, nil
}

func (d *defaultDeserializer) deserializeUINT32() (uint32, error) {
	var i uint32
	if err := binary.Read(d.reader, binary.LittleEndian, &i); err != nil {
		return 0, err
	}
	return i, nil
}

func (d *defaultDeserializer) deserializeHeader() (*HeaderResponse, error) {
	header := &HeaderResponse{}
	var err error
	ok, err := d.deserializeBOOL()
	if err != nil {
		return nil, err
	}
	if !ok {
		header.Code, err = d.deserializeUINT64()
		if err != nil {
			return nil, err
		}
		header.Msg, err = d.deserializeString()
		if err != nil {
			return nil, err
		}
	}

	header.NewHost, err = d.deserializeString()
	if err != nil {
		return nil, err
	}
	header.NewPort, err = d.deserializeUINT32()
	if err != nil {
		return nil, err
	}

	return header, nil
}
