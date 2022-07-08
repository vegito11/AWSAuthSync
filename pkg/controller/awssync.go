package controller

import (
	"time"

	crclntset "github.com/vegito11/AWSAuthSync/pkg/client/clientset/versioned"
	crscheme "github.com/vegito11/AWSAuthSync/pkg/client/clientset/versioned/scheme"
	crinf "github.com/vegito11/AWSAuthSync/pkg/client/informers/externalversions/vegito11.io/v1beta"
	crlister "github.com/vegito11/AWSAuthSync/pkg/client/listers/vegito11.io/v1beta"
	corev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

const (
	wqname = "AWSAuthSyncer"
)

type Controller struct {
	kubeClient kubernetes.Interface
	crClient   crclntset.Interface
	crSynced   cache.InformerSynced
	crLister   crlister.AWSAuthMapLister
	wq         workqueue.RateLimitingInterface
	recorder   record.EventRecorder
}

func NewController(kubeClient kubernetes.Interface, crClient crclntset.Interface, crInformer crinf.AWSAuthMapInformer) *Controller {

	utilruntime.Must(crscheme.AddToScheme(scheme.Scheme))
	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "AWSAuthMap"})

	ctr := &Controller{
		kubeClient: kubeClient,
		crClient:   crClient,
		crSynced:   crInformer.Informer().HasSynced,
		crLister:   crInformer.Lister(),
		wq:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), wqname),
		recorder:   recorder,
	}

	crInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(new interface{}) {
				klog.Info(" New AWSAuthMap has been added to the cluster")
				ctr.enqueueCR(new)
			},
			UpdateFunc: func(old, new interface{}) {
				klog.Info(" AWSAuthMap has been updated")
				ctr.enqueueCR(new)
			},
			DeleteFunc: func(new interface{}) {
				klog.Info(" AWSAuthMap has been Deleted from the cluster")
				ctr.enqueueCR(new)
			},
		},
	)

	return ctr
}

func (c *Controller) Run(ch <-chan struct{}) {
	klog.Info(" Starting Controller ⚡")

	// Informer maintain local cache , so we have to wait till cache to be synced
	// at least for one time
	if !cache.WaitForCacheSync(ch, c.crSynced) {
		klog.Error(" Waiting for cache to be synced \n")
	}
	// It will call c.worker function after every 1 second
	go wait.Until(c.worker, 1*time.Second, ch)
	<-ch
}

func (c *Controller) worker() {
	for c.processItem() {

	}
}
