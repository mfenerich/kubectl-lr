/*
Package limitrange provides a CLI plugin for Kubernetes that extends the
functionality of `kubectl` with custom commands for managing and configuring
limitrange resources.

The main feature of this plugin is the `kubectl limitrange` command, which allows users
to create and validate limitrange resources with customizable specifications,
including CPU and memory limits, default requests, and dry-run options.

Installation:
This plugin is designed to be installed via `krew`, the Kubernetes plugin manager.
Ensure that `krew` is installed and configured before proceeding.

Usage:
To create a LimitRange resource with specific limits, use:

	kubectl limitrange my-limitrange --namespace=my-namespace --max-cpu="1" --min-cpu=100m --max-memory=500Mi --dry-run=client -o yaml

License:
This project is licensed under the Apache License, Version 2.0. You may not use
this file except in compliance with the License. You may obtain a copy of the
License at:

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.
*/
package limitrange // import "k8s.io/limitrange"
