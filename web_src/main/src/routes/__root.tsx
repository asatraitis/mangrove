import { createRootRouteWithContext} from '@tanstack/react-router'

import Index from '../pages/index'
import { MeResponse } from '@dto/types'
import { IApiClient } from '@websrc/services/apiClient/apiClient'

interface RouterCtx {
  user: MeResponse,
  setUser: React.Dispatch<React.SetStateAction<MeResponse>>,
  api: IApiClient
}

export const Route = createRootRouteWithContext<RouterCtx>()({
  component: Index,
})
