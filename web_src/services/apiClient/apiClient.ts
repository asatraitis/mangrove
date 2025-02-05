import {
    Response, 
    MeResponse, 
    InitLoginResponse, 
    FinishLoginRequest, 
    InitRegistrationResponse,
    UserClientsResponse,
    CreateClientResponse,
    CreateClientRequest,
} from "@dto/types"
import { RegistrationResponseJSON } from "@simplewebauthn/browser"

interface IApiCLient {
    me(): Promise<Response<MeResponse>>
    initRegistration(registrationCode: string): Promise<Response<InitRegistrationResponse>>
    finishRegistration(userId: string, credential: RegistrationResponseJSON): Promise<Response<unknown>>
    initLogin(username: string): Promise<Response<InitLoginResponse>>
    finishLogin(finishLogin: FinishLoginRequest): Promise<Response<MeResponse>>
    userClients(): Promise<Response<UserClientsResponse>>
    createClient(client: CreateClientRequest): Promise<Response<CreateClientResponse>>
}

export default class ApiClient implements IApiCLient {
    private url: string
    private apiEndpoint = "/v1"

    static async call<T>(url: string, config?: RequestInit): Promise<Response<T>> {
        const csrfCookie = ApiClient.getCookie("csrf_token")
        const parts = csrfCookie.split(".")
        const newConfig = {...config, headers: {...(config?.headers ?? {}), "X-CSRF-Token": parts[0]}}
        try {
            const res = await fetch(url, newConfig)
            const contentType = res.headers.get("content-type")
            if (contentType && contentType.includes('application/json')) {
                const data = await res.json() as Response<T>
                return data
            }
            return {}
        } catch(err) {
            return {error: {message: `${err}`, code: "TBD"}} as Response<T>
        }
    }
    static getCookie(name: string): string {
        if (!document.cookie) {
          return "";
        }
      
        const xsrfCookies = document.cookie.split(';')
          .map(c => c.trim())
          .filter(c => c.startsWith(name + '='));
      
        if (xsrfCookies.length === 0) {
          return "";
        }
        return decodeURIComponent(xsrfCookies[0].split('=')[1]);
    }

    constructor(url: string) {
        this.url = url
    }

    async me() {
        return ApiClient.call<MeResponse>(`${this.url}${this.apiEndpoint}/me`)
    }
    async initRegistration(registrationCode: string) {
        return ApiClient.call<InitRegistrationResponse>(`${this.url}${this.apiEndpoint}/register`, {method: "POST", body: JSON.stringify({registrationCode})})
    }
    async finishRegistration(userId: string, credential: RegistrationResponseJSON) {
        return ApiClient.call<unknown>(`${this.url}${this.apiEndpoint}/register/finish`, {method: "POST", body: JSON.stringify({userId, credential})})
    }
    async initLogin(username: string) {
        return ApiClient.call<InitLoginResponse>(`${this.url}${this.apiEndpoint}/login`, {method: "POST", body: JSON.stringify({username})})
    }
    async finishLogin(finishLogin: FinishLoginRequest) {
        return ApiClient.call<MeResponse>(`${this.url}${this.apiEndpoint}/login/finish`, {method: "POST", body: JSON.stringify(finishLogin)})
    }
    async userClients() {
        return ApiClient.call<UserClientsResponse>(`${this.url}${this.apiEndpoint}/clients`)
    }
    async createClient(client: CreateClientRequest) {
        return ApiClient.call<CreateClientResponse>(`${this.url}${this.apiEndpoint}/clients`, {method: "POST", body: JSON.stringify(client)})
    }
}

export const apiClient: IApiCLient = new ApiClient("")