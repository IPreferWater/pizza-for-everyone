# Context

Simple example of how horizontal scaling can help handle heavy traffic on API

In this example the **pizza-center-api** use a defined amount of memory to create a pizza order

**clients-simulator** will call the **pizza-center-api** and display how many calls were sucessfull/timeout/refused

**k8s** folder contain the yaml to deploy the cluster

view the logs of the api
```
kubectl logs -f -l app=pizza-center-api
```