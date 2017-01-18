package simplelog

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

const minStackTraceBuffer = 1 << 14 // 16 KiB
const maxStackTraceBuffer = 1 << 24 // 16 MiB

func init() {
	c := make(chan os.Signal, 1)
	go func() {
		staticBuf := make([]byte, minStackTraceBuffer)
		ms := &runtime.MemStats{}

		for _ = range c {
			now := time.Now()

			// allocate space for formatted stack trace
			var buf = staticBuf
			bufSizeEstimate := runtime.NumGoroutine() << 11 // 2 KiB per goroutine
			if bufSizeEstimate > minStackTraceBuffer {
				if bufSizeEstimate > maxStackTraceBuffer {
					bufSizeEstimate = maxStackTraceBuffer
				}
				buf = make([]byte, bufSizeEstimate)
			}

			// dump out stack traces
			n := runtime.Stack(buf, true)
			fmt.Fprintf(os.Stderr, "%s: %d Stack traces:\n%s\n", now, n, string(buf[:n]))

			// dump out memory stats
			runtime.ReadMemStats(ms)
			fmt.Fprintf(os.Stderr, "Memory stats: %s/%s (%s total), %s across %d GCs\n\n",
				FormatByteCount(ms.Alloc), FormatByteCount(ms.Sys), FormatByteCount(ms.TotalAlloc),
				time.Duration(ms.PauseTotalNs), ms.NumGC)
		}
	}()
	signal.Notify(c, syscall.SIGQUIT)
}

func FormatByteCount(n uint64) string {
	if n < 1024 {
		return fmt.Sprintf("%d bytes", n)
	}
	v := float64(n) / 1024.0
	if v < 1024.0 {
		return fmt.Sprintf("%.2f KiB", v)
	}
	v /= 1024.0
	if v < 1024.0 {
		return fmt.Sprintf("%.2f MiB", v)
	}
	v /= 1024.0
	if v < 1024.0 {
		return fmt.Sprintf("%.2f GiB", v)
	}
	v /= 1024.0
	if v < 1024.0 {
		return fmt.Sprintf("%.2f TiB", v)
	}
	v /= 1024.0
	return fmt.Sprintf("%.2f PiB", v)
}
