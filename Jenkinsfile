/*
    This is the main pipeline section with the stages of the CI/CD
 */
pipeline {

    options {
        // Build auto timeout
        timeout(time: 60, unit: 'MINUTES')
    }

    // Some global default variables
    environment {
        IMAGE_NAME = "${JOB_NAME}"

    }

    //parameters {
        //string (name: 'GIT_BRANCH',           defaultValue: 'master',  description: 'Git branch to build')
        //credentials(name: 'DOCKER_REG', description: 'A user to build with', defaultValue: '', credentialType: "docker registry credentials ", required: true)
    //}

    // In this example, all is built and run from the master
    agent any
    // Pipeline stages
    stages {

        ////////// Step 1 //////////
        stage('Git Repo setup') {
            steps {
                checkout scm
            }
        }

        ////////// Step 2 //////////
        stage('Build  and Publish ') {

             steps {
                echo 'Starting to build docker image'
                sh 'cd src/frontend'
                //sh "docker build -t haseebh/images-analyzer:${BUILD_NUMBER} ."
                script{
                    def customImage = docker.build("haseebh/frontend:v1.${env.BUILD_ID}.0")
                    customImage.push()
                }
            }

        }
        ////////// Step 3 //////////
        stage('Deploy') {

             steps {
                echo 'Starting to build docker image'
                sh '''
                        curl \
                        --request POST \
                        -H "Content-Type: application/json" \
                        -H "token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.CnsKICAidGVhbXNfcm9sZXMiOnsiQmFja2VuZCBUZWFtIjoiVGVhbSBMZWFkIn0sCiAgImlzcyI6ImNsb3VkcGxleCIsCiAgImV4cCI6IjAiLAogICJ1c2VybmFtZSI6Imhhc2VlYkBjbG91ZHBsZXguaW8iLAogICJjb21wYW55SWQiOiI1ZGFhYzRkZjZhOWUyZjAwMWM3ZDgzYjAiLAogICAiY29tcGFueU5hbWUiOiJjbG91ZHBsZXgiLAogICJpc0FkbWluIjoiZmFsc2UiLAogICJ0b2tlbl90eXBlIjoiMSIsCiAgIm15cm9sZXMiOlsiU3VwZXItVXNlciIsIlRlYW0gTGVhZCJdCn0KICAgICAg.o0LR5YL31RfexdGCWZO0QYVJJ5iiX3-ot7pFYKEmVcw" \
                        --data '{"image":"haseebh/frontend","tag": "'v1.${BUILD_NUMBER}.0'" , "type":"container", "project_ids":["application-3j10ci"]}'\
                         https://app.cloudplex.io/rabbit/api/v1/cd/trigger/deployment

 	             '''

            }

        }



        stage('cleanup') {
            steps {
                sh "docker rmi haseebh/haseebh/frontend:v1.${BUILD_NUMBER}.0"
            }
        }
    }
}