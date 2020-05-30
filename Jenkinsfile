pipeline{
    agent{
        label "node"
    }
    stages{
        stage("Lint"){
            steps{
                golint '-set_exit_status main/krkstops.go'
                staticcheck 'main/krkstops.go'
            }
        }
        stage("Test"){
            steps{
                go 'test ./krkstops'
            }
        }
        stage("Build"){
            go build main/krkstops.go
        }
    }
}