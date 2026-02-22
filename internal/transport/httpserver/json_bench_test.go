package httpserver

import (
	"encoding/json"
	"testing"
)

func BenchmarkMarshalOrderResponse(b *testing.B) {
	resp := createOrderResponse{
		ID: "abc123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(resp)
	}
}
