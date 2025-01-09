pipeline {
    agent any
    environment {
        DOCKER_USERNAME = credentials('dockerhub-username') 
        DOCKER_ACCESS_TOKEN = credentials('dockerhub-token') 
        GOPROXY = 'http://localhost:8081/repository/go-proxy' 
    }
    stages {
        stage('Checkout Code') {
            steps {
                git branch: 'feature', url: 'https://github.com/TankEngine-ish/property_management_system_full_stack_app'
            }
        }

        stage('Run Unit Tests') {
            parallel {
                stage('Go Unit Tests') {
                    steps {
                        dir('backend') {
                            sh 'go mod tidy' // E
                            sh 'go test ./... -v'
                        }
                    }
                }
                stage('Frontend Unit Tests') {
                    steps {
                        dir('frontend') {
                            sh 'npm install'
                            sh 'npm test'
                        }
                    }
                }
            }
        }

        stage('Run E2E Tests') {
            steps {
                withEnv(['XDG_RUNTIME_DIR=/tmp']) {
                    sh '''
                        npm install
                        npx cypress run --config-file ./cypress.config.js --spec cypress/e2e/userExperience.cy.js
                    '''
                }
            }
        }

        stage('Build Docker Images') {
            steps {
                sh 'docker compose build'
            }
        }

        stage('Push Docker Images') {
            steps {
                withCredentials([string(credentialsId: 'dockerhub-token', variable: 'DOCKER_ACCESS_TOKEN')]) {
                    sh '''
                        echo "$DOCKER_ACCESS_TOKEN" | docker login -u tankengine --password-stdin

                        docker tag nextapp:1.0.0 tankengine/nextapp:1.0.0
                        docker push tankengine/nextapp:1.0.0

                        docker tag goapp:1.0.0 tankengine/goapp:1.0.0
                        docker push tankengine/goapp:1.0.0
                    '''
                }
            }
        }
    }
    post {
        always {
            echo 'Pipeline completed!'
        }
        success {
            echo 'Pipeline succeeded!'
        }
        failure {
            echo 'Pipeline failed!'
        }
    }
}
