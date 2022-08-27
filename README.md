# abstract-workload
This is an example of operator development using the [Kubebuilder operator framework](https://book.kubebuilder.io/).</br>
Detailed explanation of this project can be found [here](https://medium.com/@marom.itamar)

`AbstractWorkload` - is a Kubernetes CRD which serve as an abstraction for stateless and stateful applications.

## API reference:
```yaml
apiVersion: examples.itamar.marom/v1alpha1
kind: AbstractWorkload
metadata:
  name: 
spec:
  containerImage: # Container Image to deploy
  replicas: # Number of replicas to deploy
  workloadType: # stateful / stateless
```

## Examples
### Stateless
```yaml
apiVersion: examples.itamar.marom/v1alpha1
kind: AbstractWorkload
metadata:
  name: test-stateless
spec:
  containerImage: "nginx:latest"
  replicas: 2
  workloadType: stateless
```
### Stateful
```yaml
apiVersion: examples.itamar.marom/v1alpha1
kind: AbstractWorkload
metadata:
  name: test-stateful
spec:
  containerImage: "nginx:latest"
  replicas: 1
  workloadType: stateful
```