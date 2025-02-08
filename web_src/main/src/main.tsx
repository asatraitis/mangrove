import { StrictMode } from 'react'
import ReactDOM from 'react-dom/client'

import "@mantine/core/styles.css"

import Providers from './providers'
import App from './app'

const rootElement = document.getElementById('root')!
if (!rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement)
  root.render(
    <StrictMode>
      <Providers>
        <App />
      </Providers>
    </StrictMode>,
  )
}