pipeline {
    agent any

    environment {
        GO_VERSION = '1.23.3'
        APP_NAME = 'shopify-scraper'
    }

    options {
        timestamps()
        buildDiscarder(logRotator(numToKeepStr: '5', artifactNumToKeepStr: '5'))
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Dependencies') {
            steps {
                sh 'go mod download'
            }
        }

        stage('Build') {
            steps {
                sh 'go build .'
            }
        }
    }

    post {
        always {
            cleanWs()
        }
        failure {
            echo "Failed to build"
        }
    }
}
