import { ClientService, IClientService } from "./client"

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
            console.info({registrationOptions})
    }
    
}