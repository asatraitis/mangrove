import {Response, InitRegistrationResponse} from "@dto/types"

export interface IClientService {
    initRegistration(registrationCode: string): Promise<Response<InitRegistrationResponse>>
    finishRegistration(userId: string, credential: Credential): Promise<any>
}

export class ClientService implements IClientService {
    static async call<T>(url: string, config?: RequestInit): Promise<Response<T>> {
        return fetch(url, config).then(data => {
            const contentType = data.headers.get('content-type');
            if (contentType && contentType.includes('application/json')) {
                return data.json() as Response<T>
            }
            return {} as Response<T>
        }).catch((err) => {
            return ({
            error: {
                message: err,
                code: "ERROR_FETCH_INCOMPLETE"
            },
        } as Response<T>)})
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
    
    initRegistration(registrationCode: string): Promise<Response<InitRegistrationResponse>> {
        return ClientService.call<InitRegistrationResponse>(
            "/", 
            {
                method: "POST",
                body: JSON.stringify({registrationCode})
            }
        )
    }

    finishRegistration(userId: string, credential: Credential): Promise<Response<any>> {
        const csrfCookie = ClientService.getCookie("csrf_token")
        if (csrfCookie === "") {
            // TODO: handle no cookie
            console.error("missign CSRF cookie")
        }

        const parts = csrfCookie.split(".")
        if (parts.length != 2) {
            // TODO: handle bad cookie
            console.error("bad CSRF cookie")
        }


        return ClientService.call<any>("/finish", {
            method: "POST",
            headers: {
                "X-CSRF-Token": parts[0]
            },
            body: JSON.stringify({userId, credential})
        })
    }
}