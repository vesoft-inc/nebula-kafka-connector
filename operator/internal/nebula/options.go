/*
Copyright 2023 Vesoft Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package nebula

import "time"

type Option func(ops *Options)

type Options struct {
	Username string
	Password string
	Timeout  time.Duration
}

func loadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	return opts
}

func SetOptions(options Options) Option {
	return func(opts *Options) {
		*opts = options
	}
}

func SetTimeout(duration time.Duration) Option {
	return func(options *Options) {
		options.Timeout = duration
	}
}

func SetAuthInfo(username, password string) Option {
	return func(options *Options) {
		options.Username = username
		options.Password = password
	}
}
