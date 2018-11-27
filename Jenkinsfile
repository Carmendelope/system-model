import org.nalej.SlackHelper

pipeline {
    agent { node { label 'golang' } }

    stages {
        stage("Initialization") {
            steps {
                script {
                    env.remoteUrl = sh(returnStdout: true, script: "git remote get-url origin").trim()
                    env.repoName = (env.remoteUrl =~ /https:\/\/github.com\/([^\n\r.]*).git/)[ 0 ][ 1 ]
                    env.commitId = sh(returnStdout: true, script: "git log --pretty=format:'%H' -n 1").trim()
                    env.authorName = sh(returnStdout: true, script: "git log --pretty=format:'%aN' -n 1").trim()
                    env.authorEmail = sh(returnStdout: true, script: "git log --pretty=format:'%aE' -n 1").trim()
                    env.commitMsg = sh(returnStdout: true, script: "git log --pretty=format:'%s' -n 1").trim()
                    def slackHelper = new SlackHelper()
                    def timestamp = currentBuild.startTimeInMillis.intdiv(1000)
                    def attachment = slackHelper.createSlackAttachment("started", "", env.repoName, env.BRANCH_NAME, env.commitId, env.authorName, env.authorEmail, env.commitMsg, env.BUILD_URL, env.BUILD_NUMBER, timestamp)
                    slackSend channel: "@rnunez", attachments: attachment, message: ""
                }
            }
        }
        stage("Unit tests") {
            steps {
                container("golang") {
                    script {
                        sh "make test"
                    }
                }
            }
        }
        stage("Binary compilation") {
            steps {
                container("golang") {
                    script {
                        sh "make build"
                    }
                }
            }
        }
    }
    post {
        success {
            script {
                def slackHelper = new SlackHelper()
                def timestamp = currentBuild.startTimeInMillis.intdiv(1000)
                def attachment = slackHelper.createSlackAttachment("success", "good", env.repoName, env.BRANCH_NAME, env.commitId, env.authorName, env.authorEmail, env.commitMsg, env.BUILD_URL, env.BUILD_NUMBER, timestamp)
                slackSend channel: "@rnunez", attachments: attachment, message: ""
            }
        }
        failure {
            script {
                def slackHelper = new SlackHelper()
                def timestamp = currentBuild.startTimeInMillis.intdiv(1000)
                def attachment = slackHelper.createSlackAttachment("failure", "danger", env.repoName, env.BRANCH_NAME, env.commitId, env.authorName, env.authorEmail, env.commitMsg, env.BUILD_URL, env.BUILD_NUMBER)
                slackSend channel: "@rnunez", attachments: attachment, message: ""
            }
        }
        aborted {
            script {
                def slackHelper = new SlackHelper()
                def timestamp = currentBuild.startTimeInMillis.intdiv(1000)
                def attachment = slackHelper.createSlackAttachment("aborted", "warning", env.repoName, env.BRANCH_NAME, env.commitId, env.authorName, env.authorEmail, env.commitMsg, env.BUILD_URL, env.BUILD_NUMBER)
                slackSend channel: "@rnunez", attachments: attachment, message: ""
            }
        }
    }
}
