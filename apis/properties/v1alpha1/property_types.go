/*
Copyright 2020 The Crossplane Authors.

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

package v1alpha1

import (
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// PropertyParameters are the configurable fields of a Property.
type PropertyParameters struct {
	AccountID string `json:"accountId"`
	// +optional
	AssetID string `json:"assetId,omitempty"`
	// +optional
	ContractID string `json:"contractId,omitempty"`
	// +optional
	GroupID string `json:"groupId,omitempty"`
	// +optional
	LatestVersion int `json:"latestVersion,omitempty"`
	// +optional
	Note string `json:"note,omitempty"`
	// +optional
	ProductID string `json:"productId,omitempty"`
	// +optional
	ProductionVersion *int `json:"productionVersion,omitempty"`
	// +optional
	PropertyID string `json:"propertyId,omitempty"`
	// +optional
	PropertyName string `json:"propertyName,omitempty"`
	// +optional
	RuleFormat string `json:"ruleFormat,omitempty"`
	// +optional
	StagingVersion *int `json:"stagingVersion,omitempty"`
}

// PropertyObservation are the observable fields of a Property.
type PropertyObservation struct {
	// Response
	// Properties PropertiesItems `json:"properties"`
	// Property   *Property       `json:"-"`
	ObservableField string `json:"observableField,omitempty"`
}

// PropertiesItems is an array of properties
type PropertiesItems struct {
	Items []*Property `json:"items"`
}

// A PropertySpec defines the desired state of a Property.
type PropertySpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       PropertyParameters `json:"forProvider"`
}

// A PropertyStatus represents the observed state of a Property.
type PropertyStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          PropertyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Property is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,akamai}
type Property struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PropertySpec   `json:"spec"`
	Status PropertyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PropertyList contains a list of Property
type PropertyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Property `json:"items"`
}

// Property type metadata.
var (
	PropertyKind             = reflect.TypeOf(Property{}).Name()
	PropertyGroupKind        = schema.GroupKind{Group: Group, Kind: PropertyKind}.String()
	PropertyKindAPIVersion   = PropertyKind + "." + SchemeGroupVersion.String()
	PropertyGroupVersionKind = SchemeGroupVersion.WithKind(PropertyKind)
)

func init() {
	SchemeBuilder.Register(&Property{}, &PropertyList{})
}
