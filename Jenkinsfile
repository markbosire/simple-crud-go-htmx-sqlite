pipeline {
    agent any

    environment {
        // Define environment variables for Docker image name and credentials
        DOCKER_IMAGE_NAME = "gotaskmanager"
        DOCKER_HUB_USERNAME = "markbosire" // Replace with your Docker Hub username
        DOCKER_HUB_CREDENTIALS = "docker-hub-credentials" // Replace with your Jenkins credentials ID
    }

    stages {
        stage('Test') {
            agent {
                label 'agent'
            }
            steps {
                sh 'go test'
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

       stage('Run Docker Container') {
            steps {
                script {
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

                    // Run the new container
                    sh """
                        docker run -d \
                            -p 8081:8080 \
                            --name ${env.DOCKER_IMAGE_NAME} \
                            ${env.DOCKER_HUB_USERNAME}/${env.DOCKER_IMAGE_NAME}:${env.BUILD_ID}
                    """
                }
            }
        }
    }
}