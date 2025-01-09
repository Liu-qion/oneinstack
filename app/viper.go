package app

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Viper //
// 优先级: 命令行 > 环境变量 > 默认值
// Author [SliverHorn](https://github.com/SliverHorn)
func Viper(path ...string) *viper.Viper {
	var config string
	// 函数传递的可变参数的第一个值赋值于config
	config = "config.yaml"
	fmt.Printf("config的路径为%s\n", config)

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&ONE_CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&ONE_CONFIG); err != nil {
		panic(err)
	}
	return v
}
