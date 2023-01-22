import React from "react"
import { Outlet } from 'react-router-dom'

const Header = () => {
    return (
      <div>
        <h1>Queue System</h1>
        <div className="content">
          <Outlet />
        </div>
      </div>
    )
}

export {
    Header
}