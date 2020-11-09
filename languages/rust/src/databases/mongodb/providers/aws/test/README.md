Hackish way to deploy

1. sls deploy would give an error about missing file path in the package path. 
2. If you look at .serverless it should have 3 files- 

cloudformation-template-create-stack.json
cloudformation-template-update-stack.json
serverless-state.json

3. rename .serverless folder to p
4. cargo build --release then add ACA-r-test.zip which consists of target/release/test only to the p folder
5. sls deploy --package p
