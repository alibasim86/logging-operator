// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fluentbit

import (
	"github.com/banzaicloud/logging-operator/pkg/k8sutil"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *Reconciler) clusterRole() (runtime.Object, k8sutil.DesiredState, error) {
	if *r.Logging.Spec.FluentbitSpec.Security.RoleBasedAccessControlCreate {
		return &rbacv1.ClusterRole{
			ObjectMeta: r.FluentbitObjectMetaClusterScope(clusterRoleName),
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups: []string{""},
					Resources: []string{"pods", "namespaces"},
					Verbs:     []string{"get", "list", "watch"},
				},
			},
		}, k8sutil.StatePresent, nil
	}
	return &rbacv1.ClusterRole{
		ObjectMeta: r.FluentbitObjectMetaClusterScope(clusterRoleName),
		Rules:      []rbacv1.PolicyRule{}}, k8sutil.StateAbsent, nil
}

func (r *Reconciler) clusterRoleBinding() (runtime.Object, k8sutil.DesiredState, error) {
	if *r.Logging.Spec.FluentbitSpec.Security.RoleBasedAccessControlCreate {
		return &rbacv1.ClusterRoleBinding{
			ObjectMeta: r.FluentbitObjectMetaClusterScope(clusterRoleBindingName),
			RoleRef: rbacv1.RoleRef{
				Kind:     "ClusterRole",
				APIGroup: "rbac.authorization.k8s.io",
				Name:     r.Logging.QualifiedName(clusterRoleName),
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      r.Logging.QualifiedName(defaultServiceAccountName),
					Namespace: r.Logging.Spec.ControlNamespace,
				},
			},
		}, k8sutil.StatePresent, nil
	}
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: r.FluentbitObjectMetaClusterScope(clusterRoleBindingName),
		RoleRef:    rbacv1.RoleRef{}}, k8sutil.StateAbsent, nil
}

func (r *Reconciler) serviceAccount() (runtime.Object, k8sutil.DesiredState, error) {
	if *r.Logging.Spec.FluentbitSpec.Security.RoleBasedAccessControlCreate && r.Logging.Spec.FluentbitSpec.Security.ServiceAccount == "" {
		return &corev1.ServiceAccount{
			ObjectMeta: r.FluentbitObjectMeta(defaultServiceAccountName),
		}, k8sutil.StatePresent, nil
	}
	return &corev1.ServiceAccount{
		ObjectMeta: r.FluentbitObjectMeta(defaultServiceAccountName),
	}, k8sutil.StateAbsent, nil
}
