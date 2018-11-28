import org.nalej.SlackHelper
def slackHelper = new SlackHelper()
def packageName = "system-model"
def packagePath = "src/github.com/nalej/${packageName}"

pipeline {
    agent { node { label 'golang' } }
    options {
        checkoutToSubdirectory("${packagePath}")
    }

    stages {
        stage("Variable initialization") {
            steps {
                script {
                    dir("${packagePath}") {
                        env.GOPATH = env.WORKSPACE
                        env.remoteUrl = sh(returnStdout: true, script: "set +ex && git remote get-url origin").trim()
                        env.repoName = (env.remoteUrl =~ /https:\/\/github.com\/([^\n\r.]*).git/)[ 0 ][ 1 ]
                        env.commitId = sh(returnStdout: true, script: "set +ex && git log --pretty=format:'%H' -n 1").trim()
                        env.authorName = sh(returnStdout: true, script: "set +ex && git log --pretty=format:'%aN' -n 1").trim()
                        env.authorEmail = sh(returnStdout: true, script: "set +ex && git log --pretty=format:'%aE' -n 1").trim()
                        env.commitMsg = sh(returnStdout: true, script: "set +ex && git log --pretty=format:'%s' -n 1").trim()
                    }
                }
                script {
                    def timestamp = currentBuild.startTimeInMillis.intdiv(1000)
                    def attachment = slackHelper.createSlackAttachment("started", "", env.repoName, env.BRANCH_NAME, env.commitId, env.authorName, env.authorEmail, env.commitMsg, env.BUILD_URL, env.BUILD_NUMBER, timestamp)
                    slackSend attachments: attachment, message: ""
                }
            }
        }
        stage("Git setup") {
            steps {
                container("golang") {
                    script {
                        sh(script: """
                        set +ex && \
                        mkdir -p \$HOME/.ssh && \
                        cp /data/git-creds/id_rsa* \$HOME/.ssh/ && \
                        chmod 400 \$HOME/.ssh/id_rsa* && \
                        eval \"\$(ssh-agent -s)\" && \
                        ssh-add \$HOME/.ssh/id_rsa && \
                        ssh-keyscan -t rsa github.com >> \$HOME/.ssh/known_hosts && \
                        git config --global url."git@github.com:".insteadOf "https://github.com/"
                        """)
                    }
                }
            }
        }
        stage("Unit tests") {
            steps {
                container("golang") {
                    dir("${packagePath}") {
                        sh "dep ensure -v"
                        testResults = sh(returnStdout: true, script: "make test").trim()
                        if (env.CHANGE_ID) {
                            for (comment in pullRequest.comments) {
                                if (comment.user == "nalej-jarvis") {
                                    comment.delete()
                                }
                            }
                            commentContent = """
                            J.A.R.V.I.S. CI Test results:

                            ${testResults}
                            """
                            pullRequest.comment(commentContent)
                        }
                    }
                }
            }
        }
        stage("Binary compilation") {
            steps {
                container("golang") {
                    dir("${packagePath}") {
                        sh "make build"
                    }
                }
            }
        }
    }
    post {
        success {
            script {
                def timestamp = currentBuild.startTimeInMillis.intdiv(1000)
                def attachment = slackHelper.createSlackAttachment("success", "good", env.repoName, env.BRANCH_NAME, env.commitId, env.authorName, env.authorEmail, env.commitMsg, env.BUILD_URL, env.BUILD_NUMBER, timestamp)
                slackSend attachments: attachment, message: ""
                if (env.CHANGE_ID) {
                    pullRequest.removeLabel("J.A.R.V.I.S. is happy")
                    pullRequest.removeLabel("J.A.R.V.I.S. is sad")
                    pullRequest.addLabel("J.A.R.V.I.S. is happy")
                }
            }
        }
        failure {
            script {
                def timestamp = currentBuild.startTimeInMillis.intdiv(1000)
                def attachment = slackHelper.createSlackAttachment("failure", "danger", env.repoName, env.BRANCH_NAME, env.commitId, env.authorName, env.authorEmail, env.commitMsg, env.BUILD_URL, env.BUILD_NUMBER, timestamp)
                slackSend attachments: attachment, message: ""
                if (env.CHANGE_ID) {
                    pullRequest.removeLabel("J.A.R.V.I.S. is happy")
                    pullRequest.removeLabel("J.A.R.V.I.S. is sad")
                    pullRequest.addLabel("J.A.R.V.I.S. is sad")
                }
            }
        }
        aborted {
            script {
                def timestamp = currentBuild.startTimeInMillis.intdiv(1000)
                def attachment = slackHelper.createSlackAttachment("aborted", "warning", env.repoName, env.BRANCH_NAME, env.commitId, env.authorName, env.authorEmail, env.commitMsg, env.BUILD_URL, env.BUILD_NUMBER, timestamp)
                slackSend attachments: attachment, message: ""
                if (env.CHANGE_ID) {
                    pullRequest.removeLabel("J.A.R.V.I.S. is happy")
                    pullRequest.removeLabel("J.A.R.V.I.S. is sad")
                    pullRequest.addLabel("J.A.R.V.I.S. is sad")
                }
            }
        }
    }
}
