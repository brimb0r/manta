terraform {
  required_providers {
    manta = {
      source  = "gagno/manta"
      version = "0.0.1"
    }
  }
}

provider "manta" {
  endpoint = ""
}

output "racecar_is_palindrome" {
  value = provider::manta::is_palindrome("racecar")
}

output "hello_is_palindrome" {
  value = provider::manta::is_palindrome("hello")
}

output "version_compare" {
  value = provider::manta::semver_compare("1.2.3", "1.3.0")
}

output "merged_config" {
  value = jsondecode(provider::manta::deep_merge(
    jsonencode({ defaults = { timeout = 30, retries = 3 }, region = "us-east-1" }),
    jsonencode({ defaults = { timeout = 60 }, region = "eu-west-1", debug = true })
  ))
}

output "truncated_name" {
  value = provider::manta::truncate("my-very-long-resource-name-that-exceeds-the-limit", 24)
}

output "masked_key" {
  value = provider::manta::mask("sk-1234567890abcdef", 4)
}