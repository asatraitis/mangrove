
import { apiClient as api } from "@services/apiClient/apiClient"
import { startReg } from "@services/auth/auth"


interface IRegistrationService {}
export class RegistrationService implements IRegistrationService {
    static FORM_REF_ID = "register-form"
    static CODE_INPUT_REF_ID = "code"
    // element refs:
    formRef: HTMLFormElement
    codeInputRef: HTMLInputElement

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

        const {response, error} = await api.initRegistration(this.codeInputRef.value)
        if (error) {
            // TODO: handle error
            console.error("error", error)
            return
        }
        if (!response?.publicKey) {
            // TODO: handle empty response
            console.error("error ", "missing registration options")
            return
        }

        const [cred, err] = await startReg(response.publicKey)
        if (err) {
            // TODO: handle error
            return
        }
        if (!cred) {
            // TODO: handle empty creds
            return
        }

        const {error: finishedErr} = await api.finishRegistration(response.publicKey?.user?.id, cred)
        if (finishedErr) {
            // TODO: handle error
            return
        } else

        window.location.reload()
    }
}
