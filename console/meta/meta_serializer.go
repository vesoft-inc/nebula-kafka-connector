package meta

import (
	"bytes"
	"encoding/binary"
)

type Serializer struct {
	buffer bytes.Buffer
}

func NewSerializer() *Serializer {
	return &Serializer{}
}

func (s *Serializer) SerializeString(str string) {
	length := uint32(len(str))
	binary.Write(&s.buffer, binary.LittleEndian, length)
	binary.Write(&s.buffer, binary.LittleEndian, []byte(str))
}

func (s *Serializer) SerializeStringArray(arr []string) {
	length := uint32(len(arr))
	binary.Write(&s.buffer, binary.LittleEndian, length)
	for _, str := range arr {
		s.SerializeString(str)
	}
}

func (s *Serializer) SerializeBool(b bool) {
	if b {
		s.SerializeUINT8(1)
	} else {
		s.SerializeUINT8(0)
	}
}

func (s *Serializer) SerializeINT64(i int64) {
	binary.Write(&s.buffer, binary.LittleEndian, i)
}

func (s *Serializer) SerializeUINT32(i uint32) {
	binary.Write(&s.buffer, binary.LittleEndian, i)
}

func (s *Serializer) SerializeUINT8(i uint8) {
	binary.Write(&s.buffer, binary.LittleEndian, i)
}

func (s *Serializer) GetBytes() []byte {
	return s.buffer.Bytes()
}

type Deserializer struct {
	reader *bytes.Reader
}

func NewDeserializer(data []byte) *Deserializer {
	return &Deserializer{reader: bytes.NewReader(data)}
}

func (d *Deserializer) DeserializeString() (string, error) {
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

func (d *Deserializer) DeserializeStringArray() ([]string, error) {
	var length uint32
	if err := binary.Read(d.reader, binary.LittleEndian, &length); err != nil {
		return nil, err
	}
	arr := make([]string, length)
	for i := 0; i < int(length); i++ {
		str, err := d.DeserializeString()
		if err != nil {
			return nil, err
		}
		arr[i] = str
	}
	return arr, nil
}

func (d *Deserializer) DeserializeINT8() (int8, error) {
	var i int8
	if err := binary.Read(d.reader, binary.LittleEndian, &i); err != nil {
		return 0, err
	}
	return i, nil
}

func (d *Deserializer) DeserializeBOOL() (bool, error) {
	var i int8
	if err := binary.Read(d.reader, binary.LittleEndian, &i); err != nil {
		return false, err
	}
	return i != 0, nil
}

func (d *Deserializer) DeserializeINT64() (int64, error) {
	var i int64
	if err := binary.Read(d.reader, binary.LittleEndian, &i); err != nil {
		return 0, err
	}
	return i, nil
}

func (d *Deserializer) DeserializeUINT64() (uint64, error) {
	var i uint64
	if err := binary.Read(d.reader, binary.LittleEndian, &i); err != nil {
		return 0, err
	}
	return i, nil
}

func (d *Deserializer) DeserializeUINT32() (uint32, error) {
	var i uint32
	if err := binary.Read(d.reader, binary.LittleEndian, &i); err != nil {
		return 0, err
	}
	return i, nil
}
