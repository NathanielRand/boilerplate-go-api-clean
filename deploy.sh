#################################################
###  Remote Jobs Guru API - Deployment Script ###
#################################################

# Log service deployment
echo 'Running Go API Clean Deployment...'

# Configure GCloud Project
echo 'Configing gcloud project...'
gcloud config set project morebytes

# Build GCloud Docker Image
echo 'Build new gcloud image...'
gcloud builds submit --tag gcr.io/morebytes/boilerplate-go-api-clean

# Deploy new Docker Image to Cloud Run
echo 'Deploying to gcloud run...'
gcloud run deploy boilerplate-go-api-clean --image gcr.io/morebytes/boilerplate-go-api-clean --platform managed --region us-east1 --memory 2Gi --cpu 2 --allow-unauthenticated

