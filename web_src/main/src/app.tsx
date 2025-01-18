import { RouterProvider, createRouter } from '@tanstack/react-router'
import { MantineProvider } from '@mantine/core';
import "@mantine/core/styles.css"

// Import the generated route tree
import { routeTree } from './routeTree.gen'

import { AuthCtx, defaultAuthUser } from './contexts/auth';

// Create a new router instance
const router = createRouter({ routeTree, context: undefined! })

// Register the router instance for type safety
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

export default function App() {
    return (
        <AuthCtx.Provider value={defaultAuthUser}>
            <MantineProvider defaultColorScheme="auto">
                <RouterProvider router={router} context={{auth: defaultAuthUser}} />
            </MantineProvider>
        </AuthCtx.Provider>
    )
}