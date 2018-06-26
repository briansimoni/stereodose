pipeline {
	agent { docker { image 'golang:1.10' } }

	environment {
        GOPATH = 'true'
        DB_ENGINE    = 'sqlite'
    }

    stages {
        stage('build') {
            steps {
                sh 'go version'
            }
        }
		stage('test') {
			steps {
				sh 'ln -sf ${WORKSPACE} ${GOPATH}/src/github.com/briansimoni/stereodose'
				sh 'go test ./...'
			}
		}
    }
}
