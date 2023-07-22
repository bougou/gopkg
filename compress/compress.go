package compress

import (
	"compress/gzip"
	"io"
	"sync"
)

// CompressWithGzip takes an io.Reader as input and pipes
// it through a gzip.Writer returning an io.ReaderCloser containing the gzipped data.
// An error is returned if passing data to the gzip.Writer fails
func CompressWithGzip(data io.Reader) (io.ReadCloser, error) {
	pipeReader, pipeWriter := io.Pipe()
	gzipWriter := gzip.NewWriter(pipeWriter)

	rc := &ReadWaitCloser{
		pipeReader: pipeReader,
	}

	rc.wg.Add(1)
	var err error
	go func() {
		defer rc.wg.Done()

		_, err = io.Copy(gzipWriter, data)
		gzipWriter.Close()
		// subsequent reads from the read half of the pipe will
		// return no bytes and the error err, or EOF if err is nil.
		pipeWriter.CloseWithError(err)
	}()

	return pipeReader, err
}

type ReadWaitCloser struct {
	pipeReader *io.PipeReader
	wg         sync.WaitGroup
}

func (r *ReadWaitCloser) Close() error {
	err := r.pipeReader.Close()
	r.wg.Wait() // wait for the gzip goroutine finish
	return err
}
