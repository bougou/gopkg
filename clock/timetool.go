package clock

const (
	// RFC3339Milli follows RFC3339 format with milliseconds for precision.
	// Official defines time.RFC3339 with seconds precision and
	// RFC3339Nano with nanoseconds precision, but lacks milliseconds precision
	// RFC3339     = "2006-01-02T15:04:05Z07:00"
	// RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	RFC3339Milli = "2006-01-02T15:04:05.999Z07:00"
)
