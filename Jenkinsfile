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
        JENKINS_PRIVATE_KEY = credentials("jenkins-github-ssh-private")
    }

    options {
        skipDefaultCheckout true
        buildDiscarder(logRotator(numToKeepStr: '5'))
    }

    parameters {
        string(defaultValue: '', description: 'TAG which will be used by lerna', name: 'LERNA_TAG')
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
            when {
                expression { params.LERNA_TAG  != '' }
                branch 'feature/jenkins'
            }

            steps {
                container("node") {
                    ansiColor('xterm') {
                        sh(".jenkins/deploy.sh")
                    }
                }
            }
        }
    }
}
