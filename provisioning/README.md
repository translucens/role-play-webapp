# Deploy on GCP by Terraform

```
git clone https://github.com/mittz/role-play-webapp.git
cd role-play-webapp/provisioning
terraform init
terraform apply -var="project_id=${GOOGLE_CLOUD_PROJECT}"
```