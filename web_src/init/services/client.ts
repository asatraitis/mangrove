import {Response, InitRegistrationResponse} from "@dto/types"

export interface IClientService {
    initRegistration(registrationCode: string): Promise<Response<InitRegistrationResponse>>
}

export class ClientService implements IClientService {
    static async call<T>(url: string, config?: RequestInit): Promise<Response<T>> {
        return fetch(url, config).then(data => data.json() as Response<T>).catch(() => ({
            error: {
                message: "failed to reach the server",
                code: "ERROR_SERVER_CONN"
            },
        } as Response<T>))
    }
    async initRegistration(registrationCode: string): Promise<Response<InitRegistrationResponse>> {
        return ClientService.call<InitRegistrationResponse>(
            "/", 
            {
                method: "POST",
                body: JSON.stringify({registrationCode})
            })
    }
}