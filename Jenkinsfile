pipeline {
	agent { docker { image 'golang:1.10' } }

    stages {
        stage('build') {
            steps {
                sh 'go version'
            }
        }
		stage('test') {
			steps {
				sh 'export GOPATH=${WORKSPACE}'
				// sh 'mkdir -p ${GOPATH}/src/github.com/briansimoni/stereodose'
				// sh 'ln -sf ${WORKSPACE} ${GOPATH}/src/github.com/briansimoni/stereodose'
				sh 'go test ./...'
			}
		}
    }
}
