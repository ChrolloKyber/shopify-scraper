pipeline {
    agent any

    environment {
        GO_VERSION = '1.23.3'
        APP_NAME = 'shopify-scraper'
    }

    options {
        timestamps()
        buildDiscarder(logRotator(numToKeepStr: '30', artifactNumToKeepStr: '5'))
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
            junit '**/TEST-*.xml'
            cleanWs()
        }
        failure {
            echo "Failed to build"
        }
    }
}
