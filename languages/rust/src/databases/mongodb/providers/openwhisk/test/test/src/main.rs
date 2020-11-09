#[path = "../../../../../../../../../rust/src/configuration/configuration.rs"] mod aca_rust_configuration;
// #[path = "../../../../../../../../../rust/src/configuration/constants.rs"] mod aca_rust_constants;
// #[path = "../../../../../../../../../rust/src/databases/mongodb/configuration/configuration.rs"] mod aca_rust_mongodb_configuration;

use serde_json::json;





fn main() {
	println!("{}", aca_rust_configuration_configuration::AWS_REGION);
	
	
    // `serde_json::Values` impl `IntoResponse` by default
    // creating an application/json response
    println!("{}", json!({
    "message": "ACloudApp Hello"
    }))
}



