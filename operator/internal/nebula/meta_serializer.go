/*
Copyright 2023 Vesoft Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package nebula

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

func (s *Serializer) SerializeString(value string) error {
	length := uint32(len(value))
	if err := binary.Write(&s.buffer, binary.LittleEndian, length); err != nil {
		return err
	}
	if err := binary.Write(&s.buffer, binary.LittleEndian, []byte(value)); err != nil {
		return err
	}
	return nil
}

func (s *Serializer) SerializeStringArray(array []string) error {
	length := uint32(len(array))
	if err := binary.Write(&s.buffer, binary.LittleEndian, length); err != nil {
		return err
	}
	for _, str := range array {
		if err := s.SerializeString(str); err != nil {
			return err
		}
	}
	return nil
}

func (s *Serializer) SerializeBool(b bool) error {
	if b {
		return s.SerializeUint8(1)
	}
	return s.SerializeUint8(0)
}

func (s *Serializer) SerializeInt64(i int64) error {
	return binary.Write(&s.buffer, binary.LittleEndian, i)
}

func (s *Serializer) SerializeUint32(i uint32) error {
	return binary.Write(&s.buffer, binary.LittleEndian, i)
}

func (s *Serializer) SerializeUint8(i uint8) error {
	return binary.Write(&s.buffer, binary.LittleEndian, i)
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
	array := make([]string, length)
	for i := 0; i < int(length); i++ {
		str, err := d.DeserializeString()
		if err != nil {
			return nil, err
		}
		array[i] = str
	}
	return array, nil
}

func (d *Deserializer) DeserializeInt8() (int8, error) {
	var i int8
	if err := binary.Read(d.reader, binary.LittleEndian, &i); err != nil {
		return 0, err
	}
	return i, nil
}

func (d *Deserializer) DeserializeBool() (bool, error) {
	var i int8
	if err := binary.Read(d.reader, binary.LittleEndian, &i); err != nil {
		return false, err
	}
	return i != 0, nil
}

func (d *Deserializer) DeserializeInt64() (int64, error) {
	var i int64
	if err := binary.Read(d.reader, binary.LittleEndian, &i); err != nil {
		return 0, err
	}
	return i, nil
}

func (d *Deserializer) DeserializeUint64() (uint64, error) {
	var i uint64
	if err := binary.Read(d.reader, binary.LittleEndian, &i); err != nil {
		return 0, err
	}
	return i, nil
}

func (d *Deserializer) DeserializeUint32() (uint32, error) {
	var i uint32
	if err := binary.Read(d.reader, binary.LittleEndian, &i); err != nil {
		return 0, err
	}
	return i, nil
}
