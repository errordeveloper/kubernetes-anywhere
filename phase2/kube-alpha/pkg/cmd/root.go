package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "kube",
	Short: "deploy and manage kubernetes clusters",
	Long: `
Run me on servers to transform them into a kubernetes cluster suitable either
for tire-kicking or production usage.

This tool assumes you already have some servers running Linux and Docker. If
you want to provision servers and have kubernetes automatically installed on
them using this tool, see https://kubernetes.io/docs/XXX.

I can deploy masters (where the Kubernetes control plane runs), and workers
(where your containers get deployed to), or servers that do both at the same
time (useful for smaller clusters). Servers can either be of type "worker",
"master", or "master-and-worker".

All kubernetes components will be deployed in containers.  An etcd cluster will
be created for you.
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

type Config struct {
	disco  string
	net    string
	pki    string
	dryRun bool
}

// TODO log everything really nicely, so that the user can see what happened.

func init() {
	config := Config{}
	RootCmd.AddCommand(NewCmdInit(os.Stdout, &config))
	RootCmd.AddCommand(NewCmdJoin(os.Stdout, &config))
	RootCmd.AddCommand(NewCmdToolbox(os.Stdout, &config))

	RootCmd.PersistentFlags().StringVarP(&config.disco, "disco", "", "weave",
		`which service discovery mechanism to use for kubernetes
bootstrap (choose from "weave", "dns", "token",
"consul", default: "weave").`)
	RootCmd.PersistentFlags().StringVarP(&config.net, "net", "", "weave",
		`which pod network to create (choose from "weave",
"flannel", default: "weave")`)
	RootCmd.PersistentFlags().StringVarP(&config.pki, "pki", "", "auto",
		`certificate provider, default "auto" to ask the
discovery mechanism to bootstrap certs for you when
you "init" (chose from "vault", "amazon-cm", "containers",
"token").`)
	RootCmd.PersistentFlags().BoolVar(&config.dryRun,
		"dry-run", false, "Dry run. Useful for understanding what would be done.")

}
