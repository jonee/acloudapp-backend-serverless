// #[path = "../../../../../../../../../rust/src/configuration/configuration.rs"] mod aca_rust_configuration;
// #[path = "../../../../../../../../../rust/src/configuration/constants.rs"] mod aca_rust_constants;
// #[path = "../../../../../../../../../rust/src/databases/mongodb/configuration/configuration.rs"] mod aca_rust_mongodb_configuration;


extern crate rustc_serialize;
use rustc_serialize::json;
use rustc_serialize::json::Json;
use std::env;

use serde_json::json;

fn main() {

    // first arg contains JSON parameters
    if let Some(arg1) = env::args().nth(1) {
		let params: serde_json::Value = serde_json::from_str(&arg1).unwrap();

        if let Some(params_obj) = params.as_object() {

            if let Some(event_ow_headers) = params_obj.get("__ow_headers") {
		
				// stage - dev or prod
				if let Some(tmp_stage) = event_ow_headers.get("x-forwarded-url") {
					println!("tmp_stage {:?}", tmp_stage);
				}

				// authorization bearer token from headers if any
				if let Some(auth_header) = event_ow_headers.get("authorization") {
					println!("auth_header {}", auth_header);
				}
				
				if let Some(user_agent) = event_ow_headers.get("user-agent") {
					println!("user_agent {}", user_agent);
				}

            }

        }
    };

    // `serde_json::Values` impl `IntoResponse` by default
    // creating an application/json response
    println!("{}", json!({
    "message": "ACloudApp Hello"
    }))
}



