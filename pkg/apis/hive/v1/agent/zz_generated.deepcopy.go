// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package agent

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AgentNetwork) DeepCopyInto(out *AgentNetwork) {
	*out = *in
	if in.AgentRefs != nil {
		in, out := &in.AgentRefs, &out.AgentRefs
		*out = make([]AgentRef, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AgentNetwork.
func (in *AgentNetwork) DeepCopy() *AgentNetwork {
	if in == nil {
		return nil
	}
	out := new(AgentNetwork)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AgentRef) DeepCopyInto(out *AgentRef) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AgentRef.
func (in *AgentRef) DeepCopy() *AgentRef {
	if in == nil {
		return nil
	}
	out := new(AgentRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterNetworkEntry) DeepCopyInto(out *ClusterNetworkEntry) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterNetworkEntry.
func (in *ClusterNetworkEntry) DeepCopy() *ClusterNetworkEntry {
	if in == nil {
		return nil
	}
	out := new(ClusterNetworkEntry)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstallStrategy) DeepCopyInto(out *InstallStrategy) {
	*out = *in
	in.Networking.DeepCopyInto(&out.Networking)
	in.ProvisionRequirements.DeepCopyInto(&out.ProvisionRequirements)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstallStrategy.
func (in *InstallStrategy) DeepCopy() *InstallStrategy {
	if in == nil {
		return nil
	}
	out := new(InstallStrategy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstallStrategyStatus) DeepCopyInto(out *InstallStrategyStatus) {
	*out = *in
	if in.AgentNetworks != nil {
		in, out := &in.AgentNetworks, &out.AgentNetworks
		*out = make([]AgentNetwork, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstallStrategyStatus.
func (in *InstallStrategyStatus) DeepCopy() *InstallStrategyStatus {
	if in == nil {
		return nil
	}
	out := new(InstallStrategyStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MachineNetworkEntry) DeepCopyInto(out *MachineNetworkEntry) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MachineNetworkEntry.
func (in *MachineNetworkEntry) DeepCopy() *MachineNetworkEntry {
	if in == nil {
		return nil
	}
	out := new(MachineNetworkEntry)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Networking) DeepCopyInto(out *Networking) {
	*out = *in
	if in.MachineNetwork != nil {
		in, out := &in.MachineNetwork, &out.MachineNetwork
		*out = make([]MachineNetworkEntry, len(*in))
		copy(*out, *in)
	}
	if in.ClusterNetwork != nil {
		in, out := &in.ClusterNetwork, &out.ClusterNetwork
		*out = make([]ClusterNetworkEntry, len(*in))
		copy(*out, *in)
	}
	if in.ServiceNetwork != nil {
		in, out := &in.ServiceNetwork, &out.ServiceNetwork
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Networking.
func (in *Networking) DeepCopy() *Networking {
	if in == nil {
		return nil
	}
	out := new(Networking)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Platform) DeepCopyInto(out *Platform) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Platform.
func (in *Platform) DeepCopy() *Platform {
	if in == nil {
		return nil
	}
	out := new(Platform)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProvisionRequirements) DeepCopyInto(out *ProvisionRequirements) {
	*out = *in
	in.AgentSelector.DeepCopyInto(&out.AgentSelector)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProvisionRequirements.
func (in *ProvisionRequirements) DeepCopy() *ProvisionRequirements {
	if in == nil {
		return nil
	}
	out := new(ProvisionRequirements)
	in.DeepCopyInto(out)
	return out
}
