import {Response} from "@dto/types"
import { ClientService, IClientService } from "./client"

class UIError {
    message: string
    constructor(message: string) {
        this.message = message
    }
}

interface IRegistrationService {}

export class RegistrationService implements IRegistrationService {
    static FORM_REF_ID = "register-form"
    static CODE_INPUT_REF_ID = "code"
    // element refs:
    formRef: HTMLFormElement
    codeInputRef: HTMLInputElement

    // client
    client: IClientService


    constructor() {
        // Get element refs
        this.formRef = document.getElementById(RegistrationService.FORM_REF_ID) as HTMLFormElement
        if (!this.formRef) {
            throw new Error("failed to reference form in DOM: "+RegistrationService.FORM_REF_ID)
        }
        this.codeInputRef = this.formRef.querySelector("#code") as HTMLInputElement
        if (!this.codeInputRef) {
            throw new Error("failed to reference code input in DOM: "+RegistrationService.CODE_INPUT_REF_ID)
        }

        // init client service
        this.client = new ClientService()

        // add event handlers
        this.formRef.addEventListener("submit", this.handleSubmit.bind(this))
    }

    async handleSubmit(e: SubmitEvent) {
        e.preventDefault()
        // validate registration code
        if (this.codeInputRef.value.length != 6) {
            // TODO: set and display error
            return
        }

        const registrationOptions = await this.client.initRegistration(this.codeInputRef.value)
        if (registrationOptions.error) {
            // TODO: handle error
            console.error("error", registrationOptions.error)
            return
        }
        if (!registrationOptions.response?.publicKey) {
            // TODO: handle empty response
            console.error("error ", "missing registration options")
            return
        }

        const [prepedCredOpts, err] = typeConvPublicKeyCredentialCreationOptions(registrationOptions.response.publicKey)
        if (err != null) {
            // TODO: handle error
            console.error(err.message)
            return
        }
        if (prepedCredOpts == null) {
            // TODO: handle empty credopts
            return
        }

        const creds = await navigator.credentials.create({
            publicKey: prepedCredOpts,
        });
        console.info({creds})
        if (!creds) {
            console.error("missing registration credential")
            return
        }
        const prepedCred = prepCredentialRequest(creds)

        const finished: Response<any> = await this.client.finishRegistration(registrationOptions.response.publicKey.user.id, prepedCred)
        if (finished.error) {
            // TODO: handle error
            return
        } else

        window.location.reload()
    }

    
}

type CredOpts = [
    PublicKeyCredentialCreationOptions | null,
    UIError | null,
]
function typeConvPublicKeyCredentialCreationOptions(opts: any): CredOpts {
    if (!opts) {
        return [null, new UIError("missing opts")]
    }
    try {
        const credopts: PublicKeyCredentialCreationOptions = {
            challenge: base64UrlToBuffer(opts.challenge),
            rp: opts.rp,
            user: {
                id: base64UrlToBuffer(opts.user.id),
                name: opts.user.name,
                displayName: opts.user.displayName
            },
            pubKeyCredParams: opts.pubKeyCredParams,
            authenticatorSelection: {
                authenticatorAttachment: "cross-platform"
            },
            timeout: opts.timeout,
            attestation: "direct"
        }
        return [credopts, null]
    } catch {
        return [null, new UIError("missing opts")]
    }
}

function prepCredentialRequest(cred: any): any {
    if (!cred?.response?.attestationObject) {
        // TODO: handle this
    }
    const attestationObject = new Uint8Array(cred.response.attestationObject)

    if (!cred?.response?.clientDataJSON) {
        // TODO: handle this
    }
    const clientDataJSON = new Uint8Array(cred.response.clientDataJSON)

    if (!cred.rawId) {
        // TODO: handle this
    }
    const rawId = new Uint8Array()

    if (!cred.id && !cred.type) {
         // TODO: handle this
    }

    return {
        id: cred.id,
        type: cred.type,
        response: {
            attestationObject: bufferToBase64Url(attestationObject),
            clientDataJSON: bufferToBase64Url(clientDataJSON)
        },
        rawId: bufferToBase64Url(rawId)
    }
}

function base64UrlToBuffer(base64Url: string): ArrayBuffer {
    base64Url = base64Url.replace(/-/g, "+").replace(/_/g, "/")
    let binaryString = window.atob(base64Url)
    let len = binaryString.length
    let bytes = new Uint8Array(len)
    for (let i = 0; i < len; i++) {
        bytes[i] = binaryString.charCodeAt(i)
    }
    return bytes.buffer
}

function bufferToBase64Url(buffer: ArrayBuffer): string {
    let binary = ""
    let bytes = new Uint8Array(buffer)
    let len = bytes.byteLength
    for (let i = 0; i < len; i++) {
        binary += String.fromCodePoint(bytes[i])
    }
    return window.btoa(binary)
        .replace(/\+/g, "-")
        .replace(/\//g, "_")
        .replace(/=+$/, "")
}