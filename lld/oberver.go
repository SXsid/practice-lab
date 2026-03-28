// bad
package main

type handler struct {
	eventBus *EventBus
}

// bad one funton alot to do every funton is writien by diffrent developer and hwich orderpalce need to updaed on
func (h *handler) PlaceOrder(order any) {
	// save to db
	// nofity consumer
	// notify seller
	// genreate invoce
	// return cash back
	// update customer order and seller order book
	// regiet customer email to track updates
}

// good user observer that this funton should act as state change and the subcriber should sarte executing
type EventBus struct {
	Listerners map[string][]func(any)
}

func (e *EventBus) AddListner(event string, handler func(any)) {
	e.Listerners[event] = append(e.Listerners[event], handler)
}

func (e *EventBus) Emit(event string, data any) {
	for _, handler := range e.Listerners[event] {
		handler(data)
	}
}

func (h *handler) EventPlaceOrder(order any) {
	h.eventBus.Emit("ORDERPLACED", order)
}

type OrderHandler struct {
	eventBus *EventBus
}

func (o *OrderHandler) orderemialclient() {
	o.eventBus.AddListner("ORDERPLACED", func(a any) {
	})
}
