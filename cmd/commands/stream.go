/*
 * Copyright 2018 The original author or authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package commands

import (
	"github.com/projectriff/riff/pkg/core"
	"github.com/spf13/cobra"
)

const (
	streamCreateStreamNameIndex = iota
	streamCreateNumberOfArgs
)

func Stream() *cobra.Command {
	return &cobra.Command{
		Use:   "stream",
		Short: "Interact with stream related resources",
	}
}

func StreamCreate(fcTool *core.Client) *cobra.Command {
	createStreamOptions := core.CreateStreamOptions{}

	command := &cobra.Command{
		Use:     "create",
		Short:   "Create a new stream resource",
		Long:    "Create a new stream resource from the provided definition\n",
		Example: `  riff stream create mystream --definition "foo | bar"`,
		Args: ArgValidationConjunction(
			cobra.ExactArgs(streamCreateNumberOfArgs),
			AtPosition(streamCreateStreamNameIndex, ValidName()),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			streamName := args[streamCreateStreamNameIndex]
			createStreamOptions.Name = streamName
			err := (*fcTool).CreateStream(createStreamOptions, cmd.OutOrStdout())
			if err != nil {
				return err
			}
			PrintSuccessfulCompletion(cmd)
			return nil
		},
	}

	LabelArgs(command, "STREAM_NAME")

	command.Flags().StringVarP(&createStreamOptions.Namespace, "namespace", "n", "", "the `namespace` of the service")
	command.Flags().StringVar(&createStreamOptions.Definition, "definition", "", "the definition of the stream")
	command.MarkFlagRequired("definition")

	return command
}
