pipeline {
    agent {
        kubernetes {
        label 'node'
        defaultContainer 'jnlp'
        yamlFile '.jenkins/node.yml'
        }
    }

    environment { 
        AWS_ACCESS = credentials('aws-identity') 
        NPM_TOKEN =  credentials("npm-token")
    }

    options {
        skipDefaultCheckout true
        buildDiscarder(logRotator(numToKeepStr: '5'))
    }

    stages {
        stage('checkout') {
            steps {
                checkout scm
            }
        }
        stage("Build") {
            steps {
                container("node") {
                    ansiColor('xterm') {
                        sh ".jenkins/build.sh"
                     }
                }
            }
        }
        stage('Test') {
            steps {
                container("node") {
                    ansiColor('xterm') {
                        sh ".jenkins/test.sh"
                    }
                }
            }
        }
        stage("Deploy") {
            steps {
                echo "Deploy the package"
            }
        }
    }
}
