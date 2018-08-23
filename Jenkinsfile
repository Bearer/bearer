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
                branch 'feature/jenkins'
                beforeAgent true
            }

            steps {
                script {
                    def returnValue = input(
                        message: "Wich Lerna tag do you want to use ?",
                        ok: "Validate",
                        parameters: [
                            string(
                                defaultValue: 'test',
                                description: 'The tag to use for Lerna',
                                name: 'LERNA_TAG'
                            )
                        ]
                    )
                 }

                container("node") {
                    ansiColor('xterm') {
                        sh(".jenkins/deploy.sh")
                    }
                }
            }
        }
    }
}
