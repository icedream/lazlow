package effects

import "time"

type LazlowDurationOption struct {
	humanFriendlyName string
	description       string
	defaultValue      time.Duration
	value             time.Duration
	min               time.Duration
	max               time.Duration
	step              time.Duration
}

func NewLazlowEncoderDurationOption(
	humanFriendlyName string,
	description string,
	defaultValue time.Duration,
	min time.Duration,
	max time.Duration,
	step time.Duration) LazlowOption {
	return &LazlowDurationOption{
		humanFriendlyName: humanFriendlyName,
		description:       description,
		defaultValue:      defaultValue,
		value:             defaultValue,
		min:               min,
		max:               max,
		step:              step,
	}
}

func (option *LazlowDurationOption) HumanFriendlyName() string {
	return option.humanFriendlyName
}

func (option *LazlowDurationOption) Description() string {
	return option.description
}

func (option *LazlowDurationOption) DefaultValue() interface{} {
	return option.defaultValue
}

func (option *LazlowDurationOption) TypedValue() time.Duration {
	return option.value
}

func (option *LazlowDurationOption) Value() interface{} {
	return option.value
}

func (option *LazlowDurationOption) Min() time.Duration {
	return option.min
}

func (option *LazlowDurationOption) Max() time.Duration {
	return option.max
}

func (option *LazlowDurationOption) SetValue(value interface{}) (err error) {
	var convertedValue time.Duration

	switch v := value.(type) {
	case string:
		convertedValue, err = time.ParseDuration(v)
		if err != nil {
			return
		}
	case time.Duration:
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
