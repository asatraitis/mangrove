import {Response, InitRegistrationRequest, InitRegistrationResponse} from "@dto/types"

export interface IClientService {
    initRegistration(registrationCode:string): Promise<Response<InitRegistrationResponse>>
}

export class ClientService implements IClientService {

    async initRegistration(registrationCode: string): Promise<Response<InitRegistrationResponse>> {
        const payload: InitRegistrationRequest = {registrationCode}
        const response = await fetch("/", {
            method: "POST",
            body: JSON.stringify(payload)
        })
        // TODO: check for status code
        return await response.json() as Response<InitRegistrationResponse>
    }
}