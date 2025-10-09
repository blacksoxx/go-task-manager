output "kube_config_raw" {
  description = "Raw Kubernetes config"
  value       = azurerm_kubernetes_cluster.aks.kube_config_raw
  sensitive   = true
}

output "host" {
  description = "Kubernetes cluster host"
  value       = azurerm_kubernetes_cluster.aks.kube_config[0].host
}

output "cluster_id" {
  description = "AKS cluster ID"
  value       = azurerm_kubernetes_cluster.aks.id
}

output "cluster_name" {
  description = "AKS cluster name"
  value       = azurerm_kubernetes_cluster.aks.name
}