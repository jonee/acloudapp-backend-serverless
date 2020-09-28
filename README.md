# acloudapp-backend-serverless

A serverless backend implementation for ACloudApp- a full stack web and mobile project starter with flexible and changeable parts.  

Currently, there are codes for Golang that saves into MongoDB and deployed into an OpenWhisk provider like IBM Cloud and / or deployed into AWS. It is also a (backend) framework for supporting other programming languages, databases and cloud providers in a serverless implementation.

# Code repos
Codes should be available in
https://gitlab.com/jonee316/acloudapp-backend-serverless
https://github.com/jonee/acloudapp-backend-serverless

# Backend Features
Amongst others:
- Aimed to have no technical dept (uses go mod, official MongoDB driver, other fresh bleeding edge parts)
- Internationalization ready, uses language keys throughout, framework for translation (including email templates) in mind from the beginning.
- Serverless features
    - Flexibility with programming language up to the end point level eg you could have different programming languages for each end point.
    - Scaling built into cloud functions (such as OpenWhisk actions and AWS Lambdas).
- Authentication implementation available (JWT)
- Opinionated code heirarchy arranged in order of programming language, database and cloud providers. So just by looking at the folders you should know which codes are for Go, which Go codes are for MongoDB database and finally which Go codes for MongoDB database is meant to be deployed to an OpenWhisk provider. 
- Support for unlimited stages eg dev and prod.

# Available end points
1. Registration (which sends an email with a validation link)
curl -X POST -d '{"version":"0.1", "application_type":"ios", "security_hash":"some_security_hash", "username":"username09281", "email":"youremail+09281@youremailprovider.com", "password":"pass123", "password2":"pass123", "language":"en"}' https://{link}/{stage}/account/register

{"_id":"5f716e74feb00cf7fa412b21","access":"C","created_at":1601269364,"email":"youremail+09281@youremailprovider.com","is_blocked":false,"is_email_validated":false,"language":"en","success":1,"username":"username09281"}

2. Email validation link
curl https://{link}/{stage}/account/email_validation/NWY3MTZlNzRmZWIwMGNmN2ZhNDEyYjIx_NceQP4gB?language=en

Your email has been validated. You may now use your account to login to ACloudApp.org.

3. Login
curl -X POST -d '{"version":"0.1", "application_type":"ios", "security_hash":"some_security_hash", "username":"username09281", "email":"youremail+09281@youremailprovider.com", "password":"pass123"}' https://{link}/{stage}/account/login

{"_id":"5f723a018e030b0c088dfdc8","application_type":"ios","created_at":1601321473,"email":"youremail+09281@youremailprovider.com","expiry":1601753473,"is_valid":true,"jwt_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDE3NTM0NzMsImxvZ2luX2xvZ19pZCI6IjVmNzIzYTAxOGUwMzBiMGMwODhkZmRjOCIsInVzZXJfaWQiOiI1ZjcxNmU3NGZlYjAwY2Y3ZmE0MTJiMjEifQ.qNW8w3BomY2gzwx5dn0syOOvGYfBUsy5U9fNpqhVaas","login_count":1,"secret":"gBG0FuwR","success":1,"user_id":"5f716e74feb00cf7fa412b21","username":"username09281","version":"0.1"}

4. Is Logged in check
curl https://{link}/{stage}/account/is_logged_in -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDE3NTM0NzMsImxvZ2luX2xvZ19pZCI6IjVmNzIzYTAxOGUwMzBiMGMwODhkZmRjOCIsInVzZXJfaWQiOiI1ZjcxNmU3NGZlYjAwY2Y3ZmE0MTJiMjEifQ.qNW8w3BomY2gzwx5dn0syOOvGYfBUsy5U9fNpqhVaas"

{"expiry":1601753473,"success":1}

5. Logout
curl -X POST -d '{"version":"0.1", "application_type":"ios", "security_hash":"some_security_hash"}' https://{link}/{stage}/account/logout -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDE3NTM0NzMsImxvZ2luX2xvZ19pZCI6IjVmNzIzYTAxOGUwMzBiMGMwODhkZmRjOCIsInVzZXJfaWQiOiI1ZjcxNmU3NGZlYjAwY2Y3ZmE0MTJiMjEifQ.qNW8w3BomY2gzwx5dn0syOOvGYfBUsy5U9fNpqhVaas"

{"success":1}

6. Resend email validation (which sends an email with a validation link)
curl -X POST -d '{"version":"0.1", "application_type":"ios", "security_hash":"some_security_hash", "username":"username09281", "email":"youremail+09281@youremailprovider.com"}' https://{link}/{stage}/account/resend_email_validation

{"message_key":"EMAIL_VALIDATION_RESEND_SUCCESS","success":1}

7. Forgot password (which sends an email with a temporary password)
curl -X POST -d '{"version":"0.1", "application_type":"ios", "security_hash":"some_security_hash", "username":"username09281", "email":"youremail+09281@youremailprovider.com"}' https://{link}/{stage}/account/forgot_password

{"message_key":"FORGOT_PASSWORD_SUCCESS","success":1}

# Future feature targets
Acloudapp is envisioned to be a test bed for new features. This list would be maintained somewhere else as the project grows but some nice to have features for the backend would be:
- Better documentation
- Unit tests for apis end points
- Process (or even encryption) for validation / integrity check of parameters to be sure that the calling client is authorized and known
- Api throttle / defense for brute force attacks (eg hydra)

# Documentation, installation, requirements
In brief:
- You should be able to run the serverless examples (https://github.com/serverless/examples) for Golang and OpenWhisk and / or Golang and AWS specially one that would interface to a MongoDB database. 
- In the codes you should find 3 or 4 configuration.go.sample (or configuration.js.sample for the nodejs examples). Save your credentials in the files eg MongoDB, AWS access and secret keys and rename the files without the sample extension. (AWS account is also needed for AWS SES for email sending and S3 for object storage possibly in the future from among the many services). 

# Contributing
1. Contributors attest that their codes are free of any legal issues and that the project is shielded from any ramifications brought by the code in question.
2. You donate and transfer code ownership to the project. 

# License
The codes are licensed under GPLv3. 

# Donate
If the codes in anyway help you out, please consider donating to support the development of the project. 

# Hire me / us
