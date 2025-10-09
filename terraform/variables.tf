# Azure Authentication Variables
variable "subscription_id" {
  description = "Azure Subscription ID"
  type        = string
  sensitive   = true
}

variable "client_id" {
  description = "Azure Client ID for Service Principal"
  type        = string
  sensitive   = true
}

variable "client_secret" {
  description = "Azure Client Secret for Service Principal"
  type        = string
  sensitive   = true
}

variable "tenant_id" {
  description = "Azure Tenant ID"
  type        = string
  sensitive   = true
}

# Project Configuration
variable "project_name" {
  description = "The project name"
  type        = string
  default     = "taskmanager"
}

variable "environment" {
  description = "The environment (dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "location" {
  description = "The Azure region where resources will be created"
  type        = string
  default     = "francecentral"
}

# AKS Configuration
variable "aks_node_count" {
  description = "Number of AKS nodes"
  type        = number
  default     = 2
}

variable "aks_node_size" {
  description = "VM size for AKS nodes"
  type        = string
  default     = "Standard_B2s"
}

# Tags
variable "tags" {
  description = "Tags to apply to all resources"
  type        = map(string)
  default = {
    Project     = "TaskManager"
    Environment = "dev"
    ManagedBy   = "Terraform"
  }
}


variable "postgres_admin_username" {
  description = "PostgreSQL admin username"
  type        = string
  default     = "postgresadmin"
}

variable "postgres_admin_password" {
  description = "PostgreSQL admin password"
  type        = string
  sensitive   = true
}