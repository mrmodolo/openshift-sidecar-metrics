# Obtendo Métricas de Containers Com Ciclo de Vida Curto

## Problema

Obter métricas de um container que não está sempre em execução é complicado! Não existe no kubernets um histórico de métricas para um dado container.

## Alternativas


## Obtendo as Métricas pela API

Devem ser definidas as seguintes variáveis de ambiente:

- NAMESPACE
- POD
- APISERVER
- TOKEN
- CA

O TOKEN e a CA estão em /var/run/secrets/kubernetes.io/serviceaccount/token
e /var/run/secrets/kubernetes.io/serviceaccount/ca.crt

O valor para o APISERVER é https://kubernetes.default

```bash
curl --header "Authorization: Bearer ${TOKEN}" \
       --cacert ${CA} \
       -X GET "${APISERVER}/apis/metrics.k8s.io/v1beta1/namespaces/${NAMESPACE}/pods/${POD}"
```
A conta default deve ter privilégios de administrador no namespace!

```bash
oc policy add-role-to-user admin system:serviceaccount:${NAMESPACE}:default
```


## Referências

[Run multiple services in a container](https://docs.docker.com/config/containers/multi-service_container/)

Ver nota [2](#2)

[Multiple threads inside docker container](https://stackoverflow.com/questions/37657280/multiple-threads-inside-docker-container)

Ver nota [2](#2)

[Run Multiple Processes in a Container](https://runnable.com/docker/rails/run-multiple-processes-in-a-container)

Ver nota [2](#2)

[Runtime metrics](https://docs.docker.com/config/containers/runmetrics/)

[Docker adoption pathway - Part 1](http://livepersoninc.github.io/techblog/docker-adoption-pathway-part01.html)

[Docker adoption pathway - Part 2](http://livepersoninc.github.io/techblog/docker-adoption-pathway-part02.html)

[Getting CPU and Memory usage of a docker container from within the dockerized application](https://stackoverflow.com/questions/51248144/getting-cpu-and-memory-usage-of-a-docker-container-from-within-the-dockerized-ap)

Ver nota [1](#1)

[How to expose kubernetes metric server api to curl from inside the pod&#63;](https://stackoverflow.com/questions/58911806/how-to-expose-kubernetes-metric-server-api-to-curl-from-inside-the-pod)

[Access Clusters Using the Kubernetes API](https://kubernetes.io/docs/tasks/administer-cluster/access-cluster-api/)
Acesso com exemplos em Go, Python, Haskell e outros. Exemplo também de acesso direto com curl

[OpenShift REST API Overview](https://docs.openshift.com/container-platform/3.11/rest_api/index.html)
Exemplos de acesso

[Resource metrics pipeline](https://kubernetes.io/docs/tasks/debug-application-cluster/resource-metrics-pipeline/)
Resource usage metrics, such as container CPU and memory usage, are available in Kubernetes through the Metrics API. These metrics can be either accessed directly by user, for example by using kubectl top command, or used by a controller in the cluster, e.g. Horizontal Pod Autoscaler, to make decisions.

[OpenShift Service Accounts](https://docs.openshift.com/container-platform/3.11/dev_guide/service_accounts.html)


## Notas

[1](#1) Aqui a ideia era obter as métricas em /proc/self/cgroup, /sys/fs/cgroup/memory/memory.usage_in_bytes e /sys/fs/cgroup/cpuacct/cpuacct.usage

[2](#2) Qual é o comportamento quando temos mais de um processo no mesmo container?
