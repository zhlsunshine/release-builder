// Copyright Istio Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/howardjohn/istio-release/pkg/model"
)

func InputManifestToManifest(in model.InputManifest) (model.Manifest, error) {
	wd := in.Directory
	if wd == "" {
		var err error
		wd, err = ioutil.TempDir(os.TempDir(), "istio-release")
		if err != nil {
			return model.Manifest{}, fmt.Errorf("failed to create working directory: %v", err)
		}
	}
	outputs := map[model.BuildOutput]struct{}{}
	for _, o := range in.BuildOutputs {
		switch strings.ToLower(o) {
		case "docker":
			outputs[model.Docker] = struct{}{}
		case "helm":
			outputs[model.Helm] = struct{}{}
		case "debian":
			outputs[model.Debian] = struct{}{}
		case "archive":
			outputs[model.Archive] = struct{}{}
		default:
			return model.Manifest{}, fmt.Errorf("unknown build output: %v", o)
		}
	}
	if len(outputs) == 0 {
		outputs[model.Docker] = struct{}{}
		outputs[model.Helm] = struct{}{}
		outputs[model.Debian] = struct{}{}
		outputs[model.Archive] = struct{}{}
	}
	return model.Manifest{
		Dependencies: in.Dependencies,
		Version:      in.Version,
		Docker:       in.Docker,
		Directory:    wd,
		BuildOutputs: outputs,
	}, nil
}

func ReadManifest(manifestFile string) (model.Manifest, error) {
	manifest := model.Manifest{}
	by, err := ioutil.ReadFile(manifestFile)
	if err != nil {
		return manifest, fmt.Errorf("failed to read manifest file: %v", err)
	}
	if err := yaml.Unmarshal(by, &manifest); err != nil {
		return manifest, fmt.Errorf("failed to unmarshal manifest file: %v", err)
	}
	return manifest, nil
}

func ReadInManifest(manifestFile string) (model.InputManifest, error) {
	manifest := model.InputManifest{}
	by, err := ioutil.ReadFile(manifestFile)
	if err != nil {
		return manifest, fmt.Errorf("failed to read manifest file: %v", err)
	}
	if err := yaml.Unmarshal(by, &manifest); err != nil {
		return manifest, fmt.Errorf("failed to unmarshal manifest file: %v", err)
	}
	return manifest, nil
}