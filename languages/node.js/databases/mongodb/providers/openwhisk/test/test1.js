/*
sample binary file to base64

var fs = require('fs');

var binary = fs.readFileSync('./some.png');

var image_string = binary.toString('base64');

console.log(image_string);

buff = Buffer.from(image_string, 'base64');

fs.writeFile('some.png', buff, 'binary', function(err){});

*/

/*
sample image base64 string to s3 upload 
const AWS = require('aws-sdk');

  // image_string if any
  let imageLink;
  if (image_string && image_file_name && payload.user_id) {
    let buff = Buffer.from(image_string, 'base64');
    // console.log(buff.toString('base64'));

    let newFileName = payload.user_id + "_" + new Date().getTime() + "_" + image_file_name; // TODO sanitize filename

    // upload to s3
    const s3 = new AWS.S3({
      accessKeyId: acaConfiguration.AWS_ACCESS_KEY_ID,
      secretAccessKey: acaConfiguration.AWS_SECRET_ACCESS_KEY
    });

    const params = {
      Bucket: acaConfiguration.AWS_S3_BUCKET,
      Key: stage + "/" + newFileName,
      Body: buff,
      ACL: 'public-read-write'
    };

    imageLink = 'https://' + acaConfiguration.AWS_S3_BUCKET + '.s3.amazonaws.com/' + stage + "/" + newFileName; // anticipated since there is a race / async bug here

    // Uploading files to the bucket
    await s3.upload(params, function(err, data) {
        if (err) {
            console.log(err);
            imageLink = ""; // remove anticipated link
        }

        // console.log(`File uploaded successfully. ${data.Location}`);
        console.log(data);

        if (imageLink != data.Location) {
          console.log("s3 wrong guess");
          imageLink = data.Location;
        }
    });

/*
    const res = await new Promise((resolve, reject) => {
      s3.upload(params, (err, data) => err == null ? resolve(data) : reject(err));
    });
    if (res.location) {
      imageLink = res.location;
    }
* /
    console.log("imageLink:", imageLink);    
  }
*/


