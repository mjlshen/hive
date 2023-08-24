// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

// ServiceProviderCredentialsApplyConfiguration represents an declarative configuration of the ServiceProviderCredentials type for use
// with apply.
type ServiceProviderCredentialsApplyConfiguration struct {
	AWS *AWSServiceProviderCredentialsApplyConfiguration `json:"aws,omitempty"`
}

// ServiceProviderCredentialsApplyConfiguration constructs an declarative configuration of the ServiceProviderCredentials type for use with
// apply.
func ServiceProviderCredentials() *ServiceProviderCredentialsApplyConfiguration {
	return &ServiceProviderCredentialsApplyConfiguration{}
}

// WithAWS sets the AWS field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the AWS field is set to the value of the last call.
func (b *ServiceProviderCredentialsApplyConfiguration) WithAWS(value *AWSServiceProviderCredentialsApplyConfiguration) *ServiceProviderCredentialsApplyConfiguration {
	b.AWS = value
	return b
}