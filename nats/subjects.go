package nats

type CustomerSubject string
type OrderSubject string
type PaymentSubject string
type StoreSubject string
type UserSubject string

const (
	CustomerAddressDeleted CustomerSubject = "address:deleted"
	CustomerDeleted        CustomerSubject = "customer:deleted"
	OrderCreate            OrderSubject    = "order:create"
	OrderCreated           OrderSubject    = "order:created"
	OrderStatus            OrderSubject    = "order:status"
	PaymentCancel          PaymentSubject  = "payment:cancel"
	PaymentCreate          PaymentSubject  = "payment:create"
	PaymentUpdate          PaymentSubject  = "payment:update"
	StoreBook              StoreSubject    = "store:book"
	StoreBooked            StoreSubject    = "store:booked"
	StorePaid              StoreSubject    = "store:paid"
	StorePayment           StoreSubject    = "store:payment"
	UserDeleted            UserSubject     = "user:deleted"
)

func GetCustomerSubjects() []string {
	return []string{
		string(CustomerDeleted),
		string(CustomerAddressDeleted),
	}
}

func GetOrderSubjects() []string {
	return []string{
		string(OrderCreate),
		string(OrderCreated),
		string(OrderStatus),
	}
}

func GetPaymentSubjects() []string {
	return []string{
		string(PaymentCancel),
		string(PaymentCreate),
		string(PaymentUpdate),
	}
}

func GetStoreSubjects() []string {
	return []string{
		string(StoreBook),
		string(StoreBooked),
		string(StorePaid),
		string(StorePayment),
	}
}

func GetUserSubjects() []string {
	return []string{
		string(UserDeleted),
	}
}
