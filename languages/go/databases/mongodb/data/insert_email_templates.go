// Created by Jonee Ryan Ty
// Copyright ACloudApp

package main

import (
	"context"
	"time"

	acaGoMongoDBModels "acloudapp.org/databases/mongodb/models"
	acaGoMongoDBUtilities "acloudapp.org/databases/mongodb/utilities"

	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	var mongoClient *mongo.Client

	mapStore := make(map[string]interface{})
	mapStore["last_milestone_time"] = time.Now()
	mapStore["stage"] = "dev"
	mapStore["mongoClient"] = mongoClient
	acaGoMongoDBUtilities.DoDBConnect(mapStore)
	mongoClient = mapStore["mongoClient"].(*mongo.Client) // save it back

	var etObj1 acaGoMongoDBModels.EmailTemplate
	etObj1.Name = "VALIDATION"
	etObj1.Language = "en"
	etObj1.Subject = `%%BRANDING%% validate email`
	etObj1.TemplateHTML = `
<p>%%subject%% <br />
<br />
Hi %%username%%, <br />
<br />
Thank you very much for your registration. <br />
<br />
Please click the following link to validate your account. <br />
<a href="%%validation_link%%">%%validation_link%%</a> <br />
<br />
Thanks! <br />
%%BRANDING%% <br />
<br />
** PLEASE DO NOT REPLY TO THIS EMAIL as this is automatically generated and you will not receive a response.  **
</p>
`
	etObj1.TemplateText = `
%%subject%%

Hi %%username%%,

Thank you very much for your registration. 

Please access the following link to validate your account.
%%validation_link%%

Thanks!
%%BRANDING%%

** PLEASE DO NOT REPLY TO THIS EMAIL as this is automatically generated and you will not receive a response.  **
`
	etObj1.CreatedAt = time.Now()
	etObj1.Save(mapStore)

	var etObj2 acaGoMongoDBModels.EmailTemplate
	etObj2.Name = "EMAIL_CHANGED_VALIDATION"
	etObj2.Language = "en"
	etObj2.Subject = `%%BRANDING%% email changed`
	etObj2.TemplateHTML = `
<p>%%subject%% <br />
<br />
Hi %%username%%, <br />
<br />
You recently changed the email for your account in %%BRANDING%%. <br />
<br />
Please click the following link to revalidate your account. <br />
<a href="%%validation_link%%">%%validation_link%%</a> <br />
<br />
Thanks!
%%BRANDING%% <br />
<br />
** PLEASE DO NOT REPLY TO THIS EMAIL as this is automatically generated and you will not receive a response.  **
</p>
`
	etObj2.TemplateText = `
%%subject%%

Hi %%username%%,

You recently changed the email for your account in %%BRANDING%%.

Please click the following link to revalidate your account.
%%validation_link%%

Thanks!
%%BRANDING%%

** PLEASE DO NOT REPLY TO THIS EMAIL as this is automatically generated and you will not receive a response.  **
`
	etObj2.CreatedAt = time.Now()
	etObj2.Save(mapStore)

	var etObj3 acaGoMongoDBModels.EmailTemplate
	etObj3.Name = "RESEND_VALIDATION"
	etObj3.Language = "en"
	etObj3.Subject = `%%BRANDING%% resend validate email`
	etObj3.TemplateHTML = `
<p>%%subject%% <br />
<br />
Hi %%username%%, <br />
<br />
Thank you very much for your registration. <br />
<br />
Please click the following link to validate your account. <br />
<a href="%%validation_link%%">%%validation_link%%</a> <br />
<br />
Thanks! <br />
%%BRANDING%% <br />
<br />
** PLEASE DO NOT REPLY TO THIS EMAIL as this is automatically generated and you will not receive a response.  **
</p>
`
	etObj3.TemplateText = `
%%subject%%

Hi %%username%%,

Thank you very much for your registration. 

Please access the following link to validate your account.
%%validation_link%%

Thanks!
%%BRANDING%%

** PLEASE DO NOT REPLY TO THIS EMAIL as this is automatically generated and you will not receive a response.  **
`
	etObj3.CreatedAt = time.Now()
	etObj3.Save(mapStore)

	var etObj4 acaGoMongoDBModels.EmailTemplate
	etObj4.Name = "FORGOT_PASSWORD"
	etObj4.Language = "en"
	etObj4.Subject = `%%BRANDING%% forgot password`
	etObj4.TemplateHTML = `
<p>%%subject%% <br />
<br />
Hi %%username%%, <br />
<br />
You have recently requested a new password. (If this is a mistake then it should be safe to ignore this email.) <br />
<br />
Your new password is %%new_password%% . <br />
<br />
Thanks! <br />
%%BRANDING%% <br />
<br />
** PLEASE DO NOT REPLY TO THIS EMAIL as this is automatically generated and you will not receive a response.  **
</p>
`
	etObj4.TemplateText = `
%%subject%%

Hi %%username%%,

You have recently requested a new password. (If this is a mistake then it should be safe to ignore this email.)

Your new password is %%new_password%% .

Thanks!
%%BRANDING%%

** PLEASE DO NOT REPLY TO THIS EMAIL as this is automatically generated and you will not receive a response.  **
`
	etObj4.CreatedAt = time.Now()
	etObj4.Save(mapStore)

	var etObj5 acaGoMongoDBModels.EmailTemplate
	etObj5.Name = "MESSAGE_NOTIFICATION"
	etObj5.Language = "en"
	etObj5.Subject = `%%BRANDING%% new message notification`
	etObj5.TemplateHTML = `
<p>%%subject%% <br />
<br />
Hi %%username%%, <br />
<br />
You have received a new private message from %%from%%. <br />
<br />
Thanks! <br />
%%BRANDING%% <br />
<br />
** PLEASE DO NOT REPLY TO THIS EMAIL as this is automatically generated and you will not receive a response.  **
</p>
`
	etObj5.TemplateText = `
%%subject%%

Hi %%username%%,

You have received a new private message from %%from%%.

Thanks!
%%BRANDING%%

** PLEASE DO NOT REPLY TO THIS EMAIL as this is automatically generated and you will not receive a response.  **
`
	etObj5.CreatedAt = time.Now()
	etObj5.Save(mapStore)

	defer mongoClient.Disconnect(context.TODO())

	// main
}
