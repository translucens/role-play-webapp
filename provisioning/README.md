# Deploy on GCP by Terraform

```
git clone git@github.com:mittz/role-play-webapp.git
cd provisioning
export TF_VAR_project=${GOOGLE_CLOUD_PROJECT}
terraform init
terraform apply
```