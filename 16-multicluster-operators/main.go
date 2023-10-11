package main

import (
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	_ "k8s.io/client-go/plugin/pkg/client/auth"

	v1 "k8s.io/api/core/v1"
	"log"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.GetConfigOrDie()

	mgr, err := manager.New(cfg, manager.Options{
		Scheme:                        nil,
		MapperProvider:                nil,
		SyncPeriod:                    nil,
		Logger:                        nil,
		LeaderElection:                false,
		LeaderElectionResourceLock:    "",
		LeaderElectionNamespace:       "",
		LeaderElectionID:              "",
		LeaderElectionConfig:          nil,
		LeaderElectionReleaseOnCancel: false,
		LeaseDuration:                 nil,
		RenewDeadline:                 nil,
		RetryPeriod:                   nil,
		Namespace:                     "",
		MetricsBindAddress:            ":1234",
		HealthProbeBindAddress:        "",
		ReadinessEndpointName:         "",
		LivenessEndpointName:          "",
		Port:                          0,
		Host:                          "",
		CertDir:                       "",
		NewCache:                      nil,
		ClientBuilder:                 nil,
		ClientDisableCacheFor:         nil,
		DryRunClient:                  false,
		EventBroadcaster:              nil,
		GracefulShutdownTimeout:       nil,
	})
	if err != nil {
		return err
	}

	rec := &myReconciler{
		cache:  mgr.GetCache(),
		writer: mgr.GetClient(),
	}

	ctl, err := controller.New("my-controller", mgr, controller.Options{
		MaxConcurrentReconciles: 0,
		Reconciler:              rec,
		RateLimiter:             nil,
		Log:                     nil,
		CacheSyncTimeout:        0,
	})
	if err != nil {
		return err
	}

	if err := ctl.Watch(
		&source.Kind{Type: &v1.Secret{}},
		&handler.EnqueueRequestForObject{},
	); err != nil {
		return err
	}

	log.Printf("manager starting")
	return mgr.Start(context.Background())
}

type myReconciler struct {
	cache  client.Reader
	writer client.Writer
}

func (m *myReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log.Printf("got request: %v", request)
	sec := &v1.Secret{}
	if err := m.cache.Get(ctx, client.ObjectKey{
		Namespace: request.Namespace,
		Name:      request.Name,
	}, sec); err != nil {
		if errors.IsNotFound(err) {
			log.Printf("secret deleted")

			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// update our secret list
	cm := &v1.ConfigMap{}
	if err := m.cache.Get(ctx, client.ObjectKey{
		Namespace: "default",
		Name:      "secret-list",
	}, sec); err != nil {
		if errors.IsNotFound(err) {
			cm.Name = "secret-list"
			cm.Namespace = "default"
			if err := m.writer.Create(ctx, cm); err != nil {
				return reconcile.Result{}, err
			}
		} else {
			return reconcile.Result{}, err
		}
	}

	if cm.Data == nil {
		cm.Data = map[string]string{}
	}
	cm.Data[request.Namespace+"/"+request.Name] = "-"

	if err := m.writer.Update(ctx, cm); err != nil {
		log.Printf("update error: %v", request)
		return reconcile.Result{}, err
	}

	log.Printf("updated secret map: %v", request)

	return reconcile.Result{}, nil
}
