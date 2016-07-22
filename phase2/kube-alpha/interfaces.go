package main

type ClusterInfo struct {
	nodeIPs   []string
	masterIPs []string
	// we will also need a field for DNS name(s)
	// or any static public IP, whatever is going
	// to be used to access the API server externally
}

// - any given function can be a no-op
// - there is definetly a need for the idea of 'FirstMaster',
//   all the some implementations might not have to use it

// this is top-level interface, rigth now it's here for sanity
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

type FirstMasterPKI interface {
	Init()
	Publish()
}

type AnyMasterPKI interface {
	Fetch()
	Init()
}

type NodePKI interface {
	Fetch()
	Init()
}

type FirstMasterNetworkProvider interface {
	Init()
	Launch()
}

type AnyMasterNetworkProvider interface {
	Init()
	Launch()
}

type NodeNetworkProvider interface {
	Init()
	Launch()
}
