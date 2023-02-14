package yaml_encoder_test

import (
	"fmt"
	"testing"

	yaml_encoder "github.com/wetor/AnimeGo/third_party/yaml-encoder"
)

func TestEncoder(t *testing.T) {
	type DBConfig struct {
		Username string `yaml:"username" comment:"数据库用户名"`
		Server   struct {
			Host []string `yaml:"host"`
			Port struct {
				TestField string `yaml:"test_field" comment:"这是测试嵌套标签"`
			} `yaml:"port" comment_key:"PORT"`
		} `yaml:"server" comment:"服务器设置"`
		Password string `yaml:"password" comment:"密码"`
	}

	config := DBConfig{
		Password: "xxxxxx",
	}
	config.Server.Host = []string{"127.0.0.1", "127.0.0.2"}
	config.Server.Port.TestField = "55"

	encoder1 := yaml_encoder.NewEncoder(config,
		yaml_encoder.WithComments(yaml_encoder.CommentsOnHead),
		yaml_encoder.WithCommentsMap(map[string]string{
			"host": "主机名",
			"PORT": `端口号\n\r
- 端口号第一行
- 端口号第二行`,
		}),
	)
	content, _ := encoder1.EncodeDoc()
	fmt.Println("=============")
	fmt.Printf("%s\n", content)

	type Test struct {
		Username string `yaml:"username" `
		Server   struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port" `
		} `yaml:"server" `
		Password string `yaml:"password" `
	}

	config2 := Test{}
	encoder2 := yaml_encoder.NewEncoder(config2,
		yaml_encoder.WithComments(yaml_encoder.CommentsOnHead),
		yaml_encoder.WithCommentsMap(map[string]string{}),
	)
	content2, _ := encoder2.EncodeDoc()
	fmt.Println("=============")
	fmt.Printf("%s\n", content2)
	// Output:
	// # this is the username of database
	// username: root
	// # this is the password of database
	// password: xxxxxx
	// # 主机名
	// host: 127.0.0.1
	// # 端口号
	// port: 4444
}
