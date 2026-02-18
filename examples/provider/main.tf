terraform {
  required_providers {
    manta = {
      source  = "gagno/manta"
      version = "0.0.1"
    }
    terraform = {
      source  = "builtin/terraform"
      version = ""
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

resource "terraform_data" "testing" {

}
