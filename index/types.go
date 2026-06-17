package index

type Header struct {
	Signature  [4]byte
	Version    uint32
	EntryCount uint32
}

type Entry struct {
	Size uint32

	SHA1 [20]byte

	Path [256]byte
}
