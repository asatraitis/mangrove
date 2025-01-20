import { createFileRoute } from '@tanstack/react-router'
import Login from '../pages/login'

type LoginSearch = {
  redirect: string
}
export const Route = createFileRoute('/login')({
  validateSearch: (search: Record<string, unknown>): LoginSearch => {
    return {
      redirect: (search.redirect as string) || "/"
    }
  },
  component: Login,
})
