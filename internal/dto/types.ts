// Code generated by tygo. DO NOT EDIT.

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
// source: response.go

export interface ResponseError {
  message: string;
  code: string;
}
export interface Response<T extends any> {
  response?: T;
  error?: ResponseError;
}
