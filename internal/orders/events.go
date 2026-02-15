package orders

type Event interface {
	Name() string
}

type OrderCreated struct {
	OrderID string
	Amount  int64
}

func (e OrderCreated) Name() string { return "OrderCreated" }

type OrderPaid struct {
	OrderID string
}

func (e OrderPaid) Name() string { return "OrderPaid" }

type OrderCancelled struct {
	OrderID string
}

func (e OrderCancelled) Name() string { return "OrderCancelled" }
