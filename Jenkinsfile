pipeline {
    agent any

    parameters {
        booleanParam(defaultValue: true, description: 'Update frontend component?', name: 'UPDATE_FRONTEND')
        booleanParam(defaultValue: true, description: 'Update backend component?', name: 'UPDATE_BACKEND')
        string(defaultValue: '1.0.1', description: 'Frontend version to deploy', name: 'FRONTEND_VERSION')
        string(defaultValue: '1.0.1', description: 'Backend version to deploy', name: 'BACKEND_VERSION')
    }

    environment {
        DOCKER_USERNAME = credentials('dockerhub-username') 
        DOCKER_ACCESS_TOKEN = credentials('dockerhub-token')
        GITHUB_API_TOKEN = credentials('github-api-token') // mainly used for creating PRs via the GitHub API

        // GITHUB_CREDENTIALS = credentials('github-credentials') // if using github webhooks there's no need for that. if not - uncomment it!
        
        APP_REPO = "TankEngine-ish/property_management_system_full_stack_app"
        INFRA_REPO = "TankEngine-ish/property_management_system_infrastructure"

        DOCKER_ACCOUNT = "tankengine"

        PR_BRANCH_NAME = getBranchName()

        // GOPROXY = 'http://localhost:8081/repository/go-proxy'
        // NPM_REGISTRY = 'http://localhost:8081/repository/npm-proxy/'
        // DOCKER_HOSTED = 'localhost:5002' 
    }
    
    stages {
        stage('Checkout Code') {
            steps {
                git branch: 'feature', 
                    url: "https://github.com/${APP_REPO}"
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
                        // withEnv(["npm_config_registry=${NPM_REGISTRY}"]) { // commented out because I am not using the nexus registry at the moment.
                            sh 'npm install'
                            sh 'npm test'
                        }
                    }
                }
            }
        }

        stage('Run E2E Tests') {
            steps {
                // withEnv(["npm_config_registry=${NPM_REGISTRY}", "XDG_RUNTIME_DIR=/tmp"]) { // this is the old env for Nexus
                withEnv(["XDG_RUNTIME_DIR=/tmp"]) {
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
                    sh """
                        echo "\${DOCKER_ACCESS_TOKEN}" | docker login -u ${DOCKER_ACCOUNT} --password-stdin
                    """
                    
                    script {
                        if (params.UPDATE_FRONTEND) {
                            sh """
                                docker tag nextapp:${params.FRONTEND_VERSION} ${DOCKER_ACCOUNT}/nextapp:${params.FRONTEND_VERSION}
                                docker push ${DOCKER_ACCOUNT}/nextapp:${params.FRONTEND_VERSION}
                            """
                        }
                        
                        if (params.UPDATE_BACKEND) {
                            sh """
                                docker tag goapp:${params.BACKEND_VERSION} ${DOCKER_ACCOUNT}/goapp:${params.BACKEND_VERSION}
                                docker push ${DOCKER_ACCOUNT}/goapp:${params.BACKEND_VERSION}
                            """
                        }
                    }
                }
            }
        }

        stage('Create Infrastructure PR') {
            steps {
                script {
                    sh """
                        git clone https://github.com/${INFRA_REPO}.git
                        cd property_management_system_infrastructure
                        
                        # Create a new branch for the PR
                        git checkout -b ${env.PR_BRANCH_NAME}
                    """
                    
                    def filesToAdd = []
                    def updateMessage = "Update application versions:"
                    
                    if (params.UPDATE_BACKEND) {
                        sh """
                            cd property_management_system_infrastructure
                            
                            # Update backend version in values.yaml
                            sed -i "s/tag: \\"[0-9]\\.[0-9]\\.[0-9]*\\"/tag: \\"${params.BACKEND_VERSION}\\"/g" property-app-charts/charts/backend/values.yaml
                            
                            # Update backend version in Chart.yaml
                            sed -i "s/appVersion: \\"[0-9]\\.[0-9]\\.[0-9]*\\"/appVersion: \\"${params.BACKEND_VERSION}\\"/g" property-app-charts/charts/backend/Chart.yaml
                        """
                        filesToAdd.add("property-app-charts/charts/backend/values.yaml")
                        filesToAdd.add("property-app-charts/charts/backend/Chart.yaml")
                        updateMessage += "\\n- Backend: ${params.BACKEND_VERSION}"
                    }
                    
                    def filesToAddString = filesToAdd.join(" ")
                    def prTitle = getPrTitle()
                    def prBody = updateMessage + "\\n\\nAutomated Jenkins build: ${BUILD_URL}"

                    sh """
                        cd property_management_system_infrastructure
                        
                        # Configure git
                        git config user.email "jenkins@example.com"
                        git config user.name "Jenkins CI"
                        
                        # Add and commit changes
                        git add ${filesToAddString}
                        git commit -m "${updateMessage}"
                        
                        # Push the branch (using personal access token for authentication)
                        git push https://x-access-token:\${GITHUB_API_TOKEN}@github.com/${INFRA_REPO}.git ${env.PR_BRANCH_NAME}
                        
                        # Create PR using GitHub API
                        curl -X POST \\
                            -H "Authorization: token \${GITHUB_API_TOKEN}" \\
                            -H "Accept: application/vnd.github.v3+json" \\
                            https://api.github.com/repos/${INFRA_REPO}/pulls \\
                            -d '{
                            "title": "${prTitle}",
                            "body": "${prBody}",
                            "head": "${env.PR_BRANCH_NAME}",
                            "base": "prod"
                            }'
                    """
                }
            }
        }

        // The code below is for pushing to a nexus hosted repo, but I am using docker hub at the moment/

        // stage('Push Docker Images') {
        //     steps {
        //         withCredentials([usernamePassword(credentialsId: 'nexus-docker-credentials', usernameVariable: 'NEXUS_USERNAME', passwordVariable: 'NEXUS_PASSWORD')]) {
        //             sh '''
        //                 echo "$NEXUS_PASSWORD" | docker login $DOCKER_HOSTED -u $NEXUS_USERNAME --password-stdin
                        
        //                 docker tag nextapp:1.0.2 $DOCKER_HOSTED/nextapp:1.0.2
        //                 docker push $DOCKER_HOSTED/nextapp:1.0.2

        //                 docker tag goapp:1.0.2 $DOCKER_HOSTED/goapp:1.0.2
        //                 docker push $DOCKER_HOSTED/goapp:1.0.2
        //             '''
        //         }
        //     }
        // }

        stage('SonarQube Analysis') {
            steps {
                withSonarQubeEnv('SonarQube') { 
                    withCredentials([string(credentialsId: 'sonarqube-auth-token', variable: 'SONAR_TOKEN')]) {
                        sh '''
                            /opt/sonar-scanner/bin/sonar-scanner \\
                                -Dsonar.projectKey=property_management_system \\
                                -Dsonar.sources=backend,frontend \\
                                -Dsonar.host.url=http://localhost:9000 \\
                                -Dsonar.login=$SONAR_TOKEN
                        '''
                    }
                }
            }
        }
    }

    post {
        always {
            echo 'Pipeline completed!'
            cleanWs()
        }
        success {
            echo 'Pipeline succeeded! A pull request has been created to update the infrastructure.'
        }
        failure {
            echo 'Pipeline is a doozy! Consult your seniors for help.'
        }
    }
}

def getBranchName() {
    def components = []
    if (params.UPDATE_FRONTEND) components.add("frontend-${params.FRONTEND_VERSION}")
    if (params.UPDATE_BACKEND) components.add("backend-${params.BACKEND_VERSION}")
    
    return "update-" + components.join("-")
}

def getPrTitle() {
    def components = []
    if (params.UPDATE_BACKEND) components.add("Backend ${params.BACKEND_VERSION}")
    if (params.UPDATE_FRONTEND) components.add("Frontend ${params.FRONTEND_VERSION}")
    
    return "Update image versions: " + components.join(", ")
}