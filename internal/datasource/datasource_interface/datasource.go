package datasource_interface

/*type Interval int
const (
	Day Interval = iota
	Week
	Month
)*/

type Datasource interface {
	ExtractCurrentPrices(currencyName string) (float32, error)
}
