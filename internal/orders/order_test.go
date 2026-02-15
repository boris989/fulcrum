package orders

import "testing"

func TestNewOrderSuccess(t *testing.T) {
	o, err := NewOrder(100)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if o.Amount() != 100 {
		t.Fatalf("got %d, want %d", o.Amount(), 100)
	}

	if o.ID() == "" {
		t.Fatal("id should not be empty")
	}
}

func TestNewOrderInvalidAmount(t *testing.T) {
	_, err := NewOrder(0)

	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestPaySuccess(t *testing.T) {
	o, _ := NewOrder(100)
	err := o.Pay()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if o.Status() != StatusPaid {
		t.Fatal("status not updated")
	}
}

func TestPayCancelled(t *testing.T) {
	o, _ := NewOrder(100)
	_ = o.Cancel()
	err := o.Pay()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCancelPaid(t *testing.T) {
	o, _ := NewOrder(100)
	_ = o.Pay()

	err := o.Cancel()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestOrderCreatedEvent(t *testing.T) {
	o, _ := NewOrder(100)

	events := o.PullEvents()

	if len(events) != 1 {
		t.Fatalf("got %d, want %d", len(events), 1)
	}

	if events[0].Name() != "OrderCreated" {
		t.Fatalf("got %s, want %s", events[0].Name(), "OrderCreated")
	}
}

func TestOrderPaidEvent(t *testing.T) {
	o, _ := NewOrder(100)
	_ = o.PullEvents() //clear created

	err := o.Pay()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	events := o.PullEvents()

	if len(events) != 1 {
		t.Fatalf("got %d, want %d", len(events), 1)
	}

	if events[0].Name() != "OrderPaid" {
		t.Fatalf("got %s, want %s", events[0].Name(), "OrderPaid")
	}
}

func TestPullEventsClears(t *testing.T) {
	o, _ := NewOrder(100)

	_ = o.PullEvents()
	events := o.PullEvents()

	if len(events) != 0 {
		t.Fatalf("got %d, want %d", len(events), 0)
	}
}
