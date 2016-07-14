package disco

/*

* disco interface

 */

type IPv4Address string

type P2PDiscovery interface {
	// runs at the very beginning, and on every host
	// every host is given the address of at
	// least one other host on the real network
	Bootstrap(peers []IPv4Address)
}

type TokenDiscovery interface {
	// runs on the host where `init` is called
	// returns a token which the user will pass to other nodes
	Init() (token string)
	// the same token is passed to Bootstrap() on all
	// the other nodes out-of-band by the user
	Bootstrap(token string)
}

/*

* pod net interface

CNI

* some binaries and a config file, in same mount namespace as kubelet
* bootstrap the network (if necessary)

*/
type PodNetwork interface {
	// get a docker image with the required binaries and config files
	// will be run before kubelet is started
	// kubelet image will be run with --volumes-from result, and
	// will expect certain files to be in certain places
	// (documentation TBD)
	GetCNIContainerImage() (imageName string)
}

/*

* pki interfaces

 */

// could be implemented by e.g. amazon container registry,
// or by a registry which runs on the node where `init` was run.

type ContainerizedPKIServer interface {
	// generate certs and tokens as container images and then
	// serve them on the default docker registry port.
	//
	// the container images then need to be distributed
	// to the respective hosts when they call PullPKIContainerImages
	//
	// the containers also have pure ASM implementation of
	// `/bin/true` in them, because docker.
	CreatePKIContainerImages()
	// the bootstrap service is closed for business (temporal security)
	StopServing()
	// re-open bootstrapping for adding new nodes
	StartServing()
}

type ContainerisedPKIClient interface {
	// pull will get run on each node, this method gives the
	// client an address to pull images from from
	PKIRegistryAddress() (registryAddress string, err error)
}
