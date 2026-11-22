package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	PWMPinRight      string `mapstructure:"pwm_pin_right"`
	ControlPinRight1 string `mapstructure:"control_pin_right_1"`
	ControlPinRight2 string `mapstructure:"control_pin_right_2"`

	PWMPinLeft      string `mapstructure:"pwm_pin_left"`
	ControlPinLeft1 string `mapstructure:"control_pin_left_1"`
	ControlPinLeft2 string `mapstructure:"control_pin_left_2"`

	PWMFrequency float64       `mapstructure:"pwm_frequency"`
	MoveDuration time.Duration `mapstructure:"move_duration"`

	JoystickID int `mapstructure:"joystick_id"`
}

func Init() (*Config, error) {
	viper.AddConfigPath("./cfg/")

	viper.SetConfigName("app")
	viper.SetConfigType("env")
	err := viper.MergeInConfig()
	if err != nil {
		return nil, err
	}

	viper.AutomaticEnv()

	config := &Config{}

	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
