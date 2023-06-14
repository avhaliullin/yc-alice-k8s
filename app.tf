locals {
  app-env-vars = {
    K8S_HOST = var.k8s-host
    K8S_CA   = var.k8s-ca
  }
  bucket = "${var.folder-id}-fn-deploy"
}

resource "yandex_function" "alice-api" {
  entrypoint = "app.AliceHandler"
  memory     = 128
  name       = "k8s-alice"
  runtime    = "golang119"
  user_hash  = data.archive_file.app-code.output_base64sha256
  package {
    bucket_name = yandex_storage_object.fn-sources.bucket
    object_name = yandex_storage_object.fn-sources.key
  }
  environment        = local.app-env-vars
  service_account_id = yandex_iam_service_account.app-sa.id
  execution_timeout  = "3"
}

resource "yandex_iam_service_account" "app-sa" {
  name = "alice-app"
}

data "archive_file" "app-code" {
  output_path = "${path.module}/dist/app-code.zip"
  type        = "zip"
  source_dir  = "${path.module}/build"
}

resource "yandex_resourcemanager_folder_iam_binding" "app-k8s-access" {
  folder_id = var.folder-id
  members   = [
    "serviceAccount:${yandex_iam_service_account.app-sa.id}"
  ]
  role = "k8s.cluster-api.editor"
}

resource "yandex_iam_service_account" "func-deployer" {
  folder_id = var.folder-id
  name      = "func-deployer"
}

resource "yandex_resourcemanager_folder_iam_binding" "deployer-write-s3" {
  members = [
    "serviceAccount:${yandex_iam_service_account.func-deployer.id}"
  ]
  role      = "storage.editor"
  folder_id = var.folder-id
}

resource "yandex_iam_service_account_static_access_key" "deploy-fn" {
  service_account_id = yandex_iam_service_account.func-deployer.id
}

#resource "yandex_storage_bucket" "deploy-bucket" {
#  bucket     = local.bucket
#  access_key = yandex_iam_service_account_static_access_key.deploy-fn.access_key
#  secret_key = yandex_iam_service_account_static_access_key.deploy-fn.secret_key
#  depends_on = [yandex_resourcemanager_folder_iam_binding.deployer-write-s3]
#}

resource "yandex_storage_object" "fn-sources" {
  bucket         = local.bucket
  key            = "function.zip"
  access_key     = yandex_iam_service_account_static_access_key.deploy-fn.access_key
  secret_key     = yandex_iam_service_account_static_access_key.deploy-fn.secret_key
  content_base64 = filebase64(data.archive_file.app-code.output_path)
  #  source = data.archive_file.app-code.output_path
}

# output

output "function-alice-id" {
  value = yandex_function.alice-api.id
}

# configuration
terraform {
  required_providers {
    yandex = {
      source = "yandex-cloud/yandex"
    }
  }
}

provider "yandex" {
  folder_id = var.folder-id
  token     = var.yc-token
  //  version   = "0.53"
}

variable "folder-id" {
  type = string
}

variable "yc-token" {
  type = string
}

variable "k8s-host" {
  type = string
}

variable "k8s-ca" {
  type = string
}
