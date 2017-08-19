package main

import "crypto/sha512"
import "hash"
import "bytes"
import "os"
import "io"
import "flag"
import "log"

func WriteBlock(f io.WriteSeeker, h hash.Hash, idx, bs int64) error {
	buf := new(bytes.Buffer)
	buf.Grow(int(bs))
	gw := h.Size()
	for buf.Len() < int(bs) {
		b := h.Sum(nil)
		g, err := buf.Write(b)
		if err != nil {
			log.Fatalln("Failed to write to block buffer")
		}
		if g != gw {
			log.Fatalf("Wrong number of bytes written to block buffer %d, expected %d", g, gw)
		}
		g, err = h.Write(b)
		if err != nil {
			log.Fatalln("Failed to write to hash")
		}
		if g != gw {
			log.Fatalf("Wrong number of bytes written to hash %d, expected %d", g, gw)
		}
	}
	if buf.Len() != int(bs) {
		log.Fatalf("Constructed a block of length %d, expected block of length %d", buf.Len(), bs)
	}
	_, err := f.Seek(idx*bs, io.SeekStart)
	if err != nil {
		log.Printf("Failed to seek to block %d (bs=%d)", idx, bs)
		return err
	}
	g, err := io.Copy(f, buf)
	if err != nil {
		log.Printf("Block %d (bs=%d) failed to write, returning error %+v", idx, bs, err)
		return err
	}
	if g == 0 {
		log.Println("EOF")
		return io.EOF
	}
	return nil
}

func ReadBlock(f io.ReadSeeker, h hash.Hash, idx, bs int64) error {
	buf := new(bytes.Buffer)
	buf.Grow(int(bs))
	gw := h.Size()
	for buf.Len() < int(bs) {
		b := h.Sum(nil)
		g, err := buf.Write(b)
		if err != nil {
			log.Fatalln("Failed to write to block buffer")
		}
		if g != gw {
			log.Fatalf("Wrong number of bytes written to block buffer %d, expected %d", g, gw)
		}
		g, err = h.Write(b)
		if err != nil {
			log.Fatalln("Failed to write to hash")
		}
		if g != gw {
			log.Fatalf("Wrong number of bytes written to hash %d, expected %d", g, gw)
		}
	}
	if buf.Len() != int(bs) {
		log.Fatalf("Constructed a block of length %d, expected block of length %d", buf.Len(), bs)
	}
	_, err := f.Seek(idx*bs, io.SeekStart)
	if err != nil {
		log.Printf("Failed to seek to block %d (bs=%d)", idx, bs)
		return err
	}
	blk := make([]byte, bs)
	g, err := f.Read(blk)
	if err == io.EOF {
		log.Printf("Reached end of file at block %d", idx)
		return err
	}
	if err != nil {
		log.Printf("Failed to read block %d (bs=%d)", idx, bs)
		return err
	}
	if g != int(bs) {
		log.Printf("Read incomplete block of length %d, expected length %d, at block %d", g, bs, idx)

	}
	if bytes.Compare(blk, buf.Bytes()) != 0 {
		log.Printf("Block %d (bs=%d) returned invalid data", idx, bs)
	}
	return nil
}

func main() {
	var seed, block string
	flag.StringVar(&seed, "seed", "random", "seed for unique pattern to write")
	flag.StringVar(&block, "open", "", "block device to open for writing and reading")

	var bs, startBlk int64
	flag.Int64Var(&bs, "bs", 512, "block size")
	flag.Int64Var(&startBlk, "start", 0, "start block")

	var checkOnly bool
	flag.BoolVar(&checkOnly, "checkonly", false, "skip write step")

	flag.Parse()

	if block == "" {
		log.Fatalln("'open' argument is required")
	}

	f, err := os.OpenFile(block, os.O_RDWR, 0600)
	if err != nil {
		log.Fatalf("os.OpenFile('%s',...) returned %v", block, err)
	}

	totalBytes, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		log.Fatalf("Failed to seek to end of file. os.Seek() returned %v", err)
	}

	h := sha512.New()
	if bs%int64(h.Size()) != 0 {
		log.Fatalf("block size %d is not an even multiple of hash size %d", bs, h.Size())
	}
	if totalBytes%bs != 0 {
		log.Fatalf("block size %d is not an even divisor of file size %d", bs, totalBytes)
	}

	totalBlocks := totalBytes / bs
	log.Printf("%s (%d bytes), blocksize %d (%d blocks), starting at block %d", block, totalBytes, bs, totalBlocks, startBlk)

	h.Write([]byte(seed))
	var blk int64
	if !checkOnly {
		blk = startBlk
		for blk < totalBlocks {
			err = WriteBlock(f, h, blk, bs)
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatalf("error writing to device: %q", err)
			}
			blk++
		}
		log.Printf("wrote %d blocks (bs=%d), reading and confirming...", blk, bs)
		h.Reset()
		h.Write([]byte(seed))
	}

	blk = startBlk
	for blk < totalBlocks {
		err = ReadBlock(f, h, blk, bs)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("error reading from device: %q", err)
		}
		blk++
	}
	log.Println("done")
	f.Close()
}
