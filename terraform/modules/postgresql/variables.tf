variable "resource_group_name" {
  description = "The name of the resource group"
  type        = string
}

variable "location" {
  description = "The Azure region"
  type        = string
}

variable "project_name" {
  description = "The project name"
  type        = string
}

variable "environment" {
  description = "The environment name"
  type        = string
}

variable "admin_username" {
  description = "PostgreSQL admin username"
  type        = string
  default     = "postgresadmin"
}

variable "admin_password" {
  description = "PostgreSQL admin password"
  type        = string
  sensitive   = true
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}