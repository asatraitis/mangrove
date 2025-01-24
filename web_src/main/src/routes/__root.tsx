import { createRootRouteWithContext} from '@tanstack/react-router'

import Index from '../pages/index'
import { MeResponse } from '@dto/types'

interface RouterCtx {
  user: MeResponse,
  setUser: React.Dispatch<React.SetStateAction<MeResponse>>
}

export const Route = createRootRouteWithContext<RouterCtx>()({
  component: Index,
})
