output "resource_group_name" {
  description = "The name of the resource group"
  value       = azurerm_resource_group.main.name
}

output "acr_name" {
  description = "The ACR name"
  value       = module.acr.acr_name
}

output "acr_login_server" {
  description = "The ACR login server URL"
  value       = module.acr.acr_login_server
}

output "aks_cluster_name" {
  description = "The AKS cluster name"
  value       = module.aks.cluster_name
}

output "postgresql_fqdn" {
  description = "PostgreSQL server FQDN"
  value       = module.postgresql.server_fqdn
}

output "postgresql_database_name" {
  description = "PostgreSQL database name"
  value       = module.postgresql.database_name
}