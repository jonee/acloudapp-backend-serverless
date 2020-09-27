# acloudapp-backend-serverless

A serverless backend implementation for acloudapp- a full stack web and mobile project starter with flexible and changeable parts.  

Currently, there are codes for Golang that saves into MongoDB and deployed into an OpenWhisk provider like IBM Cloud and / or deployed into AWS. It is also a (backend) framework for supporting other programming languages, databases and cloud providers. 

# Backend Features
Amongst others:
- No technical dept as of right now (uses go mod, official MongoDB driver, other fresh bleeding edge parts)
- Internationalization ready, uses language keys throughout, framework for translation in mind from the beginning.
- Usual serverless features
    - Flexibility with programming language up to the end point level eg you could have different programming languages for each end point.
    - Scaling built into cloud functions (such as OpenWhisk actions and AWS Lambdas).
- Authentication implementation available (JWT)
- Opinionated code heirarchy arranged in order of programming language, database and cloud providers. So just by looking at the folders you should know which codes are for Go, which Go codes are for MongoDB database and finally which Go codes for MongoDB database is meant to be deployed to an OpenWhisk provider. 
- Support for unlimited stages eg dev and prod.

# Available end points 


# Future feature targets
Acloudapp is envisioned to be a test bed for new features. This list would be maintained somewhere else as the project grows but some nice to have features for the backend would be:
- Better documentation
- Unit tests for apis end points
- Process (or even encryption) for validation / integrity check of parameters to be sure that the calling client is authorized and known
- Api throttle / defense for brute force attacks (eg hydra)

# Documentation, installation, requirements
In brief:
- You should be able to run the serverless examples for Golang and OpenWhisk and / or Golang and AWS specially one that would interface to a MongoDB database. 
- In the codes you should find more than one configuration.go.sample (or configuration.js.sample for the nodejs examples). Save your credentials in the files eg MongoDB, AWS access and secret keys and rename the files without the sample extension. (AWS is also needed for AWS SES for email sending and very likely S3 for object storage in the future). 

# Contributing
1. Contributors attest that their codes are free of any legal issues and that the project is shielded from any ramifications brought by the code in question.
2. You donate and transfer code ownership to the project. 

# License
The codes are licensed under GPLv3. 

# Donate
If the codes in anyway help you out, please consider donating to support the development of the project. 

# Hire me / us
