/*
Copyright 2025.

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

package main

import (
	"context"
	"net/http"

	"github.com/go-logr/logr"
	polyfeav1alpha1 "github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/internal/controller"
	"github.com/polyfea/polyfea-controller/internal/repository"
	webserver "github.com/polyfea/polyfea-controller/internal/web-server"
	"github.com/polyfea/polyfea-controller/internal/web-server/configuration"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func startManager(cancel context.CancelFunc, mgr manager.Manager) {
	defer wg.Done()
	defer cancel()

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
	}
}

func startHTTPServer(
	ctx context.Context,
	cancel context.CancelFunc,
	microFrontendClassRepository repository.Repository[*polyfeav1alpha1.MicroFrontendClass],
	microFrontendRepository repository.Repository[*polyfeav1alpha1.MicroFrontend],
	webComponentRepository repository.Repository[*polyfeav1alpha1.WebComponent],
	logger *logr.Logger) {

	defer wg.Done()
	defer cancel()

	server := &http.Server{
		Addr: ":" + configuration.GetConfigurationValueOrDefault("POLYFEA_WEB_SERVER_PORT", "8082"),
		Handler: otelhttp.NewHandler(
			webserver.SetupRouter(microFrontendClassRepository, microFrontendRepository, webComponentRepository, logger),
			"polyfea-web-server",
		),
	}

	go func() {
		<-ctx.Done()
		err := server.Shutdown(context.Background())
		if err != nil {
			setupLog.Error(err, "problem shutting down server")
		}
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		setupLog.Error(err, "problem running server")
	}
}

// registerControllers creates repositories and registers all controllers with the manager.
func registerControllers(mgr manager.Manager) (
	repository.Repository[*polyfeav1alpha1.MicroFrontendClass],
	repository.Repository[*polyfeav1alpha1.MicroFrontend],
	repository.Repository[*polyfeav1alpha1.WebComponent],
) {
	microFrontendRepository := repository.NewInMemoryRepository[*polyfeav1alpha1.MicroFrontend]()
	microFrontendClassRepository := repository.NewInMemoryRepository[*polyfeav1alpha1.MicroFrontendClass]()
	webComponentRepository := repository.NewInMemoryRepository[*polyfeav1alpha1.WebComponent]()

	if err := (&controller.MicroFrontendReconciler{
		Client:     mgr.GetClient(),
		Scheme:     mgr.GetScheme(),
		Repository: microFrontendRepository,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "MicroFrontend")
		exitWithError()
	}
	if err := (&controller.WebComponentReconciler{
		Client:     mgr.GetClient(),
		Scheme:     mgr.GetScheme(),
		Repository: webComponentRepository,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "WebComponent")
		exitWithError()
	}
	if err := (&controller.MicroFrontendClassReconciler{
		Client:     mgr.GetClient(),
		Scheme:     mgr.GetScheme(),
		Repository: microFrontendClassRepository,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "MicroFrontendClass")
		exitWithError()
	}

	return microFrontendClassRepository, microFrontendRepository, webComponentRepository
}
