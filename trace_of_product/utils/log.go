package utils

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func LogMsg(keys, values []string) error {
	fileds := logrus.Fields{}
	if len(keys) != len(values) {
		logrus.Errorf("输入不符合规范，keys 长度为 %d values 长度为 %d，两者不相同", len(keys), len(values))
		return fmt.Errorf("输入不符合规范，keys 长度为 %d values 长度为 %d，两者不相同", len(keys), len(values))
	}
	for i := 0; i < len(keys); i++ {
		fileds[keys[i]] = values[i]
	}
	logrus.WithFields(
		fileds,
	).Println()
	return nil
}

func LogWarn(keys, values []string) error {
	fileds := logrus.Fields{}
	if len(keys) != len(values) {
		logrus.Errorf("输入不符合规范，keys 长度为 %d values 长度为 %d，两者不相同", len(keys), len(values))
		return fmt.Errorf("输入不符合规范，keys 长度为 %d values 长度为 %d，两者不相同", len(keys), len(values))
	}
	for i := 0; i < len(keys); i++ {
		fileds[keys[i]] = values[i]
	}
	logrus.WithFields(
		fileds,
	).Warnln()
	return nil
}
func LogError(keys, values []string) error {
	fileds := logrus.Fields{}
	if len(keys) != len(values) {
		logrus.Errorf("输入不符合规范，keys 长度为 %d values 长度为 %d，两者不相同", len(keys), len(values))
		return fmt.Errorf("输入不符合规范，keys 长度为 %d values 长度为 %d，两者不相同", len(keys), len(values))
	}
	for i := 0; i < len(keys); i++ {
		fileds[keys[i]] = values[i]
	}
	logrus.WithFields(
		fileds,
	).Errorln()
	return nil
}
func LogDebug(keys, values []string) error {
	fileds := logrus.Fields{}
	if len(keys) != len(values) {
		logrus.Errorf("输入不符合规范，keys 长度为 %d values 长度为 %d，两者不相同", len(keys), len(values))
		return fmt.Errorf("输入不符合规范，keys 长度为 %d values 长度为 %d，两者不相同", len(keys), len(values))
	}
	for i := 0; i < len(keys); i++ {
		fileds[keys[i]] = values[i]
	}
	logrus.WithFields(
		fileds,
	).Debugln()
	return nil
}
