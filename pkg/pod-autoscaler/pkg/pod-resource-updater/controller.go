package resourceupdater

import (
	"context"
	"fmt"
	"log"

	"github.com/lterrac/system-autoscaler/pkg/informers"
	"github.com/lterrac/system-autoscaler/pkg/pod-autoscaler/pkg/logger"

	"time"

	"github.com/lterrac/system-autoscaler/pkg/apis/systemautoscaler/v1beta1"
	podscalesclientset "github.com/lterrac/system-autoscaler/pkg/generated/clientset/versioned"
	samplescheme "github.com/lterrac/system-autoscaler/pkg/generated/clientset/versioned/scheme"
	"github.com/lterrac/system-autoscaler/pkg/podscale-controller/pkg/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
)

const controllerAgentName = "pod-resource-updater"

// Controller is the controller that recommends resources to the pods.
// For each Pod Scale assigned to the recommender, it will have a pod saved in a list.
// Every x seconds, the recommender polls the metrics from the pod by using an http request.
// For pod metrics it retrieves, it computes the new resources to assign to the pod.
type Controller struct {

	// podScalesClientset is a clientset for our own API group
	podScalesClientset podscalesclientset.Interface

	// kubernetesCLientset is the client-go of kubernetes
	kubernetesClientset kubernetes.Interface

	listers informers.Listers

	podScalesSynced cache.InformerSynced
	podSynced       cache.InformerSynced

	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder

	log *logger.Logger

	// in is the input channel.
	in chan types.NodeScales
}

// NewController returns a new sample controller
func NewController(kubernetesClientset *kubernetes.Clientset,
	podScalesClientset podscalesclientset.Interface,
	informers informers.Informers,
	in chan types.NodeScales) *Controller {

	// Create event broadcaster
	// Add sample-controller types to the default Kubernetes Scheme so Events can be
	// logged for sample-controller types.
	utilruntime.Must(samplescheme.AddToScheme(scheme.Scheme))
	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	fileLogger, err := logger.NewFileLogger("/var/podscale.json")

	//TODO: remove as soon as possible
	if err != nil {
		log.Fatal("Error while setting up Logger")
	}

	// Instantiate the Controller
	controller := &Controller{
		podScalesClientset:  podScalesClientset,
		kubernetesClientset: kubernetesClientset,
		recorder:            recorder,
		listers:             informers.GetListers(),
		podScalesSynced:     informers.PodScale.Informer().HasSynced,
		podSynced:           informers.Pod.Informer().HasSynced,
		log:                 fileLogger,
		in:                  in,
	}

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting pod resource updater controller")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh,
		c.podScalesSynced,
		c.podSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting pod resource updater workers")
	// Launch the workers to process podScale resources and recommend new pod scales
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runNodeScaleWorker, time.Second, stopCh)
	}

	return nil
}

// Shutdown is called when the controller has finished its work
func (c *Controller) Shutdown() {
	utilruntime.HandleCrash()
}

func (c *Controller) runNodeScaleWorker() {
	for nodeScale := range c.in {
		klog.Info("Processing ", nodeScale)
		for _, podScale := range nodeScale.PodScales {

			pod, err := c.listers.Pods(podScale.Spec.Namespace).Get(podScale.Spec.Pod)
			if err != nil {
				klog.Error("Error retrieving the pod: ", err)
				return
			}

			newPod, err := syncPod(pod, *podScale)
			if err != nil {
				klog.Error("Error syncing the pod: ", err)
				return
			}

			// try both updates in dry-run first and then actuate them consistently
			updatedPod, updatedPodScale, err := c.AtomicResourceUpdate(newPod, podScale)

			if err != nil {
				klog.Error("Error while updating pod and podscale: ", err)
				//TODO: We are using this channel as a workqueue. Why don't use one?
				c.in <- nodeScale
				return
			}

			//TODO: handle error
			_ = c.log.Log(updatedPodScale)

			klog.Info("Desired resources:", updatedPodScale.Spec.DesiredResources)
			klog.Info("Capped resources:", updatedPodScale.Status.CappedResources)
			klog.Info("Actual resources:", updatedPodScale.Status.ActualResources)
			klog.Info("Pod resources:", updatedPod.Spec.Containers[0].Resources)
		}
	}
}

// AtomicResourceUpdate updates a Pod and its PodScale consistently in order to keep synchronized the two resources. Before performing the real update
// it runs a request in dry-run and it checks for any potential error
func (c *Controller) AtomicResourceUpdate(pod *corev1.Pod, podScale *v1beta1.PodScale) (*corev1.Pod, *v1beta1.PodScale, error) {
	var err error
	_, _, err = c.updateResources(pod, podScale, true)

	if err != nil {
		klog.Error("Error while performing dry-run resource update: ", err)
		return nil, nil, err
	}

	return c.updateResources(pod, podScale, false)
}

// updateResources performs Pod and PodScale resource update in dry-run mode or not whether the corresponding flag is passed
func (c *Controller) updateResources(pod *corev1.Pod, podScale *v1beta1.PodScale, dryRun bool) (newPod *corev1.Pod, newPodScale *v1beta1.PodScale, err error) {

	opts := &metav1.UpdateOptions{}

	if dryRun {
		opts.DryRun = []string{metav1.DryRunAll}
	}

	newPod, err = c.kubernetesClientset.CoreV1().Pods(podScale.Spec.Namespace).Update(context.TODO(), pod, *opts)

	if err != nil {
		klog.Error("Error updating the pod: ", err)
		return nil, nil, err
	}

	newPodScale, err = c.podScalesClientset.SystemautoscalerV1beta1().PodScales(podScale.Namespace).Update(context.TODO(), podScale, *opts)

	if err != nil {
		klog.Error("Error updating the pod scale: ", err)
		return nil, nil, err
	}

	return newPod, newPodScale, nil
}
