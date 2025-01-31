import { getRouteApi, Link } from '@tanstack/react-router'
import { SlFolderAlt } from "react-icons/sl";

import { Table, Box, Flex, Text, Card, Button, ThemeIcon } from '@mantine/core';

const clientApi = getRouteApi("/clients")
export default function Clients() {
    const {response: clients, error} = clientApi.useLoaderData()

    if (error) {
        console.error(error.message)
    }
    console.log({clients, error})
    

    if (clients?.length === 0) {
        return (
            <Card radius="md" p='xl' withBorder>
                <Flex justify="center" direction="column" align="center">
                    <ThemeIcon variant="transparent" size={100}>
                        <SlFolderAlt size={80} />
                    </ThemeIcon>
                    <Text size="xl">No Clients Created</Text>
                    <Text>Created clients will appear below.</Text>
                        <Link to="/clients/new">
                            <Button mt="md">Create Client</Button>
                        </Link>
                </Flex>
            </Card>
        )
    }

    return (
        <Box>
            <Table withTableBorder highlightOnHover>
                <Table.Thead>
                    <Table.Tr>
                        <Table.Th>Name</Table.Th>
                        <Table.Th>Description</Table.Th>
                        <Table.Th>Type</Table.Th>
                        <Table.Th>Redirect URI</Table.Th>
                        <Table.Th>Status</Table.Th>
                    </Table.Tr>
                </Table.Thead>
                <Table.Tbody>
                    {
                        clients?.map(c => (
                            <Table.Tr key={c.id}>
                                <Table.Td>{c.name}</Table.Td>
                                <Table.Td>{c.description}</Table.Td>
                                <Table.Td>{c.type}</Table.Td>
                                <Table.Td>{c.redirectURI}</Table.Td>
                                <Table.Td>{c.status}</Table.Td>
                            </Table.Tr>
                        ))
                    }
                </Table.Tbody>
            </Table>
        </Box>
    )
}
