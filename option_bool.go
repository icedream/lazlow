package lazlow

type LazlowBoolOption struct {
	humanFriendlyName string
	description       string
	defaultValue      bool
	value             bool
}

func NewLazlowBoolOption(
	humanFriendlyName string,
	description string,
	defaultValue bool) LazlowOption {
	return &LazlowBoolOption{
		humanFriendlyName: humanFriendlyName,
		description:       description,
		defaultValue:      defaultValue,
		value:             defaultValue,
	}
}

func (option *LazlowBoolOption) HumanFriendlyName() string {
	return option.humanFriendlyName
}

func (option *LazlowBoolOption) Description() string {
	return option.description
}

func (option *LazlowBoolOption) DefaultValue() interface{} {
	return option.defaultValue
}

func (option *LazlowBoolOption) TypedValue() bool {
	return option.value
}

func (option *LazlowBoolOption) Value() interface{} {
	return option.value
}

func (option *LazlowBoolOption) SetValue(value interface{}) (err error) {
	convertedValue, ok := value.(bool)
	if !ok {
		err = ErrIncompatibleValueType
		return
	}

	option.value = convertedValue
	return
}
