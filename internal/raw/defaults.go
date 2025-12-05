package raw

// Defaults captures SDK-wide request tweaks applied before invoking the gRPC stubs.
type Defaults struct {
	// User identifies the caller when requests don't explicitly set one.
	User string
}
