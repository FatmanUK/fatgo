package serialisers

import (
	"math"
	"errors"
	"math/rand"  // for array wipe
)

func isLittleEndian() bool {
	var testWord uint16 = 1
	testByte := byte(testWord >> 8)
	return (testByte == 0)
}

func wipeArrayRandom(a *[]byte) {
	for k, _ := range *a {
		(*a)[k] = byte(rand.Int() % 256)
	}
}

type Hashmap = map[string]string
type Stringlist = []string

const ERR_ARRAY_TOO_SMALL = "Array too small."

/////////////////////////////////////////////////////////////////////

type Serialiser interface {
	IoInt32(*int32) error
	IoInt64(*int64) error
	IoUint32(*uint32) error
	IoUint64(*uint64) error
	IoFloat32(*float32) error
	IoFloat64(*float64) error
	IoString(*string) error
	IoHashmap(*Hashmap) error
	IoStringlist(*Stringlist) error
}

type Serialisable interface {
	serialise(Serialiser)
}

/////////////////////////////////////////////////////////////////////

type Sizer struct {
	Size *uint32
}

func (re Sizer) Init(s *uint32) Sizer {
	re.Size = s
	return re
}

func (re *Sizer) IoInt32(i *int32) error {
	*(re.Size) += 4
	return nil
}

func (re *Sizer) IoInt64(i *int64) error {
	*(re.Size) += 8
	return nil
}

func (re *Sizer) IoUint32(i *uint32) error {
	*(re.Size) += 4
	return nil
}

func (re *Sizer) IoUint64(i *uint64) error {
	*(re.Size) += 8
	return nil
}

func (re *Sizer) IoFloat32(f *float32) error {
	*(re.Size) += 4
	return nil
}

func (re *Sizer) IoFloat64(f *float64) error {
	*(re.Size) += 8
	return nil
}

func (re *Sizer) IoString(str *string) error {
	l := uint32(len(*str))
	re.IoUint32(&l)
	*(re.Size) += uint32(l)
	return nil
}

func (re *Sizer) IoHashmap(hm *Hashmap) error {
	l := uint32(len(*hm))
	re.IoUint32(&l)
	for k, v := range *hm {
		re.IoString(&k)
		re.IoString(&v)
	}
	return nil
}

func (re *Sizer) IoStringlist(sl *Stringlist) error {
	l := uint32(len(*sl))
	re.IoUint32(&l)
	for _, v := range *sl {
		re.IoString(&v)
	}
	return nil
}

/////////////////////////////////////////////////////////////////////

type Saver struct {
	Array []byte
	index uint64
}

func (re Saver) Init(s *uint32) Saver {
	re.Array = make([]byte, *s)
	return re
}

func (re *Saver) SecureWipe() {
	wipeArrayRandom(&re.Array)
}

// have to use pointer receiver for these so we can use Serialiser interface
func (re *Saver) IoInt32(i *int32) error {
	var err error = nil
	if uint64(len(re.Array)) < re.index + 4 {
		err = errors.New(ERR_ARRAY_TOO_SMALL)
		return err
	}
	a := byte(*i >> 24)
	b := byte(*i >> 16)
	c := byte(*i >> 8)
	d := byte(*i >> 0)
	mul := 1
	add := 0
	if isLittleEndian() {
		mul = -1
		add = 3
	}
	re.Array[re.index + uint64(add + (0 * mul))] = byte(a)
	re.Array[re.index + uint64(add + (1 * mul))] = byte(b)
	re.Array[re.index + uint64(add + (2 * mul))] = byte(c)
	re.Array[re.index + uint64(add + (3 * mul))] = byte(d)
	re.index += 4
	return err
}

func (re *Saver) IoInt64(i *int64) error {
	var err error = nil
	if uint64(len(re.Array)) < re.index + 4 {
		err = errors.New(ERR_ARRAY_TOO_SMALL)
		return err
	}
	a := byte(*i >> 56)
	b := byte(*i >> 48)
	c := byte(*i >> 40)
	d := byte(*i >> 32)
	e := byte(*i >> 24)
	f := byte(*i >> 16)
	g := byte(*i >> 8)
	h := byte(*i >> 0)
	mul := 1
	add := 0
	if isLittleEndian() {
		mul = -1
		add = 7
	}
	re.Array[re.index + uint64(add + (0 * mul))] = byte(a)
	re.Array[re.index + uint64(add + (1 * mul))] = byte(b)
	re.Array[re.index + uint64(add + (2 * mul))] = byte(c)
	re.Array[re.index + uint64(add + (3 * mul))] = byte(d)
	re.Array[re.index + uint64(add + (4 * mul))] = byte(e)
	re.Array[re.index + uint64(add + (5 * mul))] = byte(f)
	re.Array[re.index + uint64(add + (6 * mul))] = byte(g)
	re.Array[re.index + uint64(add + (7 * mul))] = byte(h)
	re.index += 8
	return err
}

func (re *Saver) IoUint32(i *uint32) error {
	var err error = nil
	if uint64(len(re.Array)) < re.index + 4 {
		err = errors.New(ERR_ARRAY_TOO_SMALL)
		return err
	}
	a := byte(*i >> 24)
	b := byte(*i >> 16)
	c := byte(*i >> 8)
	d := byte(*i >> 0)
	mul := 1
	add := 0
	if isLittleEndian() {
		mul = -1
		add = 3
	}
	re.Array[re.index + uint64(add + (0 * mul))] = byte(a)
	re.Array[re.index + uint64(add + (1 * mul))] = byte(b)
	re.Array[re.index + uint64(add + (2 * mul))] = byte(c)
	re.Array[re.index + uint64(add + (3 * mul))] = byte(d)
	re.index += 4
	return err
}

func (re *Saver) IoUint64(i *uint64) error {
	var err error = nil
	if uint64(len(re.Array)) < re.index + 4 {
		err = errors.New(ERR_ARRAY_TOO_SMALL)
		return err
	}
	a := byte(*i >> 56)
	b := byte(*i >> 48)
	c := byte(*i >> 40)
	d := byte(*i >> 32)
	e := byte(*i >> 24)
	f := byte(*i >> 16)
	g := byte(*i >> 8)
	h := byte(*i >> 0)
	mul := 1
	add := 0
	if isLittleEndian() {
		mul = -1
		add = 7
	}
	re.Array[re.index + uint64(add + (0 * mul))] = byte(a)
	re.Array[re.index + uint64(add + (1 * mul))] = byte(b)
	re.Array[re.index + uint64(add + (2 * mul))] = byte(c)
	re.Array[re.index + uint64(add + (3 * mul))] = byte(d)
	re.Array[re.index + uint64(add + (4 * mul))] = byte(e)
	re.Array[re.index + uint64(add + (5 * mul))] = byte(f)
	re.Array[re.index + uint64(add + (6 * mul))] = byte(g)
	re.Array[re.index + uint64(add + (7 * mul))] = byte(h)
	re.index += 8
	return err
}

func (re *Saver) IoFloat32(f *float32) error {
	var err error = nil
	if uint64(len(re.Array)) < re.index + 4 {
		err = errors.New(ERR_ARRAY_TOO_SMALL)
		return err
	}
	bits := math.Float32bits(*f)
	var u uint32 = uint32(bits)
	re.IoUint32(&u)
	return err
}

func (re *Saver) IoFloat64(f *float64) error {
	var err error = nil
	if uint64(len(re.Array)) < re.index + 4 {
		err = errors.New(ERR_ARRAY_TOO_SMALL)
		return err
	}
	bits := math.Float64bits(*f)
	var u uint64 = uint64(bits)
	re.IoUint64(&u)
	return err
}

func (re *Saver) IoString(s *string) error {
	var err error = nil
	if uint64(len(re.Array)) < re.index + 4 {
		err = errors.New(ERR_ARRAY_TOO_SMALL)
		return err
	}
	l := uint32(len(*s))
	re.IoUint32(&l)
	for c := uint64(0); c < uint64(l); c++ {
		re.Array[re.index + c] = (*s)[c]
	}
	re.index += uint64(l)
	return err
}

func (re *Saver) IoHashmap(hm *Hashmap) error {
	var err error = nil
	if uint64(len(re.Array)) < re.index + 4 {
		err = errors.New(ERR_ARRAY_TOO_SMALL)
		return err
	}
	l := uint32(len(*hm))
	re.IoUint32(&l)
	for k, v := range *hm {
		re.IoString(&k)
		re.IoString(&v)
	}
	return err
}

func (re *Saver) IoStringlist(sl *Stringlist) error {
	var err error = nil
	if uint64(len(re.Array)) < re.index + 4 {
		err = errors.New(ERR_ARRAY_TOO_SMALL)
		return err
	}
	l := uint32(len(*sl))
	re.IoUint32(&l)
	for _, v := range *sl {
		re.IoString(&v)
	}
	return err
}

/////////////////////////////////////////////////////////////////////

type Loader struct {
	Array []byte
	index uint64
}

func (re Loader) Init(s *uint32) Loader {
	re.Array = make([]byte, *s)
	return re
}

func (re *Loader) SecureWipe() {
	wipeArrayRandom(&re.Array)
}

func (re *Loader) IoInt32(i *int32) error {
	var err error = nil
	if uint64(len(re.Array)) < re.index + 4 {
		err = errors.New(ERR_ARRAY_TOO_SMALL)
		return err
	}
	mul := 1
	add := 0
	if isLittleEndian() {
		mul = -1
		add = 3
	}
	var a byte = re.Array[re.index + uint64(add + (0 * mul))]
	var b byte = re.Array[re.index + uint64(add + (1 * mul))]
	var c byte = re.Array[re.index + uint64(add + (2 * mul))]
	var d byte = re.Array[re.index + uint64(add + (3 * mul))]
	// casts required because byte is not a normal integer
	*i = (int32(a) << 24) + (int32(b) << 16) + (int32(c) << 8)
	*i += int32(d)
	re.index += 4
	return err
}

func (re *Loader) IoInt64(i *int64) error {
	var err error = nil
	if uint64(len(re.Array)) < re.index + 4 {
		err = errors.New(ERR_ARRAY_TOO_SMALL)
		return err
	}
	mul := 1
	add := 0
	if isLittleEndian() {
		mul = -1
		add = 7
	}
	var a byte = re.Array[re.index + uint64(add + (0 * mul))]
	var b byte = re.Array[re.index + uint64(add + (1 * mul))]
	var c byte = re.Array[re.index + uint64(add + (2 * mul))]
	var d byte = re.Array[re.index + uint64(add + (3 * mul))]
	var e byte = re.Array[re.index + uint64(add + (4 * mul))]
	var f byte = re.Array[re.index + uint64(add + (5 * mul))]
	var g byte = re.Array[re.index + uint64(add + (6 * mul))]
	var h byte = re.Array[re.index + uint64(add + (7 * mul))]
	// casts required because byte is not a normal integer
	*i = (int64(a) << 56) + (int64(b) << 48) + (int64(c) << 40)
	*i += (int64(d) << 32) + (int64(e) << 24) + (int64(f) << 16)
	*i += (int64(g) << 8) + int64(h)
	re.index += 8
	return err
}

func (re *Loader) IoUint32(i *uint32) error {
	var err error = nil
	if uint64(len(re.Array)) < re.index + 4 {
		err = errors.New(ERR_ARRAY_TOO_SMALL)
		return err
	}
	mul := 1
	add := 0
	if isLittleEndian() {
		mul = -1
		add = 3
	}
	var a byte = re.Array[re.index + uint64(add + (0 * mul))]
	var b byte = re.Array[re.index + uint64(add + (1 * mul))]
	var c byte = re.Array[re.index + uint64(add + (2 * mul))]
	var d byte = re.Array[re.index + uint64(add + (3 * mul))]
	// casts required because byte is not a normal integer
	*i = (uint32(a) << 24) + (uint32(b) << 16) + (uint32(c) << 8)
	*i += uint32(d)
	re.index += 4
	return err
}

func (re *Loader) IoUint64(i *uint64) error {
	var err error = nil
	if uint64(len(re.Array)) < re.index + 4 {
		err = errors.New(ERR_ARRAY_TOO_SMALL)
		return err
	}
	mul := 1
	add := 0
	if isLittleEndian() {
		mul = -1
		add = 7
	}
	var a byte = re.Array[re.index + uint64(add + (0 * mul))]
	var b byte = re.Array[re.index + uint64(add + (1 * mul))]
	var c byte = re.Array[re.index + uint64(add + (2 * mul))]
	var d byte = re.Array[re.index + uint64(add + (3 * mul))]
	var e byte = re.Array[re.index + uint64(add + (4 * mul))]
	var f byte = re.Array[re.index + uint64(add + (5 * mul))]
	var g byte = re.Array[re.index + uint64(add + (6 * mul))]
	var h byte = re.Array[re.index + uint64(add + (7 * mul))]
	// casts required because byte is not a normal integer
	*i = (uint64(a) << 56) + (uint64(b) << 48) + (uint64(c) << 40)
	*i += (uint64(d) << 32) + (uint64(e) << 24) + (uint64(f) << 16)
	*i += (uint64(g) << 8) + uint64(h)
	re.index += 8
	return err
}

func (re *Loader) IoFloat32(f *float32) error {
	var err error = nil
	if uint64(len(re.Array)) < re.index + 4 {
		err = errors.New(ERR_ARRAY_TOO_SMALL)
		return err
	}
	var u uint32
	re.IoUint32(&u)
	*f = math.Float32frombits(u)
	return err
}

func (re *Loader) IoFloat64(f *float64) error {
	var err error = nil
	if uint64(len(re.Array)) < re.index + 4 {
		err = errors.New(ERR_ARRAY_TOO_SMALL)
		return err
	}
	var u uint64
	re.IoUint64(&u)
	*f = math.Float64frombits(u)
	return err
}

func (re *Loader) IoString(s *string) error {
	var err error = nil
	if uint64(len(re.Array)) < re.index + 4 {
		err = errors.New(ERR_ARRAY_TOO_SMALL)
		return err
	}
	var l uint32
	re.IoUint32(&l)
	(*s) = string(re.Array[re.index : re.index + uint64(l)])
	re.index += uint64(l)
	return err
}

func (re *Loader) IoHashmap(hm *Hashmap) error {
	arraySize := uint64(len(re.Array))
	// bug here
	var err error = nil
	if arraySize < (re.index + 4) {
		err = errors.New(ERR_ARRAY_TOO_SMALL)
		return err
	}
	var l uint32
	err = re.IoUint32(&l)
	if err != nil { return err }
	if (*hm) == nil {
		(*hm) = make(Hashmap)
	}
	for c := uint64(0); c < uint64(l); c++ {
		var s1 string
		err = re.IoString(&s1)
		if err != nil { return err }
		var s2 string
		err = re.IoString(&s2)
		if err != nil { return err }
		(*hm)[s1] = s2
	}
	return err
}

func (re *Loader) IoStringlist(sl *Stringlist) error {
	var err error = nil
	if uint64(len(re.Array)) < re.index + 4 {
		err = errors.New(ERR_ARRAY_TOO_SMALL)
		return err
	}
	var l uint32
	re.IoUint32(&l)
	for c := uint64(0); c < uint64(l); c++ {
		var s1 string
		re.IoString(&s1)
		(*sl) = append((*sl), s1)
	}
	return err
}
