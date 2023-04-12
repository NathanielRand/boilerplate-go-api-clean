#################################################
###  Remote Jobs Guru API - Deployment Script ###
#################################################

# Log service deployment
echo 'Running Image Converter API Database Deployment...'

# Configure GCloud Project
echo 'Configing gcloud project...'
gcloud config set project webchest

# Build GCloud Docker Image
echo 'Build new gcloud image...'
gcloud builds submit --tag gcr.io/webchest/webchest-image-converter-api

# Deploy new Docker Image to Cloud Run
echo 'Deploying to gcloud run...'
gcloud run deploy webchest-image-converter-api --image gcr.io/webchest/webchest-image-converter-api --platform managed --region us-east1 --memory '512Mi' --cpu '1' --allow-unauthenticated

