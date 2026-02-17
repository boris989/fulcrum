package orders

import "errors"
import "github.com/google/uuid"

type Status string

const (
	StatusNew       Status = "NEW"
	StatusPaid      Status = "PAID"
	StatusCancelled Status = "CANCELLED"
)

type Order struct {
	id      string
	amount  int64
	status  Status
	version int64

	events []Event
}

func NewOrder(amount int64) (*Order, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	o := &Order{
		id:     uuid.NewString(),
		amount: amount,
		status: StatusNew,
	}

	o.events = append(o.events, OrderCreated{
		OrderID: o.id,
		Amount:  o.amount,
	})

	return o, nil
}

func (o *Order) ID() string {
	return o.id
}

func (o *Order) Amount() int64 {
	return o.amount
}

func (o *Order) Status() Status {
	return o.status
}

func (o *Order) Version() int64 {
	return o.version
}

func (o *Order) SetVersion(version int64) {
	o.version = version
}

func (o *Order) Pay() error {
	if o.status == StatusCancelled {
		return errors.New("cannot pay cancelled order")
	}

	if o.status == StatusPaid {
		return errors.New("order already paid")
	}

	o.status = StatusPaid

	o.events = append(o.events, OrderPaid{
		OrderID: o.id,
	})

	return nil
}

func (o *Order) Cancel() error {
	if o.status == StatusPaid {
		return errors.New("cannot cancel paid order")
	}

	if o.status == StatusCancelled {
		return errors.New("order already cancelled")
	}

	o.status = StatusCancelled

	o.events = append(o.events, OrderCancelled{
		OrderID: o.id,
	})

	return nil
}

func (o *Order) PendingEvents() []Event {
	if len(o.events) == 0 {
		return nil
	}

	cp := make([]Event, len(o.events))
	copy(cp, o.events)

	return cp
}

func (o *Order) ClearEvents() {
	o.events = nil
}

func (o *Order) PullEvents() []Event {
	ev := o.PendingEvents()
	o.ClearEvents()
	return ev
}
