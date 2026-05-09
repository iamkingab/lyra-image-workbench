import React from 'react'
import ReactDOM from 'react-dom/client'
import { AdminPage } from './components/AdminPage'
import { WorkbenchPage } from './components/WorkbenchPage'
import './styles.css'

function App() {
  return window.location.pathname === '/admin' ? <AdminPage /> : <WorkbenchPage />
}

ReactDOM.createRoot(document.getElementById('root')!).render(<App />)
