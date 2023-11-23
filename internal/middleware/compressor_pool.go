package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

// compressPWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type compressPWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

type gzWriterPool struct {
	gz sync.Pool
}

func GzipPoolMiddleware() func(http.Handler) http.Handler {
	gzp := &gzWriterPool{
		gz: sync.Pool{
			New: func() any {
				return gzip.NewWriter(io.Discard)
			},
		},
	}
	return gzp.gzipPoolMiddleware
}

func (p *gzWriterPool) newCompressPWriter(w http.ResponseWriter) *compressPWriter {
	gz, _ := p.gz.Get().(*gzip.Writer)
	gz.Reset(w)

	return &compressPWriter{
		w:  w,
		zw: gz,
	}
}

func (c *compressPWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressPWriter) Write(p []byte) (int, error) {
	// сжимаем данные только для контента с типами application/json и text/html
	ct := c.w.Header().Get("Content-Type")
	supportsCTypes := strings.Contains(ct, "application/json") || strings.Contains(ct, "text/html")

	if !supportsCTypes {
		return c.w.Write(p)
	}

	return c.zw.Write(p)
}

func (c *compressPWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressPWriter) Close() error {
	return c.zw.Close()
}

func (p *gzWriterPool) gzipPoolMiddleware(next http.Handler) http.Handler {
	gzipFn := func(w http.ResponseWriter, r *http.Request) {
		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции
		ow := w

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")

		if supportsGzip {
			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			cw := p.newCompressPWriter(w)
			// меняем оригинальный http.ResponseWriter на новый
			ow = cw
			defer cw.zw.Reset(io.Discard)
			defer p.gz.Put(cw.zw)
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
