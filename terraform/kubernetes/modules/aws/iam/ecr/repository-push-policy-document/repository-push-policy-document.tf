###############################################################################
# Variables

variable "aws_account_id" {}
variable "region" {}
variable "repository" {}

###############################################################################
# Output

output "json" {
  value = "${data.aws_iam_policy_document.ecr_pull.json}"
}

###############################################################################
# IAM

data "aws_iam_policy_document" "ecr_pull" {
  statement {
    effect = "Allow"

    actions = [
      "ecr:PutImage",
      "ecr:InitiateLayerUpload",
      "ecr:UploadLayerPart",
      "ecr:CompleteLayerUpload",
    ]

    resources = [
      "arn:aws:ecr:${var.region}:${var.aws_account_id}:repository/${var.repository}",
    ]
  }
}
