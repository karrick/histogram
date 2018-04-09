package gorill

const (
	largeBufSize = 8192 // large enough to force bufio.Writer to flush
	smallBufSize = 64
)

var (
	largeBuf []byte
	smallBuf []byte
)

func init() {
	newBuf := func(size int) []byte {
		buf := make([]byte, size)
		for i := range buf {
			buf[i] = byte(i % 256)
		}
		return buf
	}
	largeBuf = newBuf(largeBufSize)
	smallBuf = newBuf(smallBufSize)
}
