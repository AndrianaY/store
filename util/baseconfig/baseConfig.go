package baseconfig

import (
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

func ReadConfigValues(keys interface{}) {
	var notSetVariables []string

	t := reflect.TypeOf(keys).Elem()
	if t == nil {
		panic("Unable to set Env variables")
	}

	ps := reflect.ValueOf(keys)
	s := ps.Elem()

	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)

		nameTag := fieldType.Tag.Get("name")
		if nameTag == "" {
			panic("`name` tag is not set for `settingsKeys` field `" + fieldType.Name + "`")
		}

		viper.BindEnv(nameTag)

		field := s.Field(i)
		if field.IsValid() && field.CanSet() {
			if defaultTag := fieldType.Tag.Get("default"); defaultTag != "" {
				viper.SetDefault(nameTag, defaultTag)
			}

			if !viper.IsSet(nameTag) {
				notSetVariables = append(notSetVariables, "`"+nameTag+"`")
				continue
			}

			switch field.Kind() {

			case reflect.Int32:
				field.SetInt(int64(viper.GetInt32(nameTag)))

			case reflect.String:
				field.SetString(viper.GetString(nameTag))

			default:
				panic("Unsupported value config type")
			}
		} else {
			panic("Unable to set config struct keys")
		}
	}

	if len(notSetVariables) > 0 {
		panic("Variables are not set: " + strings.Join(notSetVariables, ", "))
	}
}
