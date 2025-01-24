import { startRegistration, startAuthentication, PublicKeyCredentialCreationOptionsJSON, RegistrationResponseJSON, PublicKeyCredentialRequestOptionsJSON, AuthenticationResponseJSON } from '@simplewebauthn/browser';
type WithError<T> = [T | null, Error | null]

export async function startAuth(publicKey: unknown): Promise<WithError<AuthenticationResponseJSON>>{
    try {
        const authRes = await startAuthentication({optionsJSON: publicKey as PublicKeyCredentialRequestOptionsJSON})
        return [authRes, null]
    } catch(err) {
        return [null, new Error(`${err}`)]
    }
}

export async function startReg(publicKey: unknown): Promise<WithError<RegistrationResponseJSON>> {
    try {
        const regRes = await startRegistration({optionsJSON: publicKey as PublicKeyCredentialCreationOptionsJSON})
        return [regRes, null]
    } catch(err) {
        return [null, new Error(`${err}`)]
    }
}