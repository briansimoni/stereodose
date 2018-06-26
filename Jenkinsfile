pipeline {
    agent { dockerfile true }
    stages {
        stage('build') {
            steps {
                sh 'go version'
            }
        }
		stage('test') {
			steps {
				sh 'go test ./...'
			}
		}
    }
}
