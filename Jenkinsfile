pipeline {
    agent { docker { image 'golang:1.10' } }
    stages {
        stage('build') {
            steps {
                sh 'go version'
            }
        }
    }
}
