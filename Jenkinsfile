pipeline {
    agent any

    environment {
        // Define environment variables for Docker image name and credentials
        DOCKER_IMAGE_NAME = "gotaskmanager"
        DOCKER_HUB_USERNAME = "markbosire" // Replace with your Docker Hub username
        DOCKER_HUB_CREDENTIALS = "dockerhub-credentials" // Replace with your Jenkins credentials ID
    }

    stages {
        stage('Test') {
            agent {
                label 'agent'
            }
            steps {
                sh '/usr/local/go/bin/go test'
            }
        }

        stage('Build Docker Image') {
            agent {
                label 'agent'
            }
            steps {
                script {
                    // Build the Docker image using the environment variables
                    dockerImage = docker.build("${env.DOCKER_HUB_USERNAME}/${env.DOCKER_IMAGE_NAME}:${env.BUILD_ID}")
                }
            }
        }
        stage('Push Docker Image') {
            steps {
                script {
                    // Log in to Docker Hub
                    withCredentials([usernamePassword(credentialsId: "${env.DOCKER_HUB_CREDENTIALS}", usernameVariable: 'DOCKER_USERNAME', passwordVariable: 'DOCKER_PASSWORD')]) {
                        sh """
                            docker login -u ${DOCKER_USERNAME} -p ${DOCKER_PASSWORD}
                        """
                    }

                    // Tag the image with latest and build number
                    sh """
                        docker tag ${env.DOCKER_HUB_USERNAME}/${env.DOCKER_IMAGE_NAME}:${env.BUILD_ID} ${env.DOCKER_HUB_USERNAME}/${env.DOCKER_IMAGE_NAME}:latest
                    """

                    // Push both tagged versions
                    sh """
                        docker push ${env.DOCKER_HUB_USERNAME}/${env.DOCKER_IMAGE_NAME}:${env.BUILD_ID}
                        docker push ${env.DOCKER_HUB_USERNAME}/${env.DOCKER_IMAGE_NAME}:latest
                    """

                    // Optional: Clean up local images to save space
                    sh """
                        docker image prune -f
                    """
                }
            }
        }

       stage('Run Docker Container') {
            steps {
                script {
                    // Create a named volume if it doesn't exist
                    sh """
                        docker volume create ${env.DOCKER_IMAGE_NAME}-data || true
                    """

                    // Stop and remove any running or stopped containers with the same name
                    sh '''
                        for container in $(docker ps -q --filter "name=${env.DOCKER_IMAGE_NAME}"); do
                            docker stop "$container"
                        done

                        for container in $(docker ps -a -q --filter "name=${env.DOCKER_IMAGE_NAME}"); do
                            docker rm "$container"
                        done
                    '''

                    // Pull the latest image from the Docker repository
                    sh "docker pull ${env.DOCKER_HUB_USERNAME}/${env.DOCKER_IMAGE_NAME}:${env.BUILD_ID}"

                    // Run the new container with persistent volume
                    sh """
                        docker run -d \
                            -p 8081:8080 \
                            --name ${env.DOCKER_IMAGE_NAME} \
                            -v ${env.DOCKER_IMAGE_NAME}-data:/app/data \
                            ${env.DOCKER_HUB_USERNAME}/${env.DOCKER_IMAGE_NAME}:${env.BUILD_ID}
                    """
                }
            }
        }
    }
}