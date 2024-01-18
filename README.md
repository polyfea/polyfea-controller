# polyfea-controller

This project is a Kubernetes controller developed using the Operator SDK. It's a key component of a larger polyfea project, designed to facilitate the decoupled development of microfrontends.

## Description

The controller introduces three custom resources into the Kubernetes ecosystem: Microfrontend Class, Microfrontend, and WebComponent. These resources provide a structured way to manage and interact with microfrontends within a Kubernetes environment.

The Microfrontend Class and Microfrontend resources are used to define and control the behavior of individual microfrontends. The WebComponent resource is used to manage the web components that make up the microfrontends.

One of the main features of this project is its ability to serve information to a frontend for display. This allows developers to focus on building their microfrontends without worrying about how they will be presented to the end user.

By enabling decoupled development, this project allows for greater flexibility and scalability in building and deploying microfrontends. It's an innovative solution for managing complex frontend architectures in a Kubernetes environment.

### Microfrontend Class

This CRD allows developers to create and manage MicroFrontendClass resources in a Kubernetes cluster, providing a way to define and control the behavior of microfrontend classes.

MicroFrontendClass specifies baseUri for which it will be used. Every request on the server starting with its baseUri will consider the configuration of the MicroFrontendClass. It will serve the requests and respond with the extra headers for given microfrontend class. It will also take care of the Content Security Policy (CSP) header and the extra meta tags.

The MicroFrontendClass has several properties defined in its specification:
- baseUri: The base URI for which the frontend class will be used.
- cspHeader: The Content Security Policy (CSP) header for the frontend class. If not set, a default value is used.
- extraHeaders: An array of additional headers for the frontend class. If not set, no extra headers are used.
- extraMetaTags: An array of additional meta tags for the frontend class. If not set, no extra meta tags are used.
- rolesHeader: The name of the header that contains the roles of the user. Defaults to 'x-auth-request-roles'.
- title: The title that will be used for the frontend class.
- userHeader: The name of the header that contains the user id. Defaults to 'x-auth-request-user'.

The baseUri and title properties are required for the MicroFrontendClass.

### Microfrontend

This CRD allows developers to create and manage MicroFrontend resources in a Kubernetes cluster, providing a way to define and control the behavior of micro frontends.

The MicroFrontend specifies where the micro frontend is located and how it should be loaded. It also specifies the dependencies that should be loaded before the micro frontend.

The proxy specifies whether the loading of web components should be proxied by the controller. If the proxy is set to true, the controller will proxy the loading of web components from the service. If the proxy is set to false, the controller will not proxy the loading of web components from the service.

The MicroFrontend has several properties defined in its specification:

- cacheControl: Defines the cache control header for the micro frontend. This is only used if the caching strategy is set to 'cache'.
- cacheStrategy: Defines the caching strategy for the micro frontend. The default value is 'none'.
- dependsOn: List of dependencies that should be loaded before this micro frontend.
- frontendClass: The name of the frontend class that should be used for this micro frontend. The default value is 'polyfea-controller-default'.
- modulePath: Relative path to the module file within the service.
- proxy: Specifies whether the loading of web components should be proxied by the controller. The default value is true.
- service: Reference to a service from which the modules or CSS would be served. The fully qualified name of the service should be specified in the format <schema>://<service-name>.<namespace>.<cluster>.
- staticPaths: Relative path to the static files within the service. This includes static resources that should be loaded before this micro frontend.

The frontendClass, modulePath, and service properties are required for the MicroFrontend.

### WebComponent

This CRD allows developers to create and manage WebComponent resources in a Kubernetes cluster, providing a way to define and control the behavior of web components.

The WebComponent specifies the micro frontend from which the web component should be loaded. It also specifies the conditions under which the web component should be loaded. It may also specify the styles that should be applied to the web component.
The microFrontend might be ommited if the web component is loaded from a different micro frontend than the one that is specified in the MicroFrontend resource. The assumption is that the webcomponent was already loaded e.g. native html element.

The displayRules specifies the conditions under which the web component should be loaded. There is an OR operation between the elements of the DisplayRules list. If any of the DisplayRules is matched, the web component will be loaded. And an and operation between the elements of the DisplayRule. If all of the elements of the DisplayRule are matched, the DisplayRule will be matched.

The WebComponent has several properties defined in its specification:

- attributes: A list of key-value pairs that allows you to assign specific attributes to the element.
- displayRules: Defines the conditions under which the web component should be loaded. There is an OR operation between the elements of the DisplayRules list. If any of the DisplayRules is matched, the web component will be loaded.
- element: The HTML element tag name to be used when the matcher is matched.
- microFrontend: Reference to a microfrontend from which the webcomponent would be served.
- priority: Defines the priority of the webcomponent. Used for ordering the webcomponent within the shell. The higher the number, the higher the priority. The default priority is 0.
- style: Defines the styles that should be applied to the webcomponent.

The displayRules and element properties are required for the WebComponent.

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Configuring the controller

The controller can be configured using the following environment variables:

- `POLYFEA_WEB_SERVER_PORT`: The port on which the web server will listen. Defaults to 8082.

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/polyfea-controller:tag
```

3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/polyfea-controller:tag
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller from the cluster:

```sh
make undeploy
```

## Contributing

### Prerequisites
- WSL, Linux, or macOS environment
- [Go](https://golang.org/doc/install) version v1.21+
- [Docker/Kubernetes Cluster](https://docs.docker.com/get-docker/)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Make](https://www.gnu.org/software/make/)

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

