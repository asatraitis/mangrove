// Code generated by tygo. DO NOT EDIT.

//////////
// source: finish_registration_request.go

export interface FinishRegistrationRequest {
  credential: any /* protocol.CredentialCreationResponse */;
  userId: string;
}

//////////
// source: init_registration_request.go

export interface InitRegistrationRequest {
  registrationCode: string;
}

//////////
// source: init_registration_response.go

export interface InitRegistrationResponse {
  publicKey: any /* protocol.PublicKeyCredentialCreationOptions */;
}

//////////
// source: main_me_response.go

export type UserRole = string;
export const USER_ROLE_USER: UserRole = "user";
export const USER_ROLE_ADMIN: UserRole = "admin";
export const USER_ROLE_SUPERUSER: UserRole = "superadmin";
export type UserStatus = string;
export const USER_STATUS_ACTIVE: UserStatus = "active";
export const USER_STATUS_INACTIVE: UserStatus = "inactive";
export const USER_STATUS_PENDING: UserStatus = "pending";
export const USER_STATUS_SUSPENDED: UserStatus = "suspended";
export interface MeResponse {
  id: string;
  displayName: string;
  role: UserRole;
  status: UserStatus;
}

//////////
// source: response.go

export interface ResponseError {
  message: string;
  code: string;
}
export interface Response<T extends any> {
  response?: T;
  error?: ResponseError;
}
