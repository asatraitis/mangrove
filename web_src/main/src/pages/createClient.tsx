import { TextInput, Textarea, NativeSelect, Container, Flex, Card, Button } from '@mantine/core';
import { FormEvent, useState } from 'react';
import {CreateClientRequest, CLIENT_STATUS_ACTIVE, CLIENT_KEY_ALGO_EDDSA} from "@dto/types"
import { useRouter } from '@tanstack/react-router';
import { apiClient as api } from '@websrc/services/apiClient/apiClient';


export default function CreateClient() {
    const router = useRouter()
    const [form, setForm] = useState<CreateClientRequest>({
        name: "",
        description: "",
        status: CLIENT_STATUS_ACTIVE,
        redirectURI: "",
        keyAlgo: CLIENT_KEY_ALGO_EDDSA,
        publicKey: ""
    })


    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault()
        const {response, error} = await api.createClient(form)
        if (error) {
            // TODO: handle error
            console.log(error)
            return
        }
        console.log({response})
        router.invalidate()
        router.history.push("/clients")
    }
    
    return (
        <Container size="xs">
            <Card p="lg" radius="md" withBorder>
                <form onSubmit={handleSubmit}>
                    <TextInput value={form.name} onChange={(e) => {setForm(data => ({...data, name: e.target.value}))}} label="Name" required />
                    <Textarea value={form.description} onChange={(e) => {setForm(data => ({...data, description: e.target.value}))}} mt="md" label="Description (optional)" />
                    <NativeSelect value={form.status} onChange={(e) => {setForm(data => ({...data, status: e.target.value}))}} mt="md" label="Status" required data={[CLIENT_STATUS_ACTIVE]} />
                    <TextInput value={form.redirectURI} onChange={(e) => {setForm(data => ({...data, redirectURI: e.target.value}))}} mt="md" label="Redirect URI" required />
                    <NativeSelect value={form.keyAlgo} onChange={(e) => {setForm(data => ({...data, keyAlgo: e.target.value}))}} mt="md" label="Key Algorithm" required data={[CLIENT_KEY_ALGO_EDDSA]} />
                    <Textarea value={form.publicKey} onChange={(e) => {setForm(data => ({...data, publicKey: e.target.value}))}} mt="md" label="Public Key" resize="vertical" required />
                    <Flex justify="flex-end">
                        <Button type="submit" mt="lg">Create Client</Button>
                    </Flex>
                </form>
            </Card>
        </Container>
    )
}

