/*
 * Copyright 2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package commands

import (
	"context"

	"github.com/projectriff/riff/pkg/cli"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ApplicationDeleteOptions struct {
	cli.DeleteOptions
}

func (opts *ApplicationDeleteOptions) Validate(ctx context.Context) *cli.FieldError {
	errs := &cli.FieldError{}

	errs = errs.Also(opts.DeleteOptions.Validate(ctx))

	return errs
}

func (opts *ApplicationDeleteOptions) Exec(ctx context.Context, c *cli.Config) error {
	client := c.Build().Applications(opts.Namespace)

	if opts.All {
		return client.DeleteCollection(nil, metav1.ListOptions{})
	}

	for _, name := range opts.Names {
		if err := client.Delete(name, nil); err != nil {
			return err
		}
	}

	return nil
}

func NewApplicationDeleteCommand(c *cli.Config) *cobra.Command {
	opts := &ApplicationDeleteOptions{}

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "<todo>",
		Example: "<todo>",
		Args: cli.Args(
			cli.NamesArg(&opts.Names),
		),
		PreRunE: cli.ValidateOptions(opts),
		RunE:    cli.ExecOptions(c, opts),
	}

	cli.NamespaceFlag(cmd, c, &opts.Namespace)
	cmd.Flags().BoolVar(&opts.All, cli.StripDash(cli.AllFlagName), false, "<todo>")
	
	return cmd
}
