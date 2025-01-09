pipeline {
    agent any
    environment {
        // DOCKER_USERNAME = credentials('dockerhub-username') 
        // DOCKER_ACCESS_TOKEN = credentials('dockerhub-token') 
        GOPROXY = 'http://localhost:8081/repository/go-proxy'
        NPM_REGISTRY = 'http://localhost:8081/repository/npm-proxy/'
        DOCKER_REGISTRY = 'http://localhost:5001' // Nexus Docker group // Updated to npm-proxy
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
                            sh 'go mod tidy'
                            sh 'go test ./... -v'
                        }
                    }
                }
                stage('Frontend Unit Tests') {
                    steps {
                        dir('frontend') {
                            withEnv(["npm_config_registry=${NPM_REGISTRY}"]) {
                                sh 'npm install'
                                sh 'npm test'
                            }
                        }
                    }
                }
            }
        }

        stage('Run E2E Tests') {
            steps {
                withEnv(["npm_config_registry=${NPM_REGISTRY}", "XDG_RUNTIME_DIR=/tmp"]) {
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
                withCredentials([usernamePassword(credentialsId: 'nexus-docker-credentials', usernameVariable: 'NEXUS_USERNAME', passwordVariable: 'NEXUS_PASSWORD')]) {
                    sh '''
                        echo "$NEXUS_PASSWORD" | docker login $DOCKER_REGISTRY -u $NEXUS_USERNAME --password-stdin
                        
                        docker tag nextapp:1.0.0 $DOCKER_REGISTRY/nextapp:1.0.0
                        docker push $DOCKER_REGISTRY/nextapp:1.0.0

                        docker tag goapp:1.0.0 $DOCKER_REGISTRY/goapp:1.0.0
                        docker push $DOCKER_REGISTRY/goapp:1.0.0
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
