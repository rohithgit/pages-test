package etcd

import (
	"errors"
	"strings"
	"time"

	"bitbucket-eng-sjc1.cisco.com/bitbucket/scm/specnl/spectre-base-microservice/logging"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

// This is the struct use to store the etcd client instance and keyapi instance
type SpectreEtcd struct {
	EtcdClient     *client.Client
	EtcdKeysClient *client.KeysAPI
}

var log = logging.Log.Logger

// This function establishes connection with etcd cluster and returns the necessary client and key api instance and an error if there is a problem
func NewETCDClient(urls []string, transport client.CancelableTransport) (*SpectreEtcd, error) {
	if urls == nil || len(urls) == 0 {
		panic("URLs provided to establish an etcd connection cannot be empty.")
	}
	if transport == nil {
		transport = client.DefaultTransport
	}

	cfg := client.Config{
		//Endpoints: []string{"http://127.0.0.1:2379"},
		Endpoints: urls,

		Transport: transport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
		SelectionMode:           client.EndpointSelectionPrioritizeLeader,
	}
	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// this part needs more testing
	//go func (){
	//	for {
	//		log.Info("Sync. . .")
	//		//ct, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//		err := etcdClient.AutoSync(context.Background(), 10*time.Second)
	//		//err := etcdClient.Sync(context.Background())
	//		if err == context.DeadlineExceeded || err == context.Canceled {
	//			break
	//		}
	//		log.Error("Sync error:", err)
	//		//time.Sleep(1 * time.Minute)
	//	}
	//}()

	kapi := client.NewKeysAPI(etcdClient)

	return &SpectreEtcd{
		EtcdClient:     &etcdClient,
		EtcdKeysClient: &kapi,
	}, err
}

// This is a generic Set method which can be used to store key-value pair along with custom options
func (spectreEtcd *SpectreEtcd) Set(context context.Context, key, value string, options *client.SetOptions) (*client.Response, error) {
	if spectreEtcd == nil || spectreEtcd.EtcdKeysClient == nil {
		return nil, errors.New("Etcd connection is not active.")
	}

	if !strings.HasPrefix(key, "/") {
		return nil, errors.New("Key does not contain prefix /")
	}

	etcdKeysClient := *spectreEtcd.EtcdKeysClient
	log.Infof("Setting %s key with %s value", key, value)
	resp, err := etcdKeysClient.Set(context, key, value, options)
	if err != nil {
		CheckETCDErrors(err)
		return nil, err
	} else {
		// print common key info
		log.Infof("Set is done. Metadata is %q\n", resp)
		return resp, nil
	}
}

// this is a helper method to store a key-value with default options
func (spectreEtcd *SpectreEtcd) SetKeyValue(key, value string) (*client.Response, error) {
	setOptions := &client.SetOptions{}

	log.Infof("Setting %s key with %s value", key, value)
	return spectreEtcd.Set(context.Background(), key, value, setOptions)
}

// This is a generic Get method which can be used to store key-value pair along with custom options
func (spectreEtcd *SpectreEtcd) Get(context context.Context, key string, options *client.GetOptions) (*client.Response, error) {
	if spectreEtcd == nil || spectreEtcd.EtcdKeysClient == nil {
		return nil, errors.New("Etcd connection is not active.")
	}

	etcdKeysClient := *spectreEtcd.EtcdKeysClient
	resp, err := etcdKeysClient.Get(context, key, options)
	if err != nil {
		CheckETCDErrors(err)
		return nil, err
	} else {
		// print common key info
		log.Infof("Get is done. Metadata is %q\n", resp)
		PrintNodeRecursively(resp.Node)
		return resp, nil
	}
}

// this is a helper method to get a value for a key, if it exists, with default options
func (spectreEtcd *SpectreEtcd) GetStringValue(key string) (*string, error) {
	getOptions := &client.GetOptions{
		Quorum:    true,
		Recursive: false,
	}

	log.Infof("Getting %s key", key)
	resp, err := spectreEtcd.Get(context.Background(), key, getOptions)
	if err != nil {
		return nil, err
	}
	return &(resp.Node.Value), nil
}

// this is a helper method to get a dir and its recursive nodes, if they exist
func (spectreEtcd *SpectreEtcd) GetDir(key string) (client.Nodes, error) {
	getOptions := &client.GetOptions{
		Quorum:    true,
		Recursive: true,
	}

	log.Infof("Getting Dir %s key", key)
	resp, err := spectreEtcd.Get(context.Background(), key, getOptions)
	if err != nil {
		return nil, err
	}
	if resp.Node.Dir {
		PrintNodesRecursively(resp.Node.Nodes)
		return resp.Node.Nodes, nil
	}
	return nil, errors.New("Key is not a directory")
}

// This is a generic Delete method which can be used to delete a key along with custom options
func (spectreEtcd *SpectreEtcd) Delete(context context.Context, key string, options *client.DeleteOptions) (*client.Response, error) {
	if spectreEtcd == nil || spectreEtcd.EtcdKeysClient == nil {
		return nil, errors.New("Etcd connection is not active.")
	}
	log.Infof("Deleting %s key", key)
	etcdKeysClient := *spectreEtcd.EtcdKeysClient
	resp, err := etcdKeysClient.Delete(context, key, options)
	if err != nil {
		CheckETCDErrors(err)
		return nil, err
	} else {
		// print common key info
		log.Infof("Deleting is done. Metadata is %q\n", resp)
		return resp, nil
	}
}

// this is a helper method to delete a key with default options
func (spectreEtcd *SpectreEtcd) DeleteKey(key string) (*client.Response, error) {
	delOptions := &client.DeleteOptions{
		Recursive: true,
	}
	return spectreEtcd.Delete(context.Background(), key, delOptions)
}

// this is a generic function to check for general etcd errors
func CheckETCDErrors(err error) {
	if err != nil {
		if err == context.Canceled {
			// ctx is canceled by another routine
			log.Error("Error:", err)
		} else if err == context.DeadlineExceeded {
			// ctx is attached with a deadline and it exceeded
			log.Error("Error:", err)
		} else if cerr, ok := err.(*client.ClusterError); ok {
			// process (cerr.Errors)
			log.Error("Error:", cerr)
		} else {
			// bad cluster endpoints, which are not etcd servers
			log.Error("Error:", err)
		}
	}
}

// this is a helper method to print a node recursively
func PrintNodeRecursively(node *client.Node) {
	if node.Dir {
		log.Infof("%q key is a dir\n", node.Key)
		PrintNodesRecursively(node.Nodes)
	} else {
		log.Infof("%q key has %q value\n", node.Key, node.Value)
	}
}

// this is a helper method to print nodes recursively
func PrintNodesRecursively(nodes client.Nodes) {
	if nodes != nil {
		for _, node := range nodes {
			PrintNodeRecursively(node)
		}
	}
}
