package encoders_test

func genBlob(size int) []byte {
	blob := make([]byte, size)
	for i := 0; i < size; i++ {
		blob[i] = 65 // ascii 'A'
	}
	return blob
}

var bin = genBlob(1 << 24)
