import { Table, Box, Text } from '@mantine/core';

import { UserClientsResponse } from "@dto/types"

export type Props = {
    data: UserClientsResponse
}

export default function ClientsTable({data}: Props) {
return (
            <Box mt="md">
                <Table withTableBorder highlightOnHover>
                    <Table.Thead>
                        <Table.Tr>
                            <Table.Th>ID</Table.Th>
                            <Table.Th>Name</Table.Th>
                            <Table.Th>Description</Table.Th>
                            <Table.Th>Redirect URI</Table.Th>
                            <Table.Th>Status</Table.Th>
                        </Table.Tr>
                    </Table.Thead>
                    <Table.Tbody>
                        {
                            data?.map(c => (
                                <Table.Tr key={c.id}>
                                    <Table.Td><Text truncate="end">{c.id}</Text></Table.Td>
                                    <Table.Td>{c.name}</Table.Td>
                                    <Table.Td>{c.description}</Table.Td>
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