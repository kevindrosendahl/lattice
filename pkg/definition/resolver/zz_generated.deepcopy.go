// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package resolver

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResolutionTree.
func (in *ResolutionTree) DeepCopy() *ResolutionTree {
	if in == nil {
		return nil
	}
	out := new(ResolutionTree)
	in.DeepCopyInto(out)
	return out
}