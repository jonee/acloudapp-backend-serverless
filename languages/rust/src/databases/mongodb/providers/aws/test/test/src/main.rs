#[path = "../../../../../../../../../rust/src/configuration/configuration.rs"] mod aca_rust_configuration;
// #[path = "../../../../../../../../../rust/src/configuration/constants.rs"] mod aca_rust_constants;
// #[path = "../../../../../../../../../rust/src/databases/mongodb/configuration/configuration.rs"] mod aca_rust_mongodb_configuration;

use lambda_http::{handler, lambda, IntoResponse, Request, Context};
use serde_json::json;

type Error = Box<dyn std::error::Error + Sync + Send + 'static>;

#[tokio::main]
async fn main() -> Result<(), Error> {
    lambda::run(handler(test)).await?;
    Ok(())
}


async fn test(_: Request, _: Context) -> Result<impl IntoResponse, Error> {
	println!("{}", aca_rust_configuration::AWS_REGION);
	
	
    // `serde_json::Values` impl `IntoResponse` by default
    // creating an application/json response
    Ok(json!({
    "message": "ACloudApp Hello"
    }))
}

/*
#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_handles() {
        let request = Request::default();
        let expected = json!({
        "message": "ACloudApp Hello"
        })
        .into_response();
        let response = test(request, Context::default())
            .await
            .expect("expected Ok(_) value")
            .into_response();
        assert_eq!(response.body(), expected.body())
    }
}
*/
