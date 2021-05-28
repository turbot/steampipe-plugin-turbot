connection "turbot" {
  plugin = "turbot"

  # Default credentials:
  # By default, Steampipe will use your Turbot profiles and credentials exactly
  # the same as the Turbot CLI and Turbot Terraform provider. In many cases, no
  # extra configuration is required to use Steampipe.

  # Use an existing Turbot profile configured in ~/.config/turbot
  # profile = "my-profile"

  # Define exact connection parameters to Turbot. This takes precedence over all
  # Turbot configuration, profile and environment variables.
  # workspace  = "https://turbot-acme.cloud.turbot.com/"
  # access_key = "c8e2c2ed-1ca8-429b-b369-010e3cf75aac"
  # secret_key = "a3d8385d-47f7-40c5-a90c-bfdf5b43c8dd"
}
