import { MantineProvider } from '@mantine/core';

import { AuthProvider } from './contexts/auth'
import { ServicesProvider } from './contexts/services';
import { ReactNode } from '@tanstack/react-router';

export default function Providers({children}: {children: ReactNode}) {
    return (
        <ServicesProvider>
            <AuthProvider>
                <MantineProvider defaultColorScheme="auto">
                    {children}
                </MantineProvider>
            </AuthProvider>
        </ServicesProvider>
    )
}