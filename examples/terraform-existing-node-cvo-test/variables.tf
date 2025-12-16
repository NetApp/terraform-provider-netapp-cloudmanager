variable "cloudmanager_refresh_token" {
  description = "BlueXP refresh token"
  type        = string
}

variable "cloudmanager_environment" {
  description = "BlueXP environment"
  type        = string
  default     = "stage"
}