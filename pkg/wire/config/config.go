package config

import "mals/pkg/config"

type Config struct {
	Logs      []Log      `mapstructure:"logs"`
	Models    []Model    `mapstructure:"models"`
	Lsps      []Lsp      `mapstructure:"lsps"`
	Handlers  []Handler  `mapstructure:"handlers"`
	Listeners []Listener `mapstructure:"listeners"`
}

func (o *Config) Wire(c *config.Config) error {
	if c.Logs != nil {
		o.Logs = make([]Log, len(c.Logs))
		for i, log := range c.Logs {
			wired := Log{}
			if err := wired.Wire(log); err != nil {
				return err
			}
			o.Logs[i] = wired
		}
	}

	if c.Models != nil {
		o.Models = make([]Model, len(c.Models))
		for i, model := range c.Models {
			wired := Model{}
			if err := wired.Wire(model); err != nil {
				return err
			}
			o.Models[i] = wired
		}
	}

	if c.Lsps != nil {
		o.Lsps = make([]Lsp, 0, len(c.Lsps))
		for i, lsp := range c.Lsps {
			wired := Lsp{}
			if err := wired.Wire(lsp); err != nil {
				return err
			}
			o.Lsps[i] = wired
		}
	}

	if c.Handlers != nil {
		o.Handlers = make([]Handler, 0, len(c.Handlers))
		for i, handler := range c.Handlers {
			wired := Handler{}
			if err := wired.Wire(handler); err != nil {
				return err
			}
			o.Handlers[i] = wired
		}
	}

	if c.Listeners != nil {
		o.Listeners = make([]Listener, 0, len(c.Listeners))
		for i, listener := range c.Listeners {
			wired := Listener{}
			if err := wired.Wire(listener); err != nil {
				return err
			}
			o.Listeners[i] = wired
		}
	}

	return nil
}

func (o *Config) Unwire() (*config.Config, error) {
	c := &config.Config{}

	if o.Logs != nil {
		c.Logs = make([]*config.Log, len(o.Logs))
		for i, log := range o.Logs {
			unwired, err := log.Unwire()
			if err != nil {
				return nil, err
			}
			c.Logs[i] = unwired
		}
	}

	if o.Models != nil {
		c.Models = make([]*config.Model, len(o.Models))
		for i, model := range o.Models {
			unwired, err := model.Unwire()
			if err != nil {
				return nil, err
			}
			c.Models[i] = unwired
		}
	}

	if o.Lsps != nil {
		c.Lsps = make([]*config.Lsp, len(o.Lsps))
		for i, lsp := range o.Lsps {
			unwired, err := lsp.Unwire()
			if err != nil {
				return nil, err
			}
			c.Lsps[i] = unwired
		}
	}

	if o.Handlers != nil {
		c.Handlers = make([]*config.Handler, len(o.Handlers))
		for i, handler := range o.Handlers {
			unwired, err := handler.Unwire()
			if err != nil {
				return nil, err
			}
			c.Handlers[i] = unwired
		}
	}

	if o.Listeners != nil {
		c.Listeners = make([]*config.Listener, len(o.Listeners))
		for i, listener := range o.Listeners {
			unwired, err := listener.Unwire()
			if err != nil {
				return nil, err
			}
			c.Listeners[i] = unwired
		}
	}

	return c, nil
}
