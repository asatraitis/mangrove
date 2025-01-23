import {Response, MeResponse, InitLoginResponse, FinishLoginRequest} from "@dto/types"

interface IApiCLient {
    me(): Promise<Response<MeResponse>>
    initLogin(username: string): Promise<Response<InitLoginResponse>>
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
    async initLogin(username: string) {
        return ApiClient.call<InitLoginResponse>(`${this.url}${this.apiEndpoint}/login`, {method: "POST", body: JSON.stringify({username})})
    }
    async finishLogin(finishLogin: FinishLoginRequest) {
        return ApiClient.call<MeResponse>(`${this.url}${this.apiEndpoint}/login/finish`, {method: "POST", body: JSON.stringify(finishLogin)})
    }
}