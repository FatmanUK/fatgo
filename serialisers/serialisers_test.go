package serialisers

import (
	"testing"
	"encoding/base64"
)

type Config struct {
	version float32
	secret string
	keys string
	wisp string
	bin_env string
	bin_agent string
	bin_add string
}

// TODO Improve this. Doesn't test all serialise functions.
func (re *Config) serialise(s Serialiser) {
	s.IoF(&re.version)
	s.IoS(&re.secret)
	s.IoS(&re.keys)
	s.IoS(&re.wisp)
	s.IoS(&re.bin_env)
	s.IoS(&re.bin_agent)
	s.IoS(&re.bin_add)
}

const test_data = "PkzMzQAAABF2YXVsdHMvc2VjcmV0LnR4dAAAAA92YXVsdHMv" +
                  "a2V5cy55bWwAAAASL2Rldi9zaG0vd2lzcC5iYXNoAAAADC91" +
                  "c3IvYmluL2VudgAAABIvdXNyL2Jpbi9zc2gtYWdlbnQAAAAQ" +
                  "L3Vzci9iaW4vc3NoLWFkZA=="
/*
Decodes to:
00000000: 3e4c cccd 0000 0011 7661 756c 7473 2f73  >L......vaults/s
00000010: 6563 7265 742e 7478 7400 0000 0f76 6175  ecret.txt....vau
00000020: 6c74 732f 6b65 7973 2e79 6d6c 0000 0012  lts/keys.yml....
00000030: 2f64 6576 2f73 686d 2f77 6973 702e 6261  /dev/shm/wisp.ba
00000040: 7368 0000 000c 2f75 7372 2f62 696e 2f65  sh..../usr/bin/e
00000050: 6e76 0000 0012 2f75 7372 2f62 696e 2f73  nv..../usr/bin/s
00000060: 7368 2d61 6765 6e74 0000 0010 2f75 7372  sh-agent..../usr
00000070: 2f62 696e 2f73 7368 2d61 6464            /bin/ssh-add
*/

func decode_test_data() ([]byte, error) {
	// decode test data
	b := base64.StdEncoding
	rawBytes, err := b.DecodeString(test_data)
	return rawBytes, err
}

func TestLoader(t *testing.T) {
	rawBytes, err := decode_test_data()
	if err != nil {
		t.Fatal("Failed decoding data blob.")
	}
	data := Config{}
	// Should the buffer be the responsibility of the Model object? Surely the Serialiser object?
	data.serialise(&Loader{Array: &rawBytes})
	//zero_buffer(&rawBytes)
	if data.version != 0.2 {
		t.Error("IoF returned wrong value")
	}
	if data.secret != "vaults/secret.txt" {
		t.Error("IoS returned wrong value")
	}
}

/*
func TestSaver(t *testing.T) {
	rawBytes, err := decode_test_data()
	if err != nil {
		t.Fatal("Failed decoding data blob.")
	}


func (re *Control) save_config() {
	re.view.log(
		DEBUG,
		fmt.Sprintf("Saving config: %s", re.conf_name))
	var buf_size uint64 = 0
	re.serialise(&serialisers.Sizer{&buf_size})
	buf := make([]byte, buf_size)
	// Should the buffer be the responsibility of the Model object?
	re.serialise(&serialisers.Saver{Array: &buf})
	write_binary_file(&buf, re.conf_name, re.view)
	zero_buffer(&buf)
}



}
*/

//func TestSizer(t *testing.T) {
//}

