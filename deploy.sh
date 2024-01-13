#################################################
###  Boilerplate Go API Clean - Deployment Script ###
#################################################

# Log service deployment
echo 'Running Go API Clean Deployment...'

# Configure GCloud Project
echo 'Configing gcloud project...'
gcloud config set project <project-id-here>

# Build GCloud Docker Image
echo 'Build new gcloud image...'
gcloud builds submit --tag gcr.io/<project-id-here>/boilerplate-go-api-clean

# Deploy new Docker Image to Cloud Run
echo 'Deploying to gcloud run...'
gcloud run deploy boilerplate-go-api-clean --image gcr.io/<project-id-here>/boilerplate-go-api-clean --platform managed --region us-east1 --memory 2Gi --cpu 2 --allow-unauthenticated

