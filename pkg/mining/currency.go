package mining

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	UAH Currency = "UAH"
)

func (currency Currency) String() string {
	return string(currency)
}
