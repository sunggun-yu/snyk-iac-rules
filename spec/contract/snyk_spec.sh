#!/bin/bash

cleanup() { rm -rf ./tmp; rm bundle.tar.gz; }
AfterAll 'cleanup'

Describe 'Contract test'
    It 'verifies contract between the SDK and Snyk'
        snyk_iac_test() {
            # create tmp test directory for contract tests
            mkdir tmp
            
            # create a basic rule
            ./snyk-iac-custom-rules template ./tmp --rule Contract 

            # write rule and test
            cp -r ./fixtures/custom-rules/rules/CUSTOM-3/* ./tmp/rules/Contract

            # run tests and make sure they pass
            ./snyk-iac-custom-rules test ./tmp 

            # create bundle
            ./snyk-iac-custom-rules build ./tmp --ignore "testing" --ignore "*_test.rego" 

            # authenticate with Snyk
            snyk auth $SNYK_TOKEN 

            # use bundle with Snyk
            snyk iac test --rules=./bundle.tar.gz ./tmp/rules/Contract/fixtures/test.yaml
        }

        When call snyk_iac_test
        The status should eq 1
        The output should include "Generated template" # the rule was tempalted successfully
        The output should include "PASS: 1/1" # the tests passed
        The output should include "Bundle bundle.tar.gz has been generated" # the bundle has been generated
        The output should include "Test [Critical Severity] [CUSTOM-3]" # it should include the custom rule in its output
    End
End