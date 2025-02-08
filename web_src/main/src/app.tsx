import { RouterProvider, createRouter } from '@tanstack/react-router'


// Import the generated route tree
import { routeTree } from './routeTree.gen'
import { useAuthCtx } from './contexts/auth'
import { useServices } from './contexts/services'


// Create a new router instance
const router = createRouter({ routeTree, context: undefined! })

// Register the router instance for type safety
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

export default function App() {
    const {user, setUser} = useAuthCtx()
    const services = useServices()
    
    return (
      <RouterProvider router={router} context={{user, setUser, api: services?.api}} />
    )
}