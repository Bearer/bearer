def label() {
    def value = "bearercli-${UUID.randomUUID().toString()}"
    return value
}

pipeline {
    agent {
        kubernetes {
        label label()
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
        buildDiscarder(logRotator(numToKeepStr: '50'))
    }

    parameters {
        string(defaultValue: '', description: 'TAG which will be used by lerna', name: 'LERNA_TAG')
        booleanParam(defaultValue: false, description: 'Lerna force publish', name: 'FORCE')
        booleanParam(defaultValue: false, description: 'Canary release', name: 'CANARY')

    }

    stages {
        stage('checkout') {
            steps {
                checkout scm
            }
        }
        stage("Build") {
            steps {
                container("bearercli") {
                    ansiColor('xterm') {
                        sh ".jenkins/prepare.sh"
                        sh "yarn install --frozen-lockfile"
                     }
                }
            }
        }
        stage('Test') {
             when {
                expression { params.CANARY  == false }
            }

            steps {
                container("bearercli") {
                    ansiColor('xterm') {
                        sh ".jenkins/test.sh"
                    }
                }
            }
        }
        stage("Deploy") {
            when {
                expression { params.LERNA_TAG  != '' }
                branch 'master'
            }

            steps {
                container("bearercli") {
                    ansiColor('xterm') {
                        sh(".jenkins/deploy.sh")
                    }
                }
            }
        }

        stage("Deploy Canary") {
            when {
                expression { params.CANARY  == true }
            }

            steps {
                container("bearercli") {
                    ansiColor('xterm') {
                        sh(".jenkins/deploy.sh")
                    }
                }
            }
        }
    }
}
