terraform {
  required_version = ">= 1.0"

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 2.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.0"
    }
  }

  # We'll add remote backend later
  # backend "azurerm" {}
}

provider "azurerm" {
  features {}
}

# Get current client configuration
data "azurerm_client_config" "current" {}

# Create a resource group
resource "azurerm_resource_group" "main" {
  name     = "rg-${var.project_name}-${var.environment}"
  location = var.location
  tags     = var.tags
}

# Create ACR
module "acr" {
  source              = "./modules/acr"
  resource_group_name = azurerm_resource_group.main.name
  location            = azurerm_resource_group.main.location
  project_name        = var.project_name
  environment         = var.environment
  tags                = var.tags
}

# Create AKS
module "aks" {
  source              = "./modules/aks"
  resource_group_name = azurerm_resource_group.main.name
  location            = azurerm_resource_group.main.location
  project_name        = var.project_name
  environment         = var.environment
  node_count          = var.aks_node_count
  node_size           = var.aks_node_size
  acr_id              = module.acr.acr_id
  tags                = var.tags
}

# Create PostgreSQL
module "postgresql" {
  source              = "./modules/postgresql"
  resource_group_name = azurerm_resource_group.main.name
  location            = azurerm_resource_group.main.location
  project_name        = var.project_name
  environment         = var.environment
  admin_username      = var.postgres_admin_username
  admin_password      = var.postgres_admin_password
  tags                = var.tags
}