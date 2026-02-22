package orders

import "testing"

func BenchmarkNewOrder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewOrder(1000)
	}
}

func BenchmarkPay(b *testing.B) {
	order, _ := NewOrder(1000)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = order.Pay()
	}
}
