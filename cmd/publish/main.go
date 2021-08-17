package publish

import (
	"errors"
	"fmt"
	"github.com/datreeio/datree/bl/files"
	"github.com/datreeio/datree/pkg/cliClient"
	"strings"

	"github.com/datreeio/datree/bl/messager"
	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/spf13/cobra"
)

type Messager interface {
	LoadVersionMessages(messages chan *messager.VersionMessage, cliVersion string)
}

type Printer interface {
	PrintMessage(messageText string, messageColor string)
}

type Reader interface {
	FilterFiles(paths []string) ([]string, error)
}

type LocalConfig interface {
	GetLocalConfiguration() (*localConfig.ConfigContent, error)
}

type PublishCommandContext struct {
	CliVersion  string
	LocalConfig LocalConfig
	Messager    Messager
	Printer     Printer
	Reader      Reader
	CliClient   *cliClient.CliClient
}

func New(ctx *PublishCommandContext) *cobra.Command {
	publishCommand := &cobra.Command{
		Use:   "publish <fileName>",
		Short: "Publish policies configuration for given <fileName>.",
		Long:  "Publish policies configuration for given <fileName>. Input should be the path to the Policy as Code yaml configuration file ",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				errMessage := "Requires 1 arg\n"
				cmd.Usage()
				return fmt.Errorf(errMessage)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error = nil
			defer func() {
				if err != nil {
					ctx.Printer.PrintMessage(strings.Join([]string{"\n", err.Error(), "\n"}, ""), "error")
				}
			}()

			err = publish(ctx, args[0])
			return err
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	return publishCommand
}

func publish(ctx *PublishCommandContext, path string) error {
	return errors.New("some error")

	jsonContent, err := files.ExtractYamlFileToJson(path)
	if err != nil {
		return err
	}

	ctx.CliClient.CreateEvaluation()

	//localConfigContent, err := ctx.LocalConfig.GetLocalConfiguration()
	//if err != nil {
	//	return err
	//}
	//
	//var cliId = localConfigContent.CliId
	//
	//validYamlFilesConfigurationsChan, invalidYamlFilesChan := files.ExtractFilesConfigurations(args, 100)
	//
	//var validFiles []*extractor.FileConfigurations
	//for validFile := range validYamlFilesConfigurationsChan {
	//	validFiles = append(validFiles, validFile)
	//}
	//
	//var invalidFiles []*validation.InvalidYamlFile
	//for invalidFile := range invalidYamlFilesChan {
	//	invalidFiles = append(invalidFiles, invalidFile)
	//}
	//
	////createEvaluationResponse, err := ctx.Evaluator.CreateEvaluation(localConfigContent.CliId, ctx.CliVersion, flags.K8sVersion, flags.PolicyName)
	////if err != nil {
	////	return err
	////}
	//
	//if len(invalidFiles) == 1 {
	//	return invalidFiles[0].ValidationErrors[0]
	//}
	//
	//if len(validFiles) == 1 {
	//	var file = validFiles[0]
	//	ctx.Printer.PrintMessage(file.FileName, "error")
	//	ctx.Printer.PrintMessage(file.Configurations[0])
	//}
	//
	//err = ctx.Evaluator.UpdateFailedYamlValidation(invalidYamlFiles, createEvaluationResponse.EvaluationId, stopEvaluation)
	//if err != nil {
	//	return err
	//}
	//
	//return nil
}
