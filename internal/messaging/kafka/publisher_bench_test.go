package kafka

import "testing"

func BenchmarkBuildMessage(b *testing.B) {

	for i := 0; i < b.N; i++ {
		_ = buildMessage("order-1", "order-1", []byte(`{"event": "created"}`))
	}
}
