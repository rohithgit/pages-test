The service_discovery package can be imported and consumed from other golang packages.
It can help discover services registered with consul and resolve address and Port.
The service name should have convention <microservice_name>.service.consul

To test service discovery, you will need a consul running in a container and service_discovery test code
running in another container.
Here are the steps to configure the environment.

1. Run the consul container
  sudo docker run -d -p 8400:8400 -p 8500:8500 -p 8600:53/udp --name consul1 --net my_net -h node1 progrium/consul -server –bootstrap

2. Register your service with consul. (With mantl it gets registered automatically with consul)
   For instance, we can manually register a search service like this -  
   curl -X PUT -d '{"Datacenter": "dc1", "Node": "google", "Address": "www.google.com", "Service": {"Service": "search", "Port": 80}}' http://127.0.0.1:8500/v1/catalog/register

3. Use test_service_discovery bash script and run it.

4. To test some other service, you can make changes the name of service in service_discovery_test.go
    
