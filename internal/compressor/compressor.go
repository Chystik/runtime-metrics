package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// compressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	// сжимаем данные только для контента с типами application/json и text/html
	ct := c.w.Header().Get("Content-Type")
	supportsCTypes := strings.Contains(ct, "application/json") || strings.Contains(ct, "text/html")

	if !supportsCTypes {
		return c.w.Write(p)
	}

	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// compressReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func GzipMiddleware(next http.Handler) http.Handler {
	gzipFn := func(w http.ResponseWriter, r *http.Request) {
		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции
		ow := w

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")

		if supportsGzip {
			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			cw := newCompressWriter(w)
			// меняем оригинальный http.ResponseWriter на новый
			ow = cw
			// не забываем отправить клиенту все сжатые данные после завершения middleware
			defer cw.Close()
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")

		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			r.Body = cr
			defer cr.Close()
		}

		// передаём управление хендлеру
		next.ServeHTTP(ow, r)

	}
	return http.HandlerFunc(gzipFn)
}

/*
import (
	"bufio"
	"compress/gzip"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
)

type compressor struct {
	level   int
	writers *sync.Pool
	readers *sync.Pool
}

// Compress is a middleware that compresses response body of json and text/plain
// content types to a gzip format. Compress level betwen 0 - 9,
// where 0 - no compressoin, 9 - best compression
func Compress(level int) func(http.Handler) http.Handler {
	c := newCompressor(level)
	return c.Handler
}

func newCompressor(level int) *compressor {
	return &compressor{
		level: level,
		writers: &sync.Pool{
			New: func() any {
				return encoderGzip(io.Discard, level)
			},
		},
		readers: &sync.Pool{
			New: func() any {
				gzr, _ := gzip.NewReader(nil)
				return gzr
			},
		},
	}
}

func encoderGzip(w io.Writer, level int) io.Writer {
	gw, err := gzip.NewWriterLevel(w, level)
	if err != nil {
		return nil
	}
	return gw
}

type compressWriter struct {
	http.ResponseWriter
	gzw          io.Writer
	compressable bool
}

func (cw *compressWriter) WriteHeader(statusCode int) {
	cw.Header().Set("Content-Encoding", "gzip")
	cw.ResponseWriter.WriteHeader(statusCode)
}

func (cw *compressWriter) Write(p []byte) (int, error) {

	// Compress data only for types: application/json и text/html
	ct := cw.Header().Get("Content-Type")
	supportsCTypes := strings.Contains(ct, "application/json") || strings.Contains(ct, "text/html")

	if !supportsCTypes {
		return cw.Write(p)
	}
	cw.compressable = true
	return cw.gzw.Write(p)
}

func (cw *compressWriter) writer() io.Writer {
	if cw.compressable {
		return cw.gzw
	} else {
		return cw.ResponseWriter
	}
}

type compressFlusher interface {
	Flush() error
}

func (cw *compressWriter) Flush() {
	if f, ok := cw.writer().(http.Flusher); ok {
		f.Flush()
	}

	// If the underlying writer has a compression flush signature,
	// call this Flush() method instead
	if f, ok := cw.writer().(compressFlusher); ok {
		f.Flush()

		// Also flush the underlying response writer
		if f, ok := cw.ResponseWriter.(http.Flusher); ok {
			f.Flush()
		}
	}
}

// verify Hijacker interface implementation
var _ http.Hijacker = (*compressWriter)(nil)

func (cw *compressWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := cw.writer().(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, errors.New("compress/middleware: http.Hijacker is unavailable on the writer")
}

func (cw *compressWriter) Push(target string, opts *http.PushOptions) error {
	if ps, ok := cw.writer().(http.Pusher); ok {
		return ps.Push(target, opts)
	}
	return errors.New("compress/middleware: http.Pusher is unavailable on the writer")
}

// Close closes gzip.Writer and sends data from buffer
func (cw *compressWriter) Close() error {
	if c, ok := cw.writer().(io.WriteCloser); ok {
		return c.Close()
	}
	return errors.New("compress/middleware: io.WriteCloser is unavailable on the writer")
}

// Interface for types that allow resetting io.Writers.
type ioResetterWriter interface {
	io.Writer
	Reset(w io.Writer)
}

type compressReader struct {
	r   io.ReadCloser
	gzr *gzip.Reader
}

func (cr compressReader) Read(p []byte) (n int, err error) {
	return cr.gzr.Read(p)
}

func (cr *compressReader) Close() error {
	if err := cr.r.Close(); err != nil {
		return err
	}
	return cr.gzr.Close()
}

func (c *compressor) Handler(next http.Handler) http.Handler {
	gzipFn := func(w http.ResponseWriter, r *http.Request) {

		// Original writer
		ow := w

		// Check that client can receive compressed data in gzip format
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")

		if supportsGzip {
			// Wrapp original http.ResponseWriter with new one that supports compression
			cw := &compressWriter{
				ResponseWriter: ow,
			}

			// Get gzip.Writer fomr pool
			gzw := c.writers.Get().(ioResetterWriter)
			gzw.Reset(ow)

			cw.gzw = gzw

			// Change the original writer to a new one that supports gzip
			ow = cw

			// Return writer to the pool
			defer c.writers.Put(gzw)

			// don't forget to send all compressed data to the client after completion the middleware
			defer cw.Close()
		}

		// we check that the client sent compressed data to the server in gzip format
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")

		if sendsGzip {
			// wrap the request body to io.Reader with decompression support
			cr := &compressReader{
				r: r.Body,
			}

			// Get gzip.Reader from pool
			gzr := c.readers.Get().(*gzip.Reader)
			gzr.Reset(r.Body)

			cr.gzr = gzr

			// Change r.Body to new one
			r.Body = cr

			// Return reader to the pool
			defer c.readers.Put(gzr)
			defer cr.Close()
		}
		next.ServeHTTP(ow, r)
	}
	return http.HandlerFunc(gzipFn)
}
*/
