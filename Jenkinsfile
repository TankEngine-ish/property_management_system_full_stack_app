// pipeline {
//     agent any
//     environment {
//         DOCKER_USERNAME = credentials('dockerhub-username') 
//         DOCKER_ACCESS_TOKEN = credentials('dockerhub-token') 
//         DOCKER_REGISTRY = 'docker.io'
//         // GOPROXY = 'http://localhost:8081/repository/go-proxy'
//         // NPM_REGISTRY = 'http://localhost:8081/repository/npm-proxy/'
//         // DOCKER_HOSTED = 'localhost:5002' // Hosted repository for private images, moved away from docker hub
//     }
//     stages {
//         stage('Checkout Code') {
//             steps {
//                 git branch: 'feature', url: 'https://github.com/TankEngine-ish/property_management_system_full_stack_app'
//             }
//         }

//         stage('Run Unit Tests') {
//             parallel {
//                 stage('Go Unit Tests') {
//                     steps {
//                         dir('backend') {
//                             sh 'go mod tidy'
//                             sh 'go test ./... -v'
//                         }
//                     }
//                 }
//                 stage('Frontend Unit Tests') {
//                     steps {
//                         dir('frontend') {
//                             // withEnv(["npm_config_registry=${NPM_REGISTRY}"]) {
//                                 sh 'npm install'
//                                 sh 'npm test'
//                             }
//                         }
//                     }
//                 }
//             }
//         }

//         stage('Run E2E Tests') {
//             steps {
//                 withEnv(["npm_config_registry=${NPM_REGISTRY}", "XDG_RUNTIME_DIR=/tmp"]) {
//                     sh '''
//                         npm install
//                         npx cypress run --config-file ./cypress.config.js --spec cypress/e2e/userExperience.cy.js
//                     '''
//                 }
//             }
//         }

//         stage('Build Docker Images') {
//             steps {
//                 sh 'docker compose build'
//             }
//         }


//         stage('Push Docker Images') {
//             steps {
//                 withCredentials([string(credentialsId: 'dockerhub-token', variable: 'DOCKER_ACCESS_TOKEN')]) {
//                     sh '''
//                         echo "$DOCKER_ACCESS_TOKEN" | docker login -u tankengine --password-stdin

//                         docker tag nextapp:1.0.2 tankengine/nextapp:1.0.2
//                         docker push tankengine/nextapp:1.0.2

//                         docker tag goapp:1.0.2 tankengine/goapp:1.0.2
//                         docker push tankengine/goapp:1.0.2
//                     '''
//                 }
//             }
//         }
//     }

//     // The above code is for pushing to docker hub, but I am using nexus as my docker registry now, so I changed the code to push to a nexus hosted repo - code below: //

//         // stage('Push Docker Images') {
//         //     steps {
//         //         withCredentials([usernamePassword(credentialsId: 'nexus-docker-credentials', usernameVariable: 'NEXUS_USERNAME', passwordVariable: 'NEXUS_PASSWORD')]) {
//         //             sh '''
//         //                 echo "$NEXUS_PASSWORD" | docker login $DOCKER_HOSTED -u $NEXUS_USERNAME --password-stdin
                        
//         //                 docker tag nextapp:1.0.2 $DOCKER_HOSTED/nextapp:1.0.2
//         //                 docker push $DOCKER_HOSTED/nextapp:1.0.2

//         //                 docker tag goapp:1.0.2 $DOCKER_HOSTED/goapp:1.0.2
//         //                 docker push $DOCKER_HOSTED/goapp:1.0.2
//         //             '''
//         //         }
//         //     }
//         // }

//         stage('SonarQube Analysis') {
//             steps {
//                 withSonarQubeEnv('SonarQube') { 
//                     withCredentials([string(credentialsId: 'sonarqube-auth-token', variable: 'SONAR_TOKEN')]) { // mapped this variable to the token's id in Jenkins
//                         sh '''
//                             /opt/sonar-scanner/bin/sonar-scanner \
//                                 -Dsonar.projectKey=property_management_system \
//                                 -Dsonar.sources=backend,frontend \
//                                 -Dsonar.host.url=http://localhost:9000 \
//                                 -Dsonar.login=$SONAR_TOKEN
//                         '''
//                     }
//                 }
//             }
//         }
//     }

//     post {
//         always {
//             echo 'Pipeline completed!'
//         }
//         success {
//             echo 'Pipeline succeeded!'
//         }
//         failure {
//             echo 'Pipeline failed!'
//         }
//     }
// }


pipeline {
    agent any
    environment {
        DOCKER_USERNAME = credentials('dockerhub-username') 
        DOCKER_ACCESS_TOKEN = credentials('dockerhub-token') 
        // GOPROXY = 'http://localhost:8081/repository/go-proxy'
        // NPM_REGISTRY = 'http://localhost:8081/repository/npm-proxy/'
        // DOCKER_HOSTED = 'localhost:5002' 
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
                    sh '''
                        echo "$DOCKER_ACCESS_TOKEN" | docker login -u tankengine --password-stdin

                        docker tag nextapp:1.0.2 tankengine/nextapp:1.0.2
                        docker push tankengine/nextapp:1.0.2

                        docker tag goapp:1.0.2 tankengine/goapp:1.0.2
                        docker push tankengine/goapp:1.0.2
                    '''
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
                            /opt/sonar-scanner/bin/sonar-scanner \
                                -Dsonar.projectKey=property_management_system \
                                -Dsonar.sources=backend,frontend \
                                -Dsonar.host.url=http://localhost:9000 \
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
        }
        success {
            echo 'Pipeline succeeded!'
        }
        failure {
            echo 'Pipeline failed!'
        }
    }
}