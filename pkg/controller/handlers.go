package controller

import (
	"context"

	"github.com/vegito11/AWSAuthSync/pkg/apis/vegito11.io/v1beta"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

const (
	aws_auth_name = "aws-auth-test"
	aws_auth_ns   = "kube-system"
)

func (c *Controller) enqueueCR(obj interface{}) {

	var err error

	if _, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.wq.Add(obj)
}

func (c *Controller) getExistingMap() ([]v1beta.MapUser, []v1beta.MapRole, error) {

	aws_auth_map, getErr := c.kubeClient.CoreV1().ConfigMaps(aws_auth_ns).Get(context.Background(), aws_auth_name, metav1.GetOptions{})

	if getErr != nil {
		klog.Error(" Error getting aws auth map : ", getErr.Error())
		return nil, nil, getErr
	}

	var existingUsers []v1beta.MapUser
	convErr := yaml.Unmarshal([]byte(aws_auth_map.Data["mapUsers"]), &existingUsers)
	if convErr != nil {
		klog.Error("Error converting existing awsauth MapUser to go Struct ")
		return nil, nil, convErr
	}

	var existingRoles []v1beta.MapRole
	convErr = yaml.Unmarshal([]byte(aws_auth_map.Data["mapRoles"]), &existingRoles)
	if convErr != nil {
		klog.Error("Error converting existing awsauth MapRoles to go Struct ")
		return nil, nil, convErr
	}

	return existingUsers, existingRoles, nil
}

func (c *Controller) deleteAuthEntries(delobj v1beta.AWSAuthMap) error {
	aws_auth_map, _ := c.kubeClient.CoreV1().ConfigMaps(aws_auth_ns).Get(context.Background(), aws_auth_name, metav1.GetOptions{})

	existingUsers, existingRoles, getErr := c.getExistingMap()

	if getErr != nil {
		return getErr
	}

	for _, delusr := range delobj.Spec.MapUsers {
		for ind, user := range existingUsers {
			if delusr.UserARN == user.UserARN {
				existingUsers = append(existingUsers[:ind], existingUsers[ind+1:]...)
				break
			}
		}
	}

	for _, delrole := range delobj.Spec.MapRoles {
		for ind, role := range existingRoles {
			if delrole.RoleARN == role.RoleARN {
				existingRoles = append(existingRoles[:ind], existingRoles[ind+1:]...)
				break
			}
		}
	}

	usersStr, _ := yaml.Marshal(&existingUsers)
	aws_auth_map.Data["mapUsers"] = string(usersStr)

	roleStr, _ := yaml.Marshal(&existingRoles)
	aws_auth_map.Data["mapRoles"] = string(roleStr)

	_, upErr := c.kubeClient.CoreV1().ConfigMaps(aws_auth_ns).Update(context.TODO(), aws_auth_map, metav1.UpdateOptions{})

	if upErr != nil {
		klog.Error(" Error while updating configmap :", upErr.Error())
		return upErr
	}

	return nil
}

func (c *Controller) addAuthEntries(addobj v1beta.AWSAuthMap) error {
	aws_auth_map, _ := c.kubeClient.CoreV1().ConfigMaps(aws_auth_ns).Get(context.Background(), aws_auth_name, metav1.GetOptions{})

	existingUsers, existingRoles, getErr := c.getExistingMap()

	if getErr != nil {
		return getErr
	}

	for _, addusr := range addobj.Spec.MapUsers {
		flg := true
		for _, user := range existingUsers {
			if addusr.UserARN == user.UserARN {
				flg = false
				break
			}
		}

		if flg {
			existingUsers = append(existingUsers, addusr)
		}
	}

	for _, addrole := range addobj.Spec.MapRoles {
		flg := true
		for _, role := range existingRoles {
			if addrole.RoleARN == role.RoleARN {
				flg = false
				break
			}
		}

		if flg {
			existingRoles = append(existingRoles, addrole)
		}
	}

	usersStr, _ := yaml.Marshal(&existingUsers)
	aws_auth_map.Data["mapUsers"] = string(usersStr)

	roleStr, _ := yaml.Marshal(&existingRoles)
	aws_auth_map.Data["mapRoles"] = string(roleStr)

	_, upErr := c.kubeClient.CoreV1().ConfigMaps(aws_auth_ns).Update(context.TODO(), aws_auth_map, metav1.UpdateOptions{})

	if upErr != nil {
		klog.Error(" Error while updating configmap :", upErr.Error())
		return upErr
	}

	return nil

}
