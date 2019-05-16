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
	"fmt"

	"github.com/projectriff/riff/pkg/cli"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/resource"
	buildv1alpha1 "github.com/projectriff/system/pkg/apis/build/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ApplicationCreateOptions struct {
	cli.ResourceOptions

	Image     string
	CacheSize string

	LocalPath   string
	GitRepo     string
	GitRevision string
	SubPath     string
}

func (opts *ApplicationCreateOptions) Validate(ctx context.Context) *cli.FieldError {
	errs := &cli.FieldError{}

	errs = errs.Also(opts.ResourceOptions.Validate(ctx))

	if opts.Image == "" {
		errs = errs.Also(cli.ErrMissingField("image"))
	} else if false {
		// TODO validate image
	}

	if opts.CacheSize != "" {
		// must parse as a resource quantity
		if _, err := resource.ParseQuantity(opts.CacheSize); err != nil {
			errs = errs.Also(cli.ErrInvalidValue(opts.CacheSize, "cache-size"))
		}
	}

	// git-repo and local-path are mutually exclusive
	if opts.GitRepo == "" && opts.LocalPath == "" {
		errs = errs.Also(cli.ErrMissingOneOf("git-repo", "local-path"))
	} else if opts.GitRepo != "" && opts.LocalPath != "" {
		errs = errs.Also(cli.ErrMultipleOneOf("git-repo", "local-path"))
	}

	// git-revision is required for git-repo
	if opts.GitRepo != "" && opts.GitRevision == "" {
		errs = errs.Also(cli.ErrMissingField("git-revision"))
	}

	if opts.LocalPath != "" {
		if opts.SubPath != "" {
			// sub-path cannot be used with local-path
			errs = errs.Also(cli.ErrDisallowedFields("sub-path"))
		}
		if opts.CacheSize != "" {
			// cache-size cannot be used with local-path
			errs = errs.Also(cli.ErrDisallowedFields("cache-size"))
		}
	}

	return errs
}

func NewApplicationCreateCommand(c *cli.Config) *cobra.Command {
	opts := &ApplicationCreateOptions{}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "<todo>",
		Example: "<todo>",
		Args: cli.Args(
			cli.NameArg(&opts.Name),
		),
		PreRunE: cli.ValidateOptions(opts),
		RunE: func(cmd *cobra.Command, args []string) error {
			application := &buildv1alpha1.Application{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: opts.Namespace,
					Name:      opts.Name,
				},
				Spec: buildv1alpha1.ApplicationSpec{
					Image: opts.Image,
				},
			}
			if opts.CacheSize != "" {
				quantity := resource.MustParse(opts.CacheSize)
				application.Spec.CacheSize = &quantity
			}
			if opts.GitRepo != "" {
				application.Spec.Source = buildv1alpha1.Source{
					Git: &buildv1alpha1.GitSource{
						URL:      opts.GitRepo,
						Revision: opts.GitRevision,
					},
					SubPath: opts.SubPath,
				}
			}
			if opts.LocalPath != "" {
				// TODO implement
				return fmt.Errorf("not implemented")
			}

			application, err := c.Build().Applications(opts.Namespace).Create(application)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Created application %q\n", application.Name)
			return nil
		},
	}

	cli.NamespaceFlag(cmd, c, &opts.Namespace)
	cmd.Flags().StringVar(&opts.Image, "image", "", "<todo>")
	cmd.Flags().StringVar(&opts.CacheSize, "cache-size", "", "<todo>")
	cmd.Flags().StringVar(&opts.LocalPath, "local-path", "", "<todo>")
	cmd.Flags().StringVar(&opts.GitRepo, "git-repo", "", "<todo>")
	cmd.Flags().StringVar(&opts.GitRevision, "git-revision", "master", "<todo>")
	cmd.Flags().StringVar(&opts.SubPath, "sub-path", "", "<todo>")

	return cmd
}
