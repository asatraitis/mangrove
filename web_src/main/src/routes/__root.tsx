import { createRootRouteWithContext} from '@tanstack/react-router'

import { AuthUser } from '../contexts/auth'
import Index from '../pages/index'

interface RouterCtx {
  auth: AuthUser
}

export const Route = createRootRouteWithContext<RouterCtx>()({
  component: Index,
})
