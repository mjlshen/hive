package utils

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"

	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/openshift/hive/pkg/constants"
	"github.com/openshift/hive/pkg/test/clusterdeployment"
	"github.com/openshift/hive/pkg/test/generic"
)

func TestIsDeleteProtected(t *testing.T) {
	cases := []struct {
		name           string
		absent         bool
		value          string
		expectedResult bool
	}{
		{
			name:           "absent",
			absent:         true,
			expectedResult: false,
		},
		{
			name:           "true",
			value:          "true",
			expectedResult: true,
		},
		{
			name:           "false",
			value:          "false",
			expectedResult: false,
		},
		{
			name:           "empty",
			value:          "",
			expectedResult: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var options []clusterdeployment.Option
			if !tc.absent {
				options = append(
					options,
					clusterdeployment.Generic(generic.WithAnnotation(constants.ProtectedDeleteAnnotation, tc.value)),
				)
			}
			cd := clusterdeployment.Build(options...)
			actualResult := IsDeleteProtected(cd)
			assert.Equal(t, tc.expectedResult, actualResult, "unexpected result")
		})
	}
}

func TestIsClusterPausedOrRelocating(t *testing.T) {
	cases := []struct {
		name     string
		cd       *hivev1.ClusterDeployment
		expected bool
	}{
		{
			name:     "no annotation",
			cd:       clusterdeployment.Build(),
			expected: false,
		},
		{
			name: "syncset annotation true",
			cd: clusterdeployment.Build(
				clusterdeployment.Generic(generic.WithAnnotation(constants.SyncsetPauseAnnotation, "true")),
			),
			expected: true,
		},
		{
			name: "syncset annotation false",
			cd: clusterdeployment.Build(
				clusterdeployment.Generic(generic.WithAnnotation(constants.SyncsetPauseAnnotation, "false")),
			),
			expected: false,
		},
		{
			name: "syncset annotation not parsable",
			cd: clusterdeployment.Build(
				clusterdeployment.Generic(generic.WithAnnotation(constants.SyncsetPauseAnnotation, "other")),
			),
			expected: false,
		},
		{
			name: "relocate annotation",
			cd: clusterdeployment.Build(
				clusterdeployment.Generic(generic.WithAnnotation(constants.RelocateAnnotation, "some-relocate/outgoing")),
			),
			expected: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := IsClusterPausedOrRelocating(tc.cd, logrus.StandardLogger())
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestIsRelocating(t *testing.T) {
	cases := []struct {
		name                 string
		obj                  hivev1.MetaRuntimeObject
		expectedRelocateName string
		expectedStatus       hivev1.RelocateStatus
		expectError          bool
	}{
		{
			name: "no annotation",
			obj:  clusterdeployment.Build(),
		},
		{
			name: "valid annotation",
			obj: clusterdeployment.Build(
				clusterdeployment.Generic(generic.WithAnnotation(constants.RelocateAnnotation, "test-relocate/outgoing")),
			),
			expectedRelocateName: "test-relocate",
			expectedStatus:       hivev1.RelocateOutgoing,
		},
		{
			name: "malformed annotation",
			obj: clusterdeployment.Build(
				clusterdeployment.Generic(generic.WithAnnotation(constants.RelocateAnnotation, "bad-value")),
			),
			expectError: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actualRelocateName, actualStatus, actualError := IsRelocating(tc.obj)
			if tc.expectError {
				assert.Error(t, actualError, "expected error")
				return
			}
			require.NoError(t, actualError, "unexpected error")
			assert.Equal(t, tc.expectedRelocateName, actualRelocateName, "unexpected relocate name")
			assert.Equal(t, tc.expectedStatus, actualStatus, "unexpected relocate status")
		})
	}
}

func TestSetRelocateAnnotation(t *testing.T) {
	cases := []struct {
		name                    string
		obj                     hivev1.MetaRuntimeObject
		relocateName            string
		relocateStatus          hivev1.RelocateStatus
		expectedAnnotationValue string
		expectedChanged         bool
	}{
		{
			name:                    "new annotation",
			obj:                     clusterdeployment.Build(),
			relocateName:            "test-relocate",
			relocateStatus:          hivev1.RelocateOutgoing,
			expectedAnnotationValue: "test-relocate/outgoing",
			expectedChanged:         true,
		},
		{
			name: "replace annotation",
			obj: clusterdeployment.Build(
				clusterdeployment.Generic(generic.WithAnnotation(constants.RelocateAnnotation, "other-relocate/outgoing")),
			),
			relocateName:            "test-relocate",
			relocateStatus:          hivev1.RelocateOutgoing,
			expectedAnnotationValue: "test-relocate/outgoing",
			expectedChanged:         true,
		},
		{
			name: "no change",
			obj: clusterdeployment.Build(
				clusterdeployment.Generic(generic.WithAnnotation(constants.RelocateAnnotation, "test-relocate/outgoing")),
			),
			relocateName:            "test-relocate",
			relocateStatus:          hivev1.RelocateOutgoing,
			expectedAnnotationValue: "test-relocate/outgoing",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actualChanged := SetRelocateAnnotation(tc.obj, tc.relocateName, tc.relocateStatus)
			actualAnnotationValue := tc.obj.GetAnnotations()[constants.RelocateAnnotation]
			assert.Equal(t, tc.expectedAnnotationValue, actualAnnotationValue, "unexpected annotation value")
			assert.Equal(t, tc.expectedChanged, actualChanged, "unexpected changed result")
		})
	}
}

func TestClearRelocateAnnotation(t *testing.T) {
	cases := []struct {
		name            string
		obj             hivev1.MetaRuntimeObject
		expectedChanged bool
	}{
		{
			name: "no annotation",
			obj:  clusterdeployment.Build(),
		},
		{
			name: "existing annotation",
			obj: clusterdeployment.Build(
				clusterdeployment.Generic(generic.WithAnnotation(constants.RelocateAnnotation, "test-relocate/outgoing")),
			),
			expectedChanged: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actualChanged := ClearRelocateAnnotation(tc.obj)
			assert.NotContains(t, tc.obj.GetAnnotations(), constants.RelocateAnnotation, "unexpected relocate annotation")
			assert.Equal(t, tc.expectedChanged, actualChanged, "unexpected changed result")
		})
	}
}

func TestIsClusterDeploymentErrorUpdateEvent(t *testing.T) {
	cases := []struct {
		name     string
		old, new client.Object

		exp bool
	}{{
		name: "no conditions",
		old:  &hivev1.ClusterDeployment{},
		new:  &hivev1.ClusterDeployment{},
		exp:  false,
	}, {
		name: "all in desired state",
		old:  &hivev1.ClusterDeployment{},
		new: &hivev1.ClusterDeployment{
			Status: hivev1.ClusterDeploymentStatus{
				Conditions: []hivev1.ClusterDeploymentCondition{{
					Type:   hivev1.ProvisionFailedCondition,
					Status: corev1.ConditionFalse,
				}, {
					Type:   hivev1.AWSPrivateLinkReadyClusterDeploymentCondition,
					Status: corev1.ConditionTrue,
				}},
			},
		},
		exp: false,
	}, {
		name: "condition in undesired, previously missing",
		old: &hivev1.ClusterDeployment{
			Status: hivev1.ClusterDeploymentStatus{
				Conditions: []hivev1.ClusterDeploymentCondition{{
					Type:   hivev1.ProvisionFailedCondition,
					Status: corev1.ConditionFalse,
				}},
			},
		},
		new: &hivev1.ClusterDeployment{
			Status: hivev1.ClusterDeploymentStatus{
				Conditions: []hivev1.ClusterDeploymentCondition{{
					Type:   hivev1.ProvisionFailedCondition,
					Status: corev1.ConditionFalse,
				}, {
					Type:   hivev1.AWSPrivateLinkReadyClusterDeploymentCondition,
					Status: corev1.ConditionFalse,
				}},
			},
		},
		exp: true,
	}, {
		name: "condition in undesired, previously in desired",
		old: &hivev1.ClusterDeployment{
			Status: hivev1.ClusterDeploymentStatus{
				Conditions: []hivev1.ClusterDeploymentCondition{{
					Type:   hivev1.ProvisionFailedCondition,
					Status: corev1.ConditionFalse,
				}, {
					Type:   hivev1.AWSPrivateLinkReadyClusterDeploymentCondition,
					Status: corev1.ConditionTrue,
				}},
			},
		},
		new: &hivev1.ClusterDeployment{
			Status: hivev1.ClusterDeploymentStatus{
				Conditions: []hivev1.ClusterDeploymentCondition{{
					Type:   hivev1.ProvisionFailedCondition,
					Status: corev1.ConditionFalse,
				}, {
					Type:   hivev1.AWSPrivateLinkReadyClusterDeploymentCondition,
					Status: corev1.ConditionFalse,
				}},
			},
		},
		exp: true,
	}, {
		name: "condition in undesired, previously unknown",
		old: &hivev1.ClusterDeployment{
			Status: hivev1.ClusterDeploymentStatus{
				Conditions: []hivev1.ClusterDeploymentCondition{{
					Type:   hivev1.ProvisionFailedCondition,
					Status: corev1.ConditionFalse,
				}, {
					Type:   hivev1.AWSPrivateLinkReadyClusterDeploymentCondition,
					Status: corev1.ConditionUnknown,
				}},
			},
		},
		new: &hivev1.ClusterDeployment{
			Status: hivev1.ClusterDeploymentStatus{
				Conditions: []hivev1.ClusterDeploymentCondition{{
					Type:   hivev1.ProvisionFailedCondition,
					Status: corev1.ConditionFalse,
				}, {
					Type:   hivev1.AWSPrivateLinkReadyClusterDeploymentCondition,
					Status: corev1.ConditionFalse,
				}},
			},
		},
		exp: false,
	}, {
		name: "condition in undesired, previously unknown with new error",
		old: &hivev1.ClusterDeployment{
			Status: hivev1.ClusterDeploymentStatus{
				Conditions: []hivev1.ClusterDeploymentCondition{{
					Type:   hivev1.ProvisionFailedCondition,
					Status: corev1.ConditionFalse,
				}, {
					Type:   hivev1.AWSPrivateLinkReadyClusterDeploymentCondition,
					Status: corev1.ConditionUnknown,
					Reason: "SomeReason",
				}},
			},
		},
		new: &hivev1.ClusterDeployment{
			Status: hivev1.ClusterDeploymentStatus{
				Conditions: []hivev1.ClusterDeploymentCondition{{
					Type:   hivev1.ProvisionFailedCondition,
					Status: corev1.ConditionFalse,
				}, {
					Type:   hivev1.AWSPrivateLinkReadyClusterDeploymentCondition,
					Status: corev1.ConditionFalse,
					Reason: "SomeAnotherReason",
				}},
			},
		},
		exp: true,
	}, {
		name: "condition in undesired, previously in undesired",
		old: &hivev1.ClusterDeployment{
			Status: hivev1.ClusterDeploymentStatus{
				Conditions: []hivev1.ClusterDeploymentCondition{{
					Type:   hivev1.ProvisionFailedCondition,
					Status: corev1.ConditionFalse,
				}, {
					Type:   hivev1.AWSPrivateLinkReadyClusterDeploymentCondition,
					Status: corev1.ConditionFalse,
				}},
			},
		},
		new: &hivev1.ClusterDeployment{
			Status: hivev1.ClusterDeploymentStatus{
				Conditions: []hivev1.ClusterDeploymentCondition{{
					Type:   hivev1.ProvisionFailedCondition,
					Status: corev1.ConditionFalse,
				}, {
					Type:   hivev1.AWSPrivateLinkReadyClusterDeploymentCondition,
					Status: corev1.ConditionFalse,
				}},
			},
		},
		exp: false,
	}, {
		name: "condition in undesired, previously in undesired with new error",
		old: &hivev1.ClusterDeployment{
			Status: hivev1.ClusterDeploymentStatus{
				Conditions: []hivev1.ClusterDeploymentCondition{{
					Type:   hivev1.ProvisionFailedCondition,
					Status: corev1.ConditionFalse,
				}, {
					Type:   hivev1.AWSPrivateLinkReadyClusterDeploymentCondition,
					Status: corev1.ConditionFalse,
					Reason: "SomeReason",
				}},
			},
		},
		new: &hivev1.ClusterDeployment{
			Status: hivev1.ClusterDeploymentStatus{
				Conditions: []hivev1.ClusterDeploymentCondition{{
					Type:   hivev1.ProvisionFailedCondition,
					Status: corev1.ConditionFalse,
				}, {
					Type:   hivev1.AWSPrivateLinkReadyClusterDeploymentCondition,
					Status: corev1.ConditionFalse,
					Reason: "SomeAnotherReason",
				}},
			},
		},
		exp: true,
	}, {
		name: "condition in undesired, previously in undesired with same error",
		old: &hivev1.ClusterDeployment{
			Status: hivev1.ClusterDeploymentStatus{
				Conditions: []hivev1.ClusterDeploymentCondition{{
					Type:   hivev1.ProvisionFailedCondition,
					Status: corev1.ConditionFalse,
				}, {
					Type:   hivev1.AWSPrivateLinkReadyClusterDeploymentCondition,
					Status: corev1.ConditionFalse,
					Reason: "SomeReason",
				}},
			},
		},
		new: &hivev1.ClusterDeployment{
			Status: hivev1.ClusterDeploymentStatus{
				Conditions: []hivev1.ClusterDeploymentCondition{{
					Type:   hivev1.ProvisionFailedCondition,
					Status: corev1.ConditionFalse,
				}, {
					Type:   hivev1.AWSPrivateLinkReadyClusterDeploymentCondition,
					Status: corev1.ConditionFalse,
					Reason: "SomeReason",
				}},
			},
		},
		exp: false,
	}}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got := IsClusterDeploymentErrorUpdateEvent(event.UpdateEvent{ObjectOld: test.old, ObjectNew: test.new})
			assert.Equal(t, test.exp, got)
		})
	}
}
