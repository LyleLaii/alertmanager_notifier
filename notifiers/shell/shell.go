package shell

import (
	"alertmanager_notifier/config"
	"alertmanager_notifier/log"
	"alertmanager_notifier/notifiers"
	"alertmanager_notifier/template"
	"fmt"
	"os/exec"
)

type Shell struct {
	name   string
	config *config.ShellConfig
	tmpl  *template.Template
	logger log.Logger
}


func New(name string, c *config.ShellConfig, tmpl *template.Template, logger log.Logger) (*Shell, error) {
	return &Shell{
		name:   name,
		config: c,
		tmpl: tmpl,
		logger: logger,
	}, nil
}

func (s Shell) Notify(am *notifiers.AlertMessage)  {
	for index, _ := range am.AlertInfo.Alerts {
		var argList []string
		for _, value := range s.config.Args {
			result, err := s.tmpl.ParseTmplString(index, value, am)
			if err != nil {
				s.logger.Error("shellExecutor", err)
			}
			argList = append(argList, result)
		}

		s.logger.Debug("shellExecutor", fmt.Sprintf("%+v, %+v", s.config.Command, argList))

		cmd := exec.Command(s.config.Command, argList...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			s.logger.Warn("shellExecutor", fmt.Sprintf("exec command return err: %+v", err))
		}
		s.logger.Info("shellExecutor", string(output))
	}

}