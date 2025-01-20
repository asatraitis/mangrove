import {Response, MeResponse} from "@dto/types"

interface IApiCLient {
    me(): Promise<Response<MeResponse>>
}

export default class ApiClient implements IApiCLient {
    private url: string

    static async call<T>(url: string, config?: RequestInit): Promise<Response<T>> {
        try {
            const csrfCookie = ApiClient.getCookie("csrf_token")
            const parts = csrfCookie.split(".")
            const newConfig = {...config, headers: {...(config?.headers ?? {}), "X-CSRF-Token": parts[0]}}

            const res = await fetch(url, newConfig)
            const contentType = res.headers.get("content-type")
            if (contentType && contentType.includes('application/json')) {
                const data = await res.json() as Response<T>
                if (!res.ok) {
                    throw new Error(data.error?.message)
                }
                if (data.error) {
                    throw new Error(data.error?.message)
                }
                return data
            }
            return {}
        } catch(err) {
            throw new Error(`request to ${url} failed. ${err}`)
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
        return ApiClient.call<MeResponse>(`${this.url}/v1/me`)
    }
}