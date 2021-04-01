package lazlow

type LazlowIntegerOption struct {
	humanFriendlyName string
	description       string
	defaultValue      int64
	value             int64
	min               int64
	max               int64
	step              int64
}

func NewLazlowEncoderIntegerOption(
	humanFriendlyName string,
	description string,
	defaultValue int64,
	min int64,
	max int64,
	step int64) LazlowOption {
	return &LazlowIntegerOption{
		humanFriendlyName: humanFriendlyName,
		description:       description,
		defaultValue:      defaultValue,
		value:             defaultValue,
		min:               min,
		max:               max,
		step:              step,
	}
}

func (option *LazlowIntegerOption) HumanFriendlyName() string {
	return option.humanFriendlyName
}

func (option *LazlowIntegerOption) Description() string {
	return option.description
}

func (option *LazlowIntegerOption) DefaultValue() interface{} {
	return option.defaultValue
}

func (option *LazlowIntegerOption) TypedValue() int64 {
	return option.value
}

func (option *LazlowIntegerOption) Value() interface{} {
	return option.value
}

func (option *LazlowIntegerOption) Min() int64 {
	return option.min
}

func (option *LazlowIntegerOption) Max() int64 {
	return option.max
}

func (option *LazlowIntegerOption) SetValue(value interface{}) (err error) {
	var convertedValue int64

	switch v := value.(type) {
	case uint8:
		convertedValue = int64(v)
	case uint16:
		convertedValue = int64(v)
	case uint32:
		convertedValue = int64(v)
	case int8:
		convertedValue = int64(v)
	case int16:
		convertedValue = int64(v)
	case int32:
		convertedValue = int64(v)
	case int64:
		convertedValue = v
	default:
		err = ErrIncompatibleValueType
		return
	}

	if convertedValue%option.step != 0 ||
		convertedValue < option.min ||
		convertedValue > option.max {
		err = ErrOutOfRange
		return
	}

	option.value = convertedValue
	return
}
