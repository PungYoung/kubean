# Kubean test case

| Function              | Case Name                                                                                                           | Declaration                          | Requirement ID | Status                 | Code Link                                                                   | Test details                                                                                    |
|-----------------------|---------------------------------------------------------------------------------------------------------------------|--------------------------------------|---------------|------------------------|-----------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------|
| Create basic cluster  | Support one node cluster: master and worker in 1 node                                                               |master and worker in one node         | C-001         |                        |                                                                                     | [detail](./testcase_details/create_cluster.md#create-basic-cluster)                             |
|                       | Support k8s: 1.24                                                                                                   |                                      | C-012         |<ul><li>[x] </li></ul>  |[Code Link](../../test/kubean_sonobouy_nightlye2e/kubean_cluster_sonobouy_test.go)  | [detail](./testcase_details/create_cluster.md#support-k8s-124)                             |
|                       | Support CRI: containerd                                                                                             |                                      | C-013         | <ul><li>[x] </li></ul> | [Code Link](../../test/kubean_functions_e2e/kubean_cluster_install_test.go)        | [detail](./testcase_details/create_cluster.md#create-basic-cluster)                             |
|                       | Support CNI: calico                                                                                                 |                                      | C-001/N-37    | <ul><li>[x] </li></ul> | [Code Link](../../test/kubean_sonobouy_nightlye2e/kubean_cluster_sonobouy_test.go) | [detail](./testcase_details/create_cluster.md#support-cni-calico)                             |
|                       | Support calico tunnel mode: IPIP                                                                                    |                                      | C-001         |<ul><li>[x] </li></ul>|[Code Link](../../test/kubean_calico_nightlye2e/kubean_network_calico_test.go)       | [detail](./testcase_details/create_cluster.md#create-basic-cluster)                             |
|                       | Not overwirte hostname                                                                                              |                                      | C-001         | <ul><li>[x] </li></ul> | [Code Link](../../test/kubean_reset_e2e/kubean_cluster_reset_test.go)              | [detail](./testcase_details/create_cluster.md#not-overwrite-hostname)                             |
|                       | Support disable cluster ca auto_renew                                                                               |                                      | C-015         |                        |                                                                                    | [detail](./testcase_details/create_cluster.md#create-basic-cluster)                             |
|                       | Ssh authorization: user name and password                                                                           |                                      |               | <ul><li>[x] </li></ul> | [Code Link](../../test/kubean_sonobouy_nightlye2e/kubean_cluster_sonobouy_test.go) | [detail](./testcase_details/create_cluster.md#create-basic-cluster)                             |
| Create cluster：extend| Create cluster topology ：1 master and 1 worker                                                                      |                                      |               | <ul><li>[x] </li></ul> | [Code Link](../../test/kubean_sonobouy_nightlye2e/kubean_cluster_sonobouy_test.go) | [detail](./testcase_details/create_cluster.md#create-cluster-with-one-master-and-one-worker)    |
|                       | Create cluster topology ：3 master and 2 worker                                                                      |                                     | C-012        |                        |                                                                                    | [detail](./testcase_details/create_cluster.md#create-cluster-topology-3-master-and-2-worker)    |
|                       | Support k8s: 1.23                                                                                                    |                                      | C-012        | <ul><li>[x] </li></ul> |[Code Link](../../test/kubean_sonobouy_nightlye2e/kubean_cluster_sonobouy_test.go)   | [detail](./testcase_details/create_cluster.md#support-k8s-1230)               |
|                       | Support k8s: 1.25                                                                                                    |                                      | C-012        |                        |                                                                                     |                                                                                     |
|                       | Support CRI: docker                                                                                                  |                                      | C-001         | <ul><li>[x] </li></ul> | [Code Link](../../test/kubean_sonobouy_nightlye2e/kubean_cluster_sonobouy_test.go) | [detail](./testcase_details/create_cluster.md#support-cri-docker)             |
|                       | SSH authorization: private key                                                                                       |                                       | C-001         | <ul><li>[x] </li></ul>  | [Code Link](../../test/kubean_add_remove_worker_nightlye2e/kubean_add_remove_worker_test.go) | [detail](./testcase_details/create_cluster.md#ssh-authorization-private-key)  |
|                       | Support CRI: Cilium                                                                                                  | cilium：kernel not lower than centos8 | C-001         |                        |                                                                                    | [detail](./testcase_details/create_cluster.md#create_cluster.md#support-cricilium )             |
|                       | Support calico tunnel mode: Vxlan                                                                                   |                                      | C-001         |                        |                                                                                      |                                                                                                 |
|                       | Support calico tunnel mode: IPIP_CrossSubnet                                                                        |                                      | C-001         |                        |                                                                                      |                                                                                                 |
|                       | Support calico tunnel mode: Vxlan_CrossSubnet                                                                       |                                      | C-001         |                        |                                                                                      |                                                                                                 |
|                       | Support cilium tunnel mode: Vxlan                                                                                   |                                      | C-001         |                        |                                                                                      |                                                                                                 |
|                       | Support cilium tunnel mode: IPIP_CrossSubnet                                                                        |                                      | C-001         |                        |                                                                                      |                                                                                                 |
|                       | Support cilium tunnel mode: Vxlan_CrossSubnet                                                                       |                                      | C-001         |                        |                                                                                      |                                                                                                 |
|                       | Support cilium tunnel mode: IPIP                                                                                    |                                      | C-001         |                        |                                                                                      |                                                                                                 |
|                       | Cluster node on different sub net                                                                                   |                                      | C-001         |                        |                                                                                      |                                                                                                 |
|                       | Support CRI: runc                                                                                                    |                                      | C-042         |                        |                                                                                      |                                                                                                 |
|                       | Support CRI: kata                                                                                                    |                                      | C-042         |                        |                                                                                      |                                                                                                 |
|                       | Support CRI: gvisor                                                                                                  |                                      | C-042         |                        |                                                                                      |                                                                                                 |
|                       | Support calico: kube_pods_subnet                                                                                    | kube_pods_subnet                     | C-001         | <ul><li>[x] </li></ul> | [Code Link](../../test/kubean_functions_e2e/kubean_cluster_install_test.go)           | [detail](./testcase_details/create_cluster.md#support-kube_pods_subnet)                         |
|                       | Support calico: kube_service_addresses                                                                              | kube_service_addresses               | C-001         |                        |                                                                                      |                                                                                                  |
|                       | Support cilium: kube_pods_subnet                                                                                    |                                     | C-001          |                        |                                                                                      |                                                                                                  |
|                       | Support cilium: kube_service_addresses                                                                              |                                      | C-001         |                        |                                                                                      |                                                                                                  |
|                       | Support cluster ca auto_renew                                                                                       | auto_renew_certificates:true         | C-001         |                        |                                                                                      |                                                                                                  |
|                       | Support overwrite hostname                                                                                          | override_system_hostname=true        | C-001         | <ul><li>[x] </li></ul> | [Code Link](../../test/kubean_functions_e2e/kubean_cluster_install_test.go)          | [detail](./testcase_details/create_cluster.md#support-overwrite-hostname)                        |
|                       | Support Readhat8 OS                                                                                                 |                                      | C-001         |                        |                                                                                      | [detail](./testcase_details/create_cluster.md#support-readhat8-os)                              |
|                       | Support Centos8 OS                                                                                                  |                                      | C-001         |                        |                                                                                      | [detail](./testcase_details/create_cluster.md#support-centos8-os)                               |
|                       | Support Kylin2 OS on arm                                                                                            |                                      | C-001         |                        |                                                                                      |                                                                                                 |
|                       | Support create cluster on public clouds                                                                             |                                      | C-001         |                        |                                                                                      |                                                                                                 |
|                       | Support create cluster on physical machine                                                                          |                                      | C-001         |                        |                                                                                      |                                                                                                 |
|                       | Support calico ipv4/ipv6 dual-stack net                                                                             |                                      | C-001         |                        |                                                                                      |                                                                                                 |
|                       | Support cilium ipv4/ipv6 dual-stack net                                                                             |                                      | C-001         |                        |                                                                                      |                                                                                                 |
|                       | Support kube-vip                                                                                                    |                                      | C-001         |                        |                                                                                      |                                                                                                 |
|                       | Support ipvs mode                                                                                                   |                                      | C-001         |                        |                                                                                      |                                                                                                 |
|                       | Support iptable mode                                                                                                |                                      | C-001         |                        |                                                                                      |                                                                                                 |
|                       | Support non root user install                                                                                       |                                      | C-001         |                        |                                                                                      |                                                                                                 |
|                       | Support local DNS server                                                                                            |                                      | C-001         |                        |                                                                                           |                                                                                                 |
|                       | Support local NTP server                                                                                            |                                      | C-001         |                        |                                                                                          |                                                                                                 |
|                       | Support set kubelet_max_pods                                                                                        |                                      | C-001         |                        |                                                                                      |                                                                                                 |
|                       | Support set kubernetes_audit                                                                                        |                                      | C-001         |                        |                                                                                      |                                                                                                 |
|                       | Support local insecure registry                                                                                     |                                      | C-001         |                        |                                                                                   |                                                                                                 |
|                       | Create cluster with same node name, should be fail                                                                  |                                      | C-001         |                        |                                                                                      | [detail](./testcase_details/create_cluster.md#create-cluster-set-all-the-nodes-with-same-name)  |
| Reset                 | Cluster reset                                                                                                       |                                      |               | <ul><li>[x] </li></ul> | [Code Link](../../test/kubean_reset_e2e/kubean_cluster_reset_test.go)                | [detail](./testcase_details/cluster_operation.md#cluster-reset)                                 |
|                       | Retry 0 times when job fail                                                                                         | backoffLimit=0                       |               | <ul><li>[x] </li></ul> | [Code Link](../../test/kubeanOps_functions_nightlye2e/kubean_bol0_test.go)           | [detail](./testcase_details/cluster_operation.md#retry-0-times-when-job-fail)                   |
|                       | Retry 1 times when job fail                                                                                         | backoffLimit=1                       |               |<ul><li>[x] </li></ul>  |[Code Link](../../test/kubeanOps_functions_nightlye2e/kubean_bol1_test.go)            |                                                                                                 |
|                       | Retry 2 times when job fail                                                                                         | backoffLimit=2                       |               |<ul><li>[x] </li></ul>  |[Code Link](../../test/kubeanOps_functions_nightlye2e/kubean_bol2_test.go)            |                                                                                                 |
| Cluster Function      | Posthook  the kubeconfig  to kubean                                                                                 |                                      |               |<ul><li>[x] </li></ul>  |[Code Link](../../test/kubean_sonobouy_nightlye2e/kubean_cluster_sonobouy_test.go)    | [detail](./testcase_details/kubean_func.md#posthook-cluster-kubeconfig)                    |
|                       | Limit concurrent operations of the cluster                                                                          |                                      |               |<ul><li>[x] </li></ul>  |[Code Link](../../test/kubeanOps_functions_nightlye2e/kubeanOps_test.go)              |                        |                                                                        |                                                                                                 |
|                       | Delete the backup resources reverse order                                                                           | default reserve 5 copies             |               |<ul><li>[x] </li></ul>  |[Code Link](../../test/kubeanOps_functions_nightlye2e/kubeanOps_test.go)              |                                                                                                 |
|                       | Anti-modification flag can be set when backup resource be modified                                                  |                                      |               |                        |                                                                                      |                                                                                                 |
| Add/Remove node       | Add master                                                                                                          |                                      | C-004         |                        |                                                                                      |                                                                                                 |
|                       | Add worker                                                                                                          |                                      | C-004         | <ul><li>[x] </li></ul> |[Code Link](../../test/kubean_add_remove_worker_nightlye2e/kubean_add_remove_worker_test.go) | [detail](./testcase_details/cluster_operation.md#add-worker)                    |
|                       | Remove online worker                                                                                                |                                      | C-004         |<ul><li>[x] </li></ul>  |[Code Link](../../test/kubean_add_remove_worker_nightlye2e/kubean_add_remove_worker_test.go) | [detail](./testcase_details/cluster_operation.md#remove-online-worker)            |
|                       | Remove offline worker                                                                                               |                                      | C-004         |                        |                                                                                      | [detail](./testcase_details/cluster_operation.md#remove-offline-worker)                         |
|                       | Remove online master                                                                                                |                                      | C-004         |                        |                                                                                      |                                                                                                 |
|                       | Remove offline master                                                                                               |                                      | C-004         |                        |                                                                                      |                                                                                                 |
|                       | Online master down in remove procedure                                                                              |                                      | C-004         |                        |                                                                                      |                                                                                                 |
|                       | Online worker down in remove procedure                                                                              |                                      | C-004         |                        |                                                                                      | [detail](./testcase_details/cluster_operation.md#online-worker-down-in-remove-procedure)        |
|                       | Readd master node to cluster                                                                                        |                                      | C-015         |                        |                                                                                      |                                                                                                 |
|                       | Readd a worker to cluster                                                                                           |                                      | C-015         |                        |                                                                                      | [detail](./testcase_details/cluster_operation.md#readd-a-worker-to-cluster)                     |
|                       | Manual update ca                                                                                                    | kubeadmin renew                      | C-015         |                        |                                                                                      |                                                                                                 |
| Hot upgrade           | Hot upgrade k8s Z version: online                                                                                   |                                      | C-003         | <ul><li>[x] </li></ul> |[Code Link](../../test/kubean_sonobouy_nightlye2e/kubean_cluster_sonobouy_test.go)    |
|                       | Hot upgrade k8s Y version: online                                                                                   |                                      | C-003         | <ul><li>[x] </li></ul> |[Code Link](../../test/kubean_sonobouy_nightlye2e/kubean_cluster_sonobouy_test.go)    | [detail](./testcase_details/cluster_operation.md#hot-upgrade-k8s-y-version-online)       |
|                       | Hot upgrade k8s Y  version: offline                                                                                 |                                      | C-003         |                        |                                                                                      |                                                                                                 |
|                       | Hot upgrade k8s Z version: offline                                                                                  |                                      | C-003         |                        |                                                                                      |                                                                                                 |
|                       | Hot upgrade CNI                                                                                                     | To be determined                     | C-003         |                        |                                                                                      |                                                                                                 |
|                       | Hot degrade CNI                                                                                                     | To be determined                     | C-003         |                        |                                                                                      |                                                                                                 |
|                       | Hot upgrade CRI-containerd online                                                                                   |                                      | C-003         |                        |                                                                                      |                                                                                                 |
|                       | Hot upgrade CRI-docker online                                                                                       |                                      | C-003         |                        |                                                                                      |                                                                                                 |
|                       | Hot upgrade CRI-containerd offline                                                                                  |                                      | C-003         |                        |                                                                                      |                                                                                                 |
|                       | Hot upgrade CRI-docker offline                                                                                      |                                      | C-003         |                        |                                                                                      |                                                                                                 |
|                       | Hot upgrade Network solutions（CNI,Metricserver,dns)                                                                | To be determined                     | C-003         |                        |                                                                                      |                                                                                                 |
|                       | Use low verison of kubespray image to create cluster，then use high verion of kubespray image to upgrade k8s version | solution to be determined            | C-003         |                        |                                                                                      |                                                                                                 |
| High availability     | Hight availability of not firt master : Node crash                                                                  |                                      | L-019         |                        |                                                                                      |                                                                                                 |
|                       | Hight availability of not firt master: CPU Insufficient                                                             |                                      | L-019         |                        |                                                                                      |                                                                                                 |
|                       | Hight availability of not firt master: Disk Insufficient                                                            |                                      | L-019         |                        |                                                                                      |                                                                                                 |
|                       | Hight availability of not firt master: Memory Insufficient                                                          |                                      | L-019         |                        |                                                                                      |                                                                                                 |
|                       | Hight availability of not firt master: Network unstable                                                             |                                      | L-019         |                        |                                                                                      |                                                                                                 |
|                       | Hight availability of leader etcd                                                                                   |                                      | L-019         |                        |                                                                                      |                                                                                                 |
|                       | Hight availability of follower etcd                                                                                 |                                      | L-019         |                        |                                                                                      |                                                                                                 |
|                       | Hight availability of worker: Node crash                                                                            |                                      | L-018         |                        |                                                                                      |                                                                                                 |
|                       | Hight availability of firt master                                                                                   | To be determind                      | L-019         |                        |                                                                                      |                                                                                                 |
| Others                | Hardware accelerator                                                                                                | To be determind                      |               |                        |                                                                                      |                                                                                                 |
|                       | Get log of cluster creation after cluster created                                                                   |                                      | C-002         |                        |                                                                                      |                                                                                                 |
|                       | Get log of cluster creation while create procedure                                                                  |                                      | C-002         |                        |                                                                                      |                                                                                                 |
|                       | Ntp func when create cluster                                                                                        |                                      | C-001         |                        |                                                                                      | [detail](./testcase_details/create_cluster.md#ntp-func-when-create-cluster)                     |
|                       | Ntp func while cluster in use                                                                                       |                                      | C-001         |                        |                                                                                      | [detail](./testcase_details/cluster_operation.md#readd-a-worker-to-cluster)                     |