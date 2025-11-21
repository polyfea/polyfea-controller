# polyfea-controller

This project is a Kubernetes controller built with the Operator SDK. It is a core component of the polyfea platform and enables clean, decoupled development of microfrontends.

## Description

The controller introduces three Custom Resource Definitions (CRDs) into Kubernetes:

**MicroFrontendClass** – Defines shared configuration such as routing base URIs, CSP headers, meta tags, headers, and optional PWA settings.
**MicroFrontend** – Describes an individual microfrontend and its deployment-specific configuration.
**WebComponent** – Represents the web components that compose a microfrontend.

Together, these resources provide a structured, Kubernetes-native way to define, manage, and integrate microfrontends within a cluster.

A key feature of the controller is its ability to dynamically serve metadata required by frontend applications. This allows microfrontends to be developed independently while still integrating seamlessly at runtime.

By enabling strong separation of concerns, the polyfea-controller improves flexibility, scalability, and maintainability when building complex microfrontend architectures on Kubernetes.

### Microfrontend Class

This Custom Resource Definition (CRD) enables developers to define and manage **MicroFrontendClass** resources within a Kubernetes cluster. A MicroFrontendClass describes how a microfrontend should behave when served through the platform.

Each MicroFrontendClass specifies a **baseUri**, which determines the URL prefix it applies to. All incoming requests whose paths begin with this baseUri will use the configuration of the corresponding MicroFrontendClass. This includes injecting additional response headers, applying the appropriate Content Security Policy (CSP), and rendering configured meta tags.

#### Specification Overview

A MicroFrontendClass supports the following configuration properties:

- **baseUri**: The base URI associated with the frontend class. All requests starting with this URI will use this configuration.
- **cspHeader**: The Content Security Policy header to apply. If omitted, a secure default value is used.
- **extraHeaders**: Additional HTTP headers to include in responses. If omitted, no extra headers are added.
- **extraMetaTags**: Additional meta tags rendered in the resulting HTML. If omitted, none are added.
- **rolesHeader**: The header containing user roles. Defaults to x-auth-request-roles.
- **title**: The title associated with the frontend class.
- **userHeader**: The header containing the user identifier. Defaults to x-auth-request-user.
- **progressiveWebApp**: Configures optional Progressive Web App (PWA) capabilities. This includes defining the web app manifest, setting up pre-cache and runtime caching strategies, and specifying how frequently the service worker should reconcile changes. If omitted, PWA support is disabled.

#### Required Fields

The baseUri and title fields are mandatory for every MicroFrontendClass.

### MicroFrontend

The **MicroFrontend** Custom Resource Definition (CRD) allows developers to define and manage individual microfrontends within a Kubernetes cluster. It specifies where a microfrontend is hosted, how its assets should be loaded, and which other microfrontends it depends on.

A MicroFrontend defines the origin service hosting its module and static assets, the loading behavior (including proxying), and any required dependencies. When proxying is enabled, the controller transparently fetches and serves the microfrontend’s modules and static resources, simplifying integration and improving isolation.

#### Specification Overview

A MicroFrontend resource includes the following properties:

* **cacheControl**: The cache-control header applied to the microfrontend when the caching strategy is set to `cache`.
* **cacheStrategy**: The caching strategy to use for this microfrontend. Defaults to `none`.
* **dependsOn**: A list of other microfrontends that must be loaded before this one.
* **frontendClass**: The name of the MicroFrontendClass that defines shared configuration (headers, CSP, PWA settings, etc.). Defaults to `polyfea-controller-default`.
* **modulePath**: The relative path to the microfrontend’s module file within the referenced service.
* **proxy**: Determines whether the controller proxies the loading of modules and web component resources. Defaults to `true`.
* **service**: A reference to the service hosting the microfrontend's module and CSS files. This can be either:
  * An **in-cluster service** specified using `name`, with optional `namespace` (defaults to the MicroFrontend's namespace), `port` (defaults to 80), `scheme` (defaults to http), and `domain` (defaults to svc.cluster.local).
  * An **external service** specified using `uri` with the complete URL (e.g., `https://cdn.example.com`).
  * **Note**: Either `name` or `uri` must be specified, but not both.
* **staticPaths**: A list of static resources (scripts, stylesheets, or other link-based assets) that must be loaded before the microfrontend. Each entry may specify attributes, load-waiting behavior, and whether the resource should be proxied.
* **cacheOptions**: Optional PWA-style cache configuration (pre-cache and runtime caching rules).

#### Required Fields

The following fields are mandatory:

* **service**
* **modulePath**
* **frontendClass**

These values ensure the controller knows where the microfrontend lives, which class configuration it belongs to, and where to load its primary module.

### WebComponent

The **WebComponent** Custom Resource Definition (CRD) allows developers to define and manage individual web components within a Kubernetes-managed microfrontend environment. It specifies where a web component originates, how it should be configured, and under which conditions it should be displayed.

A WebComponent may reference the **MicroFrontend** that provides its implementation. This field is optional—if omitted, the controller assumes the component is already available (for example, when it is a native browser element or when it was loaded by another microfrontend).

Display behavior is controlled through **displayRules**, which define when the component should be rendered. Each rule group in the top-level list represents an OR condition: if *any* DisplayRule matches, the component is included. Within a single DisplayRule, the sets `allOf`, `anyOf`, and `noneOf` are combined using AND semantics:

* all matchers in `allOf` must match,
* at least one matcher in `anyOf` must match (if provided),
* none of the matchers in `noneOf` may match.

Only when all conditions within a DisplayRule evaluate to true is that rule considered matched.

### Specification Overview

A WebComponent resource supports the following properties:

* **attributes**: A list of attribute key/value pairs applied to the final HTML element. The value may contain any valid JSON type.
* **displayRules**: Conditions that determine when the component should be loaded. The list is evaluated using OR semantics.
* **element**: The HTML tag name for the component (e.g., `my-menu-item`).
* **microFrontend**: (Optional) The MicroFrontend providing this component. If omitted, the controller assumes the component is already loaded or native.
* **priority**: Controls ordering when multiple components match. Higher values indicate higher priority. Defaults to `0`.
* **style**: Inline style definitions applied to the component.

### Required Fields

The following fields must be provided:

* **element**
* **displayRules**

## Getting Started

To run the controller, you need access to a Kubernetes cluster. For local development, you can use [Kind](https://sigs.k8s.io/kind) to spin up a lightweight cluster.
The controller uses the **current kubeconfig context**, so it will operate on whichever cluster is returned by:

```sh
kubectl cluster-info
```

---

### Configuring the Controller

The controller can be configured through the following environment variable:

* **`POLYFEA_WEB_SERVER_PORT`** – The port used by the controller’s internal web server. Defaults to `8082`.

---

### Running on the Cluster

#### 1. Install Custom Resource Instances

```sh
kubectl apply -f config/samples/
```

#### 2. Build and Push the Controller Image

```sh
make docker-build docker-push IMG=<registry>/polyfea-controller:<tag>
```

#### 3. Deploy the Controller

```sh
make deploy IMG=<registry>/polyfea-controller:<tag>
```

---

### Uninstalling

#### Remove CRDs

```sh
make uninstall
```

#### Undeploy the Controller

```sh
make undeploy
```

---

## Contributing

### Prerequisites

You’ll need the following tools:

* WSL, Linux, or macOS environment
* [Go](https://golang.org/doc/install) **1.24+**
* [Docker](https://docs.docker.com/get-docker/)
* [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* [Make](https://www.gnu.org/software/make/)

---

### How It Works

This project follows the Kubernetes [Operator Pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).
It uses Kubernetes **controllers**, each implementing a reconcile loop that continually drives the cluster state toward the desired state declared in the custom resources.

---

### Running Locally

#### 1. Install CRDs

```sh
make install
```

#### 2. Run the Controller (Foreground)

```sh
make run
```

> 💡 You can also run both steps in one go:
>
> ```sh
> make install run
> ```

---

### Modifying the API Definitions

When you update API types, regenerate manifests (CRDs, RBAC, samples) using:

```sh
make manifests
```

For a full list of available build and development commands:

```sh
make --help
```

For more detailed guidance, refer to the official
**[Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)**.

---

## License

Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.

You may obtain a copy of the License at:

```
http://www.apache.org/licenses/LICENSE-2.0
```

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an **"AS IS" BASIS**,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.

See the License for the specific language governing permissions and
limitations under the License.
