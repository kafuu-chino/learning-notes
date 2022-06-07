package apollo

import (
	"testing"

	remote "github.com/shima-park/agollo/viper-remote"
	"github.com/spf13/viper"
)

type Config struct {
	Test     string `mapstructure:"test"`
	DBConfig string `mapstructure:"db"`
}
type DBConfig struct {
	Host string `mapstructure:"host"` // db.host
}

func TestApollo(t *testing.T) {
	// 设置配置id
	remote.SetAppID("1")

	v := viper.New()

	// 设置配置类型
	v.SetConfigType("prop")
	err := v.AddRemoteProvider("apollo", "http://192.168.178.128:8080/", "application")
	if err != nil {
		t.Error(err)
		return
	}

	err = v.ReadRemoteConfig()
	if err != nil {
		t.Error(err)
	}

	// 直接反序列化到结构体中
	var conf Config
	err = v.Unmarshal(&conf)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%+v\n", conf)

	// 获取所有key，所有配置
	t.Log("AllKeys", v.AllKeys(), "AllSettings", v.AllSettings())
}
