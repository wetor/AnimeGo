package websocket

import "github.com/wetor/AnimeGo/internal/logger"

type Command struct {
	Action        string `json:"action"`
	actionFuncMap map[string]func() error
}

func (c *Command) Init() {
	c.Action = ""
	c.actionFuncMap = make(map[string]func() error)
}

func (c *Command) SetActionFunc(action string, f func() error) {
	c.actionFuncMap[action] = f
}

func (c *Command) Execute() error {
	switch c.Action {
	case "pause":
		logger.PauseLogNotify()
	case "resume":
		logger.EnableLogNotify()
	case "terminate":
	}
	if f, ok := c.actionFuncMap[c.Action]; ok {
		err := f()
		if err != nil {
			return err
		}
	}
	return nil
}
