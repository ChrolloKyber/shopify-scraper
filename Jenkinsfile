pipeline {
    agent any

    environment {
        GO_VERSION = '1.23.3'
        APP_NAME = 'shopify-scraper'
    }

    options {
        timestamps()
        ansiColor('xterm')
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Set Up Go') {
            steps {
                script {
                    // Set up Go environment if needed, for Jenkins with multiple Go versions
                    if (env.GOROOT == null) {
                        sh "wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
                        sh "sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz"
                        env.PATH = "/usr/local/go/bin:${env.PATH}"
                    }
                }
            }
        }

        stage('Dependencies') {
            steps {
                sh 'go mod download'
            }
        }

        stage('Build') {
            steps {
                sh 'go build -v -o bin/${APP_NAME} ./...'
            }
        }

        stage('Lint') {
            steps {
                sh 'go install golang.org/x/lint/golint@latest'
                sh 'golint ./... || true'
            }
        }

        stage('Test') {
            steps {
                sh 'go test -v ./...'
            }
        }

        stage('Docker Build') {
            when {
                expression { fileExists('Dockerfile') }
            }
            steps {
                script {
                    docker.build("${APP_NAME}:latest")
                }
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
