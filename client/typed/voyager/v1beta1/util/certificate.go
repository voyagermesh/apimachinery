package util

import (
	"encoding/json"
	"fmt"

	"github.com/appscode/kutil"
	api "github.com/appscode/voyager/apis/voyager/v1beta1"
	cs "github.com/appscode/voyager/client/typed/voyager/v1beta1"
	"github.com/golang/glog"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/jsonmergepatch"
	"k8s.io/apimachinery/pkg/util/wait"
)

func EnsureCertificate(c cs.VoyagerV1beta1Interface, meta metav1.ObjectMeta, transform func(alert *api.Certificate) *api.Certificate) (*api.Certificate, error) {
	return CreateOrPatchCertificate(c, meta, transform)
}

func CreateOrPatchCertificate(c cs.VoyagerV1beta1Interface, meta metav1.ObjectMeta, transform func(alert *api.Certificate) *api.Certificate) (*api.Certificate, error) {
	cur, err := c.Certificates(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		glog.V(3).Infof("Creating Certificate %s/%s.", meta.Namespace, meta.Name)
		return c.Certificates(meta.Namespace).Create(transform(&api.Certificate{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Certificate",
				APIVersion: api.SchemeGroupVersion.String(),
			},
			ObjectMeta: meta,
		}))
	} else if err != nil {
		return nil, err
	}
	return PatchCertificate(c, cur, transform)
}

func PatchCertificate(c cs.VoyagerV1beta1Interface, cur *api.Certificate, transform func(*api.Certificate) *api.Certificate) (*api.Certificate, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, err
	}

	modJson, err := json.Marshal(transform(cur.DeepCopy()))
	if err != nil {
		return nil, err
	}

	patch, err := jsonmergepatch.CreateThreeWayJSONMergePatch(curJson, modJson, curJson)
	if err != nil {
		return nil, err
	}
	if len(patch) == 0 || string(patch) == "{}" {
		return cur, nil
	}
	glog.V(3).Infof("Patching Certificate %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	result, err := c.Certificates(cur.Namespace).Patch(cur.Name, types.MergePatchType, patch)
	return result, err
}

func TryPatchCertificate(c cs.VoyagerV1beta1Interface, meta metav1.ObjectMeta, transform func(*api.Certificate) *api.Certificate) (result *api.Certificate, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.Certificates(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = PatchCertificate(c, cur, transform)
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to patch Certificate %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to patch Certificate %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}

func TryUpdateCertificate(c cs.VoyagerV1beta1Interface, meta metav1.ObjectMeta, transform func(*api.Certificate) *api.Certificate) (result *api.Certificate, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.Certificates(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.Certificates(cur.Namespace).Update(transform(cur.DeepCopy()))
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to update Certificate %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update Certificate %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}
