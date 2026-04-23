package serialisers

import (
	"bytes"
	"testing"
	"reflect"
)

type Tester struct {
	i32a int32
	i32b int32
	i32c int32
	i32d int32
	i64a int64
	i64b int64
	u32a uint32
	u32b uint32
	u64a uint64
	u64b uint64
	f32a float32
	f32b float32
	f64a float64
	f64b float64
	s1 string
	s2 string
	hm1 Hashmap
	sl1 Stringlist
}

func (re *Tester) serialise(s Serialiser) error {
	var err error
	err = s.IoInt32(&re.i32a)
	if err != nil { return err }
	err = s.IoInt32(&re.i32b)
	if err != nil { return err }
	err = s.IoInt32(&re.i32c)
	if err != nil { return err }
	err = s.IoInt32(&re.i32d)
	if err != nil { return err }
	err = s.IoInt64(&re.i64a)
	if err != nil { return err }
	err = s.IoInt64(&re.i64b)
	if err != nil { return err }
	err = s.IoUint32(&re.u32a)
	if err != nil { return err }
	err = s.IoUint32(&re.u32b)
	if err != nil { return err }
	err = s.IoUint64(&re.u64a)
	if err != nil { return err }
	err = s.IoUint64(&re.u64b)
	if err != nil { return err }
	err = s.IoFloat32(&re.f32a)
	if err != nil { return err }
	err = s.IoFloat32(&re.f32b)
	if err != nil { return err }
	err = s.IoFloat64(&re.f64a)
	if err != nil { return err }
	err = s.IoFloat64(&re.f64b)
	if err != nil { return err }
	err = s.IoString(&re.s1)
	if err != nil { return err }
	err = s.IoString(&re.s2)
	if err != nil { return err }
	err = s.IoHashmap(&re.hm1)
	if err != nil { return err }
	err = s.IoStringlist(&re.sl1)
	if err != nil { return err }
	return nil
}

func testSizer() (uint32, Tester, error) {
	var size uint32 = 0
	hm := make(Hashmap)
	hm["team"] = "SysAdmins"
	hm["key"] = "FLB-51937"
	var sl []string
	sl = append(sl, "test1")
	sl = append(sl, "test2")
	tr := Tester{
		84, 69, 83, 84,
		-23548, -69375,
		12, 24, 4734, 2756,
		3.14, -1.2, -3.6, 3.1415926585,
		"var e byte = re.Array[re.index",
		" + uint64(add + (4 * mul))]",
		hm,
		sl,
	}
	sz := Sizer{}.Init(&size)
	err := tr.serialise(&sz)
	return size, tr, err
}

func testSaver(tr Tester, ob *bytes.Buffer, size uint32) ([]byte, error) {
	sv := Saver{}.Init(&size)
	err := tr.serialise(&sv)
	(*ob).Write(sv.Array)
	return sv.Array, err
}

func testLoader(tr Tester, ob *bytes.Buffer, size uint32, svArray *[]byte) (int, bool, []byte, Tester, error) {
	ld := Loader{}.Init(&size)
	n, err := (*ob).Read(ld.Array)
	var trLoader Tester
	err2 := trLoader.serialise(&ld)
	(*ob).Reset()
	(*ob).Write(ld.Array)
	return n, (err == nil), ld.Array, trLoader, err2
}

func TestSizer(t *testing.T) {
	size, _, err := testSizer()
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("Sized test struct at %d bytes.", size)
}

func TestSaver(t *testing.T) {
	size, tr, err := testSizer()
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("Sized test struct at %d bytes.", size)
	var outBuffer bytes.Buffer
	_, err = testSaver(tr, &outBuffer, size)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("Buffer: %q", outBuffer.String())
}

func TestLoader(t *testing.T) {
	size, tr, err := testSizer()
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("Sized test struct at %d bytes.", size)
	var outBuffer bytes.Buffer
	svArray, err := testSaver(tr, &outBuffer, size)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Logf("Buffer: %q", outBuffer.String())
	numBytesCopied, isCopiedOk, ldArray, trLoader, err := testLoader(tr, &outBuffer, size, &svArray)
	if err != nil {
		t.Errorf(err.Error())
	}
	if isCopiedOk {
		t.Logf("No error while copying Loader array")
	} else {
		t.Errorf("Error thrown while copying Loader array")
	}
	if uint32(numBytesCopied) == size {
		t.Logf("Correctly copied %d bytes.", numBytesCopied)
	} else {
		errStr := "Read wrong number of bytes."
		errStr += " Wanted %d, got %d."
		t.Errorf(errStr, size, numBytesCopied)
	}
	if reflect.DeepEqual(svArray, ldArray) {
		t.Logf("Saver and Loader arrays are identical.")
	} else {
		errStr := "Saver and Loader arrays differ."
		t.Errorf("%s", errStr)
	}
	t.Logf("Buffer: %q", outBuffer.String())
	if reflect.DeepEqual(tr, trLoader) {
		t.Logf("Post-Load data identical to Pre-Save data.")
	} else {
		errStr := "Post-Load data and Pre-Save data differ."
		t.Errorf("%s", errStr)
	}
}

func TestIsLittleEndian(t *testing.T) {
	// this is going to be run on LE systems only
	if isLittleEndian() {
		t.Logf("Endianness correctly detected.")
	} else {
		errStr := "Endianness detection failed."
		t.Errorf("%s", errStr)
	}
}

func TestWipe(t *testing.T) {
	var size uint32 = 2048
	sv := Saver{}.Init(&size)
	ld := Loader{}.Init(&size)
	sv.SecureWipe()
	ld.SecureWipe()
}
