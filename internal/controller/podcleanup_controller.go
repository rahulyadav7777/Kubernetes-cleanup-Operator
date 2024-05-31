package controller

import (
	"context"
	"strings"
	"time"

	batchv1 "github.com/example/pod-cleanup-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// PodCleanupReconciler reconciles a PodCleanup object
type PodCleanupReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=batch.example.com,resources=podcleanups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch.example.com,resources=podcleanups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=batch.example.com,resources=podcleanups/finalizers,verbs=update

//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;delete

func (r *PodCleanupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the PodCleanup instance
	podCleanup := &batchv1.PodCleanup{}
	err := r.Get(ctx, req.NamespacedName, podCleanup)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Define the cleanup logic for different pod states
	var podList corev1.PodList
	if err := r.List(ctx, &podList, client.InNamespace(req.Namespace)); err != nil {
		log.Error(err, "unable to list pods")
		return ctrl.Result{}, err
	}

	for _, pod := range podList.Items {
		if shouldCleanup(pod) {
			log.Info("Deleting pod", "pod", pod.Name)
			if err := r.Delete(ctx, &pod); client.IgnoreNotFound(err) != nil {
				log.Error(err, "failed to delete pod", "pod", pod.Name)
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{RequeueAfter: wait.Jitter(10*time.Minute, 0.1)}, nil
}

func shouldCleanup(pod corev1.Pod) bool {
	return isEvicted(pod) || isCrashLoopBackOff(pod) || isImagePullError(pod) || isFailed(pod)
}

func isEvicted(pod corev1.Pod) bool {
	return pod.Status.Phase == corev1.PodFailed && strings.Contains(pod.Status.Reason, "Evicted")
}

func isCrashLoopBackOff(pod corev1.Pod) bool {
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.State.Waiting != nil && containerStatus.State.Waiting.Reason == "CrashLoopBackOff" {
			return true
		}
	}
	return false
}

func isImagePullError(pod corev1.Pod) bool {
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.State.Waiting != nil && (containerStatus.State.Waiting.Reason == "ImagePullBackOff" || containerStatus.State.Waiting.Reason == "ErrImagePull") {
			return true
		}
	}
	return false
}

func isFailed(pod corev1.Pod) bool {
	return pod.Status.Phase == corev1.PodFailed
}

func (r *PodCleanupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.PodCleanup{}).
		Complete(r)
}
