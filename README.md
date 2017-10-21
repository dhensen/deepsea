# Deepsea

Web hosting solution built on top of Kubernetes.

## Infrastructure

I define infrastructure loosely as Virtual machines + Kubernetes.
Local: minikube/vagrant
Live: digital-ocean

## Provisioning

On top of Kubernetes there will be done some provisioning:
- Nginx Ingress Controller
- Kube Lego
- MySQL (I could also do this at Pod level, or optional depending on customers requirements)
- Deepsea
- Docker Registry (Minikube only)
- Secrets and ConfigMaps for above programs

Configure which Kubernetes cluster to use and provision using terraform/ansible/???

## Design
See [design](DESIGN.md)

## Todo
See [todo](TODO.md)
