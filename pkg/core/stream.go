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

package core

import (
	"io"

	streamsv1alpha1 "github.com/projectriff/stream-controller/pkg/apis/streamcontroller/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreateStreamOptions struct {
	Namespace  string
	Name       string
	Definition string
}

func (c *client) CreateStream(options CreateStreamOptions, log io.Writer) error {
	ns := c.explicitOrConfigNamespace(options.Namespace)

	stream := &streamsv1alpha1.Stream{
		ObjectMeta: v1.ObjectMeta{
			Name:      options.Name,
			Namespace: ns,
		},
		Spec: streamsv1alpha1.StreamSpec{
			Definition: options.Definition,
		},
	}

	_, err := c.streams.ProjectriffV1alpha1().Streams(ns).Create(stream)
	if err != nil {
		return err
	}

	return nil
}
