// Copyright 2023 Nautes Authors
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

//Package controllers use to set k8s cluster info to secret store
package controllers

import (
	"context"
	"errors"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	operatorerror "github.com/nautes-labs/cluster-operator/pkg/error"
	factory "github.com/nautes-labs/cluster-operator/pkg/secretclient/factory"
	secretclient "github.com/nautes-labs/cluster-operator/pkg/secretclient/interface"
	clustercrd "github.com/nautes-labs/pkg/api/v1alpha1"
	nautesconfig "github.com/nautes-labs/pkg/pkg/nautesconfigs"
)

const (
	clusterFinalizerName   = "cluster.nautes.resource.nautes.io/finalizers"
	clusterConditionType   = "SecretStoreReady"
	clusterConditionReason = "RegularUpdate"
)

// ClusterReconciler reconciles a Cluster object
type ClusterReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	Configs nautesconfig.NautesConfigs
}

//+kubebuilder:rbac:groups=nautes.resource.nautes.io,resources=clusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=nautes.resource.nautes.io,resources=clusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=nautes.resource.nautes.io,resources=clusters/finalizers,verbs=update
//+kubebuilder:rbac:groups=nautes.resource.nautes.io,resources=coderepoes,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *ClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	cluster := &clustercrd.Cluster{}
	if err := r.Get(ctx, req.NamespacedName, cluster); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !cluster.DeletionTimestamp.IsZero() {
		if err := r.deleteCluster(ctx, *cluster); err != nil {
			setStatus(cluster, nil, err)
			if err := r.Status().Update(ctx, cluster); err != nil {
				logger.Error(err, "update condition failed")
			}
			return ctrl.Result{}, err
		}

		controllerutil.RemoveFinalizer(cluster, clusterFinalizerName)
		if err := r.Update(ctx, cluster); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if err := isLegal(cluster); err != nil {
		setStatus(cluster, nil, err)
		if err := r.Status().Update(ctx, cluster); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}

	controllerutil.AddFinalizer(cluster, clusterFinalizerName)
	if err := r.Update(ctx, cluster); err != nil {
		return ctrl.Result{}, err
	}

	result, err := r.syncCluster(ctx, *cluster)
	if err := r.updateStatus(ctx, req.NamespacedName,
		func(c *clustercrd.Cluster) { setStatus(c, result, err) }); err != nil {
		return ctrl.Result{}, err
	}

	requeueAfter := time.Second * 60
	if err != nil {
		var operatorErr operatorerror.Errors
		if errors.As(err, &operatorErr) {
			requeueAfter = operatorErr.RequeueAfter
			logger.Error(err, "get cluster info failed")
		} else {
			return ctrl.Result{}, err
		}
	}

	logger.V(1).Info("reconcile finish", "RequeueAfter", fmt.Sprintf("%gs", requeueAfter.Seconds()))
	return ctrl.Result{RequeueAfter: requeueAfter}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&clustercrd.Cluster{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}

func (r *ClusterReconciler) newSecretClient(ctx context.Context) (secretclient.SecretClient, error) {
	cfg, err := r.Configs.GetConfigByClient(r.Client)
	if err != nil {
		return nil, fmt.Errorf("read error global config failed: %w", err)
	}

	secClient, err := factory.GetFactory().NewSecretClient(ctx, cfg, r.Client)
	if err != nil {
		return nil, fmt.Errorf("get secret client failed: %w", err)
	}

	return secClient, nil
}

func (r *ClusterReconciler) syncCluster(ctx context.Context, cluster clustercrd.Cluster) (*clustercrd.MgtClusterAuthStatus, error) {
	lastCluster, err := getSyncHistory(&cluster)
	if err != nil {
		return nil, err
	}

	secretClient, err := r.newSecretClient(ctx)
	if err != nil {
		return nil, err
	}

	result, err := secretClient.Sync(ctx, &cluster, lastCluster)
	if err != nil {
		return nil, fmt.Errorf("sync cluster info to secret store failed: %w", err)
	}

	return &clustercrd.MgtClusterAuthStatus{
		LastSuccessSpec: cluster.SpecToJsonString(),
		SecretID:        result.SecretID,
	}, nil
}

func (r *ClusterReconciler) deleteCluster(ctx context.Context, cluster clustercrd.Cluster) error {
	if !controllerutil.ContainsFinalizer(&cluster, clusterFinalizerName) {
		return nil
	}

	lastCluster, err := getSyncHistory(&cluster)
	if err != nil {
		return err
	}

	secretClient, err := r.newSecretClient(ctx)
	if err != nil {
		return err
	}

	err = secretClient.Delete(ctx, lastCluster)
	if err != nil {
		return fmt.Errorf("delete cluster info from secret store failed: %w", err)
	}
	return nil
}

type setStatusFunction func(*clustercrd.Cluster)

func (r *ClusterReconciler) updateStatus(ctx context.Context, key types.NamespacedName, setStatus setStatusFunction) error {
	cluster := &clustercrd.Cluster{}
	if err := r.Get(ctx, key, cluster); err != nil {
		return err
	}

	setStatus(cluster)

	return r.Status().Update(ctx, cluster)
}

func isLegal(cluster *clustercrd.Cluster) error {
	lastCluster, err := getSyncHistory(cluster)
	if err != nil {
		return err
	}

	return clustercrd.ValidateCluster(cluster, lastCluster, false)
}

func setStatus(cluster *clustercrd.Cluster, result *clustercrd.MgtClusterAuthStatus, err error) {
	if err != nil {
		condition := metav1.Condition{
			Type:    clusterConditionType,
			Status:  "False",
			Reason:  clusterConditionReason,
			Message: err.Error(),
		}
		cluster.Status.SetConditions([]metav1.Condition{condition}, map[string]bool{clusterConditionType: true})
	} else {
		condition := metav1.Condition{
			Type:   clusterConditionType,
			Status: "True",
			Reason: clusterConditionReason,
		}
		cluster.Status.SetConditions([]metav1.Condition{condition}, map[string]bool{clusterConditionType: true})

		if result != nil {
			cluster.Status.MgtAuthStatus = result
			conditions := cluster.Status.GetConditions(map[string]bool{clusterConditionType: true})
			if len(conditions) == 1 {
				cluster.Status.MgtAuthStatus.LastSuccessTime = conditions[0].LastTransitionTime
			}
		}
	}
}

func getSyncHistory(cluster *clustercrd.Cluster) (*clustercrd.Cluster, error) {
	if cluster.Status.MgtAuthStatus == nil {
		return nil, nil
	}

	lastCluster, err := clustercrd.GetClusterFromString(cluster.Name, cluster.Namespace, cluster.Status.MgtAuthStatus.LastSuccessSpec)
	if err != nil {
		return nil, fmt.Errorf("get history from status failed: %w", err)
	}

	return lastCluster, nil
}
