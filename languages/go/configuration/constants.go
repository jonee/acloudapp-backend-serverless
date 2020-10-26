// Created by Jonee Ryan Ty
// Copyright ACloudApp

/**
 * Constants
 */

package configuration

const (
	// should agree with client and server values eg server constants.go, ios Constants.swift
	/* server + client set values */
	PASSWORD_MIN_LENGTH = 6
	// USERNAME_MAX_LENGTH      = 12
	USERNAME_MIN_LENGTH = 4

	// REGEX_USERNAME_CLIENT = "^@[-0-9a-zA-Z\\_]*$" // client side starts with @, server side does not
	REGEX_USERNAME_CLIENT = "^[-0-9a-zA-Z\\_]*$" // letters, numbers and underscore only
	REGEX_USERNAME_SERVER = "^[-0-9a-zA-Z\\_]*$"

	REGEX_EMAIL = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	// REGEX_EMAIL = "[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]"

	// REGEX_WEBSITE = "((?:http|https)://)?(?:www\\.)?[\\w\\d\\-_]+\\.\\w{2,3}(\\.\\w{2})?(/(?<=/)(?:[\\w\\d\\-./_]+)?)?"

	/* server configs */
	BRANDING  = "ACloudApp.org"
	BRANDING2 = "acloudapp" // this is used for url eg acloudapp://

	PASSWORD_TEMPORARY_EXPIRY = 12     // hours
	LOGGED_IN_EXPIRY          = 24 * 5 // 5 days

	DEV_CHEAT_SLUG             = "123xyz123"
	LOGGED_IN_EXPIRY_DEV_CHEAT = 24 * 365

	DEFAULT_LANGUAGE = "en"
)

/*
const (
	AclUnknown = 0
	AclOpen    = 1
	AclUser    = 2
	AclManager = 3
	AclAdmin   = 4
)

var AclString = map[int]string{
	AclUnknown: "Unknown",
	AclOpen:    "Open",
	AclUser:    "User",
	AclManager: "Manager",
	AclAdmin:   "Admin",
}
*/
