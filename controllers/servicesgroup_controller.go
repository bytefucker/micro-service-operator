/*
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
*/

package controllers

import (
	"context"

	k8sv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	appsv1 "github.com/bytefucker/micro-service-operator/api/v1"
)

// ServicesGroupReconciler reconciles a ServicesGroup object
type ServicesGroupReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.github.com,resources=servicesgroups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.github.com,resources=servicesgroups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.github.com,resources=servicesgroups/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ServicesGroup object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *ServicesGroupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	group := &appsv1.ServicesGroup{}
	if err := r.Get(ctx, req.NamespacedName, group); err != nil {
		// 如果没有实例，就返回空结果，这样外部就不再立即调用Reconcile方法了
		if errors.IsNotFound(err) {
			log.Info("instance not found, maybe removed")
			return reconcile.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	//开始创建服务
	for _, s := range group.Spec.Services {
		deployment := &k8sv1.Deployment{}
		if err := r.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: s.Name}, deployment); err != nil {
			log.Info("namespace=%s,name=%s deployment不存在，开始创建", req.Namespace, s.Name)
			createDeployment(ctx, r, group, &s)
		} else {
			log.Info("namespace=%s,name=%s deployment已存在，开始更新", req.Namespace, s.Name)
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServicesGroupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&appsv1.ServicesGroup{}).Complete(r)
}

// 创建一个Deployment
func createDeployment(ctx context.Context, r *ServicesGroupReconciler, sg *appsv1.ServicesGroup, service *appsv1.Service) error {
	deployment := &k8sv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: sg.Namespace,
			Name:      service.Name,
		},
		Spec: k8sv1.DeploymentSpec{
			Replicas: service.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": service.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": service.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            service.Name,
							Image:           service.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolSCTP,
									ContainerPort: service.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
	}
	if err := controllerutil.SetControllerReference(sg, deployment, r.Scheme); err != nil {
		log.Log.Error(err, "SetControllerReference error")
		return err
	}
	if err := r.Create(ctx, deployment); err != nil {
		log.Log.Error(err, "createDeployment error")
		return err
	}
	return nil
}
