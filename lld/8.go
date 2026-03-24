package main

// wrong

func (h *Handler) HandleOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	itemID := r.URL.Query().Get("item_id")
	amount := 1000.0

	// handler has to know about ALL these services
	user := h.userService.GetUser(userID)
	item := h.inventoryService.CheckStock(itemID)
	payment := h.paymentService.Charge(user, amount)
	order := h.orderService.Create(user, item, payment)
	h.notifier.Send(user.Email, "Order confirmed!")
	h.inventoryService.DeductStock(itemID)
}

// facade pattern or a orchestion interface /class
func (h *Handler) HandleOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	itemID := r.URL.Query().Get("item_id")

	// calls ONE thing
	order, err := h.orderFacade.PlaceOrder(userID, itemID, 1000.0)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(order)
}

// right way handler not know about the all service we are using
type OrderFacade struct {
	userService      UserService
	inventoryService InventoryService
	paymentService   PaymentService
	orderService     OrderService
	notifier         Notifier
}

func (f *OrderFacade) PlaceOrder(userID string, itemID string, amount float64) (*Order, error) {
	// ALL the orchestration lives here
	user, err := f.userService.GetUser(userID)
	if err != nil {
		return nil, err
	}

	if err := f.inventoryService.CheckStock(itemID); err != nil {
		return nil, err
	}

	payment, err := f.paymentService.Charge(user, amount)
	if err != nil {
		return nil, err
	}

	order, err := f.orderService.Create(user, itemID, payment)
	if err != nil {
		return nil, err
	}

	f.inventoryService.DeductStock(itemID)
	f.notifier.Send(user.Email, "Order confirmed!")

	return order, nil
}
