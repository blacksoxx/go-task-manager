output "acr_id" {
  description = "The ACR resource ID"
  value       = azurerm_container_registry.acr.id
}

output "acr_name" {
  description = "The ACR name"
  value       = azurerm_container_registry.acr.name
}

output "acr_login_server" {
  description = "The ACR login server URL"
  value       = azurerm_container_registry.acr.login_server
}