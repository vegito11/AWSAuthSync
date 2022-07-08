package controller

import (
	"context"

	"github.com/vegito11/AWSAuthSync/pkg/apis/vegito11.io/v1beta"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

/* Main Worker function which takes item from Queue and process it.  */
func (c *Controller) processItem() bool {

	currentAuthMap, shutdown := c.wq.Get()

	defer c.wq.Done(currentAuthMap)

	if shutdown {
		return false
	}

	// we already handled error in enqueue So we will not get any error here
	key, _ := cache.MetaNamespaceKeyFunc(currentAuthMap)
	ns, name, err := cache.SplitMetaNamespaceKey(key)

	if err != nil {
		klog.Error(" Error while spliting key into namespace and name")
		return false
	}

	authmap, getErr := c.crClient.Vegito11V1beta().AWSAuthMaps(ns).Get(context.Background(), name, metav1.GetOptions{})

	// ===================== 2) Process deleted Map ------------------
	if apierrors.IsNotFound(getErr) {
		klog.V(4).Infof(" AuthMap with %s name is not found in ns %s - processing delete event ", name, ns)

		obj := currentAuthMap.(*v1beta.AWSAuthMap)

		delErr := c.deleteAuthEntries(*obj)

		return delErr == nil

	} else if authmap.Name != "" { /* ===== 1) Process update/add event  */
		addErr := c.addAuthEntries(*authmap)

		return addErr == nil
		// c.updateStatus(ns, name)
	}

	return true
}

/* func (c *Controller) updateStatus(ns, name string) {
	authMap, _ := c.crClient.Vegito11V1beta().AWSAuthMaps(ns).Get(context.Background(), name, metav1.GetOptions{})
	authCopy := authMap.DeepCopy()
	if authMap.Status.State == "" {
		authCopy.Status.State = "Created"
	} else if authMap.Status.State == "Created" {
		authCopy.Status.State = "Updated"
	} else {
		return
	}
	c.crClient.Vegito11V1beta().AWSAuthMaps(ns).UpdateStatus(context.TODO(), authCopy, metav1.UpdateOptions{})

} */
