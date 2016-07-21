package main

type ClusterInfo struct {
	NodeIPs   []string
	MasterIPs []string
	// we will also need a field for DNS name(s)
	// or any static public IP, whatever is going
	// to be used to access the API server externally
}

// - any given function can be a no-op
// - there is definetly a need for the idea of 'FirstMaster',
//   all the some implementations might not have to use it

// this is top-level interface, right now it's here for sanity
// perhaps we should actually have something like this and use
// to package collections of implementations, as probably very
// few of the implementations will work well together, or even
// make sense to compose
type BootstrapCluster interface {
	InitFirstMaster()
	InitAnyMaster()
	InitNode()
}

type FirstMasterRendezvousProvider interface {
	Setup() // install and/or configure
	Launch()
}

type AnyMasterRendezvousProvider interface {
	Setup()
	Launch()
}

type NodeRendezvousProvider interface {
	Setup()
	Launch()
}

// For PKI, on one special node we initialize the PKI material and make it
// available to the other nodes (somehow); on all other nodes we just fetch
// that material and then install it.
type FirstMasterPKI interface {
	Init()
	Publish()
}

type AnyServerPKI interface {
    // XXX is there information conveyed in the type of server (master vs node)
    // we run this on?
	Fetch()
	Install()
}

// Pod network setup involves installing a CNI plugin, and then doing any
// backend-specific launch actions.
type PodNetworkProvider interface {
    // XXX is there information conveyed in the type of server (master vs node)
    // we run this on?
	Install()
	Launch()
}
