pipeline {
	agent { docker { image 'golang:1.10' } }

	environment {
		GOPATH = $WORKSPACE
	}

    stages {
        stage('build') {
            steps {
                sh 'go version'
            }
        }
		stage('test') {
			steps {
				// sh 'mkdir -p ${GOPATH}/src/github.com/briansimoni/stereodose'
				// sh 'ln -sf ${WORKSPACE} ${GOPATH}/src/github.com/briansimoni/stereodose'
				sh 'go test ./...'
			}
		}
    }
}
