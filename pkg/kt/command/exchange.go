package command

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/jonyhy96/kt-connect/pkg/kt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"

	"github.com/jonyhy96/kt-connect/pkg/kt/cluster"
	"github.com/jonyhy96/kt-connect/pkg/kt/connect"
	optionspkg "github.com/jonyhy96/kt-connect/pkg/kt/options"
	"github.com/jonyhy96/kt-connect/pkg/kt/util"
	urfave "github.com/urfave/cli"
)

// newExchangeCommand return new exchange command
func newExchangeCommand(cli kt.CliInterface, options *optionspkg.DaemonOptions, action ActionInterface) urfave.Command {
	return urfave.Command{
		Name:  "exchange",
		Usage: "exchange kubernetes deployment to local",
		Flags: []urfave.Flag{
			urfave.StringFlag{
				Name:        "expose",
				Usage:       "expose port [port] or [remote:local]",
				Destination: &options.ExchangeOptions.Expose,
			},
		},
		Action: func(c *urfave.Context) error {
			if options.Debug {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			}
			if err := combineKubeOpts(options); err != nil {
				return err
			}
			exchange := c.Args().First()
			expose := options.ExchangeOptions.Expose

			if len(exchange) == 0 {
				return errors.New("exchange is required")
			}
			if len(expose) == 0 {
				return errors.New("-expose is required")
			}
			return action.Exchange(exchange, cli, options)
		},
	}
}

//Exchange exchange kubernetes workload
func (action *Action) Exchange(exchange string, cli kt.CliInterface, options *optionspkg.DaemonOptions) error {
	ch := SetUpCloseHandler(cli, options, "exchange")

	kubernetes, err := cluster.Create(options.KubeConfig)
	if err != nil {
		return err
	}
	randomString := util.RandomString(5)
	app, err := kubernetes.Deployment(exchange, options.Namespace)
	if err != nil {
		return err
	}

	port, err := strconv.Atoi(options.ExchangeOptions.Expose)
	if err != nil {
		log.Error().Msgf("transform arguments --expose %s to int error %v", options.ExchangeOptions.Expose, err)
		return err
	}

	workload := app.GetName() + "-kt-" + strings.ToLower(randomString)

	sideCar := core.Container{
		Name:    "sidecar" + "-gor-" + strings.ToLower(randomString),
		Image:   options.SideCarImage,
		Command: []string{"gor"},
		Args: []string{
			"--input-raw", ":" + options.ExchangeOptions.Expose,
			"--output-http", "http://" + workload + ":" + options.ExchangeOptions.Expose,
		},
		SecurityContext: &core.SecurityContext{
			Capabilities: &core.Capabilities{
				Add: []core.Capability{"NET_ADMIN"},
			},
		},
	}
	app.Spec.Template.Spec.Containers = append(app.Spec.Template.Spec.Containers, sideCar)
	log.Info().Msgf("inject gor sideCar to deployment")
	updated, err := kubernetes.UpdateDeployment(options.Namespace, app)
	if err != nil {
		log.Error().Msgf("inject gor sideCar to deployment error %v", err)
		return err
	}

	app = updated // update app to app with sideCar version

	// record origin deployment for clean up
	options.RuntimeOptions.Patch = &optionspkg.Patch{
		DeploymentName: app.Name,
		SideCar:        sideCar,
	}
	// record context in order to remove after command exit
	options.RuntimeOptions.Origin = app.GetName()
	options.RuntimeOptions.Replicas = *app.Spec.Replicas

	shadowPodLabels := getExchangeLabels(options.Labels, workload, app)

	podIP, podName, sshcm, credential, err := kubernetes.GetOrCreateShadow(workload, options.Namespace, options.Image, shadowPodLabels, options.Debug, false)
	log.Info().Msgf("create exchange shadow %s in namespace %s", workload, options.Namespace)

	if err != nil {
		return err
	}

	service, err := kubernetes.CreateService(workload, options.Namespace, port, shadowPodLabels)
	if err != nil {
		log.Error().Msgf("create exchange shadow's service error %v", err)
		return err
	}
	log.Info().Msgf("create exchange shadow's service %s in namespace %s", service.Name, options.Namespace)

	// record data
	options.RuntimeOptions.Shadow = workload
	options.RuntimeOptions.SSHCM = sshcm
	options.RuntimeOptions.Service = workload

	shadow := connect.Create(options)
	if err = shadow.Inbound(options.ExchangeOptions.Expose, podName, podIP, credential); err != nil {
		return err
	}

	// watch background process, clean the workspace and exit if background process occur exception
	go func() {
		<-util.Interrupt()
		CleanupWorkspace(cli, options)
		os.Exit(0)
	}()
	s := <-ch
	log.Info().Msgf("Terminal Signal is %s", s)

	return nil
}

func getExchangeLabels(customLabels string, workload string, origin *v1.Deployment) map[string]string {
	labels := map[string]string{
		"kt":           workload,
		"kt-component": "exchange",
		"control-by":   "kt",
	}
	// extra labels must be applied after origin labels
	for k, v := range util.String2Map(customLabels) {
		labels[k] = v
	}
	splits := strings.Split(workload, "-")
	labels["version"] = splits[len(splits)-1]
	return labels
}
