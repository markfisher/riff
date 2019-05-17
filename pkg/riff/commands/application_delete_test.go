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

package commands_test

import (
	"github.com/projectriff/riff/pkg/cli"
	"github.com/projectriff/riff/pkg/riff/commands"
	"github.com/projectriff/riff/pkg/testing"
	buildv1alpha1 "github.com/projectriff/system/pkg/apis/build/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestApplicationDeleteOptions(t *testing.T) {
	table := testing.OptionsTable{
		{
			Name: "invalid delete",
			Options: &commands.ApplicationDeleteOptions{
				DeleteOptions: testing.InvalidDeleteOptions,
			},
			ExpectFieldError: testing.InvalidDeleteOptionsFieldError,
		},
		{
			Name: "valid delete",
			Options: &commands.ApplicationDeleteOptions{
				DeleteOptions: testing.ValidDeleteOptions,
			},
			ShouldValidate: true,
		},
	}

	table.Run(t)
}

func TestApplicationDeleteCommand(t *testing.T) {
	applicationName := "test-application"
	applicationOtherName := "test-other-application"
	defaultNamespace := "default"

	table := testing.CommandTable{
		{
			Name:        "invalid args",
			Args:        []string{},
			ShouldError: true,
		},
		{
			Name: "delete all applications",
			Args: []string{cli.AllFlagName},
			GivenObjects: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Name:      applicationName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectDeleteCollections: []testing.DeleteCollectionRef{{
				Group:     "build.projectriff.io",
				Resource:  "applications",
				Namespace: defaultNamespace,
			}},
		},
		{
			Name: "delete application",
			Args: []string{applicationName},
			GivenObjects: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Name:      applicationName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectDeletes: []testing.DeleteRef{{
				Group:     "build.projectriff.io",
				Resource:  "applications",
				Namespace: defaultNamespace,
				Name:      applicationName,
			}},
		},
		{
			Name: "delete applications",
			Args: []string{applicationName, applicationOtherName},
			GivenObjects: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Name:      applicationName,
						Namespace: defaultNamespace,
					},
				},
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Name:      applicationOtherName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectDeletes: []testing.DeleteRef{{
				Group:     "build.projectriff.io",
				Resource:  "applications",
				Namespace: defaultNamespace,
				Name:      applicationName,
			}, {
				Group:     "build.projectriff.io",
				Resource:  "applications",
				Namespace: defaultNamespace,
				Name:      applicationOtherName,
			}},
		},
		{
			Name: "application does not exist",
			Args: []string{applicationName},
			ExpectDeletes: []testing.DeleteRef{{
				Group:     "build.projectriff.io",
				Resource:  "applications",
				Namespace: defaultNamespace,
				Name:      applicationName,
			}},
			ShouldError: true,
		},
		{
			Name: "delete error",
			Args: []string{applicationName},
			GivenObjects: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Name:      applicationName,
						Namespace: defaultNamespace,
					},
				},
			},
			WithReactors: []testing.ReactionFunc{
				testing.InduceFailure("delete", "applications"),
			},
			ExpectDeletes: []testing.DeleteRef{{
				Group:     "build.projectriff.io",
				Resource:  "applications",
				Namespace: defaultNamespace,
				Name:      applicationName,
			}},
			ShouldError: true,
		},
	}

	table.Run(t, commands.NewApplicationDeleteCommand)
}
