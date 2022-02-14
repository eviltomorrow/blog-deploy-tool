package cmd

import (
	"context"
	"log"
	"path/filepath"
	"time"

	"github.com/eviltomorrow/blog-deploy-tool/pkg/client"
	"github.com/eviltomorrow/blog-deploy-tool/pkg/pb"
	"github.com/eviltomorrow/blog-deploy-tool/pkg/system"
	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use:   "blog-deploy",
	Short: "",
	Long:  "  \r\nblog-deploy is a tool for blog html publish",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cfg.Load(path, nil); err != nil {
			log.Fatalf("[Fatal] Load config failure, nest error: %v\r\n", err)
		}
		creds, err := client.WithTLS(ServerName, filepath.Join(system.Pwd, "certs", "ca.crt"), filepath.Join(system.Pwd, "certs", "client.crt"), filepath.Join(system.Pwd, "certs", "client.pem"))
		if err != nil {
			log.Fatalf("With TLS failure, nest error: %v\r\n", err)
		}

		stub, close, err := client.NewCmd(cfg.Server.Host, cfg.Server.Port, creds, 10*time.Second)
		if err != nil {
			log.Fatalf("New cmd client failure, nest error: %v\r\n", err)
		}
		defer close()

		output, err := stub.Execute(context.Background(), &pb.Input{
			Text:    cfg.Cmd.Text,
			Timeout: cfg.Cmd.Timeout,
		})
		if err != nil {
			log.Fatalf("Execute cmd failure, nest error: %v\r\n", err)
		}
		log.Printf("%s", output.Data)
	},
}

func init() {
	clientCmd.CompletionOptions = cobra.CompletionOptions{
		DisableDefaultCmd: false,
	}
	clientCmd.Flags().StringVarP(&path, "config", "c", "config.toml", "blog-deploy's config file")
}

func NewClient() {
	cobra.CheckErr(clientCmd.Execute())
}
