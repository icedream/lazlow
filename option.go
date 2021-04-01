package lazlow

type LazlowOption interface {
	HumanFriendlyName() string
	Description() string
	DefaultValue() interface{}
	SetValue(value interface{}) (err error)
	Value() interface{}
}
