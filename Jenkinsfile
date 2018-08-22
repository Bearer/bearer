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
        stage("build") {
            steps {
                sh ".jenkins/build.sh"
            }
        }
        stage('test') {
            steps {
                container("node") {
                    ansiColor('xterm') {
                        sh "ls -l"
                    }
                }
            }
        }
        stage("deploy") {
            steps {
                echo "Deploy the package"
            }
        }
    }
}
