package exec

import (
	"github.com/jonyhy96/kt-connect/pkg/kt/exec/kubectl"
	"github.com/jonyhy96/kt-connect/pkg/kt/exec/ssh"
	"github.com/jonyhy96/kt-connect/pkg/kt/exec/sshuttle"
)

// CliInterface ...
type CliInterface interface {
	Kubectl() kubectl.CliInterface
	SSHUttle() sshuttle.CliInterface
	SSH() ssh.CliInterface
}

// Cli ...
type Cli struct {
	KubeOptions []string
}

// Kubectl ...
func (c *Cli) Kubectl() kubectl.CliInterface {
	return &kubectl.Cli{KubeOptions: c.KubeOptions}
}

// SSHUttle ...
func (c *Cli) SSHUttle() sshuttle.CliInterface {
	return &sshuttle.Cli{}
}

// SSH ...
func (c *Cli) SSH() ssh.CliInterface {
	return &ssh.Cli{}
}
