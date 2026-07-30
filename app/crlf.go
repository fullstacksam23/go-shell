package main

import "io"

type crlfWriter struct{ w io.Writer }

func (c crlfWriter) Write(p []byte) (int, error) {
	var out []byte
	for i := 0; i < len(p); i++ {
		if p[i] == '\n' && (i == 0 || p[i-1] != '\r') {
			out = append(out, '\r', '\n')
		} else {
			out = append(out, p[i])
		}
	}
	_, err := c.w.Write(out)
	if err != nil {
		return 0, err
	}
	return len(p), nil // report the original length so callers (e.g. cmd.Run) don't error
}
