package cfg

import "github.com/spf13/viper"

// StrVal represents string value of config
type StrVal string

// Load loads value from given key
func (s *StrVal) Load(key string) {
	v := StrVal(viper.GetString(key))
	*s = v
}

// String convert StrVal to string type
func (s StrVal) String() string {
	return string(s)
}

type IntVal int

// Load loads configuration from viper using a key
func (i *IntVal) Load(key string) {
	v := IntVal(viper.GetInt(key))
	*i = v
}

// Int convert IntVal to int type
func (i IntVal) Int() int {
	return int(i)
}
