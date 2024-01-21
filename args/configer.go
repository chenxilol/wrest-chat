package args

import (
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/opentdp/go-helper/logman"
)

type Configer struct {
	Koanf  *koanf.Koanf
	Parser *yaml.YAML
	File   string
}

type ConfigData struct {
	Bot *IBot
	LLM *ILLM
	Log *ILog
	Web *IWeb
	Wcf *IWcf
}

func NewConfiger() *Configer {

	c := &Configer{
		Koanf: koanf.NewWithConf(koanf.Conf{
			StrictMerge: true,
			Delim:       ".",
		}),
		Parser: yaml.Parser(),
		File:   "config.yml",
	}

	if len(os.Args) > 1 {
		c.File = os.Args[1]
	}

	return c

}

func LoadConfig() error {

	c := NewConfiger()

	logman.Info("load config", "file", c.File)

	// 文件不存在
	if _, err := os.Stat(c.File); os.IsNotExist(err) {
		logman.Warn("load config", "skip", c.File)
		return nil // 忽略错误
	}

	// 从文件读入参数
	err := c.Koanf.Load(file.Provider(c.File), c.Parser)
	if err != nil {
		logman.Error("load config", "error", err)
		return err
	}

	// 将参数写入内存
	c.Koanf.Unmarshal("Bot", Bot)
	c.Koanf.Unmarshal("LLM", LLM)
	c.Koanf.Unmarshal("Log", Log)
	c.Koanf.Unmarshal("Web", Web)
	c.Koanf.Unmarshal("Wcf", Wcf)

	return nil

}

func SaveConfig() error {

	c := NewConfiger()

	logman.Info("save config", "file", c.File)

	// 从内存读入参数
	obj := ConfigData{Bot, LLM, Log, Web, Wcf}
	err := c.Koanf.Load(structs.Provider(obj, ""), nil)
	if err != nil {
		logman.Error("load struct", "error", err)
		return err
	}

	// 序列化参数信息
	buf, err := c.Koanf.Marshal(c.Parser)
	if err != nil {
		logman.Error("save config", "error", err)
		return err
	}

	// 将参数写入文件
	err = os.WriteFile(c.File+"-bak", buf, 0644)
	if err != nil {
		logman.Error("save config", "error", err)
		return err
	}

	return nil

}
