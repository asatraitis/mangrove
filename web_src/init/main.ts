import { RegistrationService } from "./services/registration"

document.addEventListener("DOMContentLoaded", init)

function init() {
    try {
        new RegistrationService()
    } catch (err) {
        console.error(err, " failed to init")
    }
}

