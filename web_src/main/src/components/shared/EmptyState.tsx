import { ReactNode } from "react"
import { Flex, Text, Title, ThemeIcon } from '@mantine/core';
import { SlFolderAlt } from "react-icons/sl";

export type Props = {
    title?: string
    description?: string
    children?: ReactNode
}
// TODO: take a look at icon - make it a prop? 
export default function EmptyState({title, description, children}: Props) {
    return (
        <Flex justify="center" direction="column" align="center">
            <ThemeIcon variant="transparent" size={100}>
                <SlFolderAlt size={80} />
            </ThemeIcon>
            {title && <Title order={3}>{title}</Title>}
            {description && <Text>{description}</Text>}
            {children}
        </Flex>
    )
}